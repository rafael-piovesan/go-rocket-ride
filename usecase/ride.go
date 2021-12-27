package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/audit"
	"github.com/rafael-piovesan/go-rocket-ride/entity/idempotency"
	"github.com/rafael-piovesan/go-rocket-ride/entity/originip"
	"github.com/rafael-piovesan/go-rocket-ride/entity/stagedjob"
	"github.com/rafael-piovesan/go-rocket-ride/pkg/stripemock"
)

type rideUseCase struct {
	cfg   rocketride.Config
	store rocketride.Datastore
}

func NewRideUseCase(cfg rocketride.Config, ds rocketride.Datastore) rocketride.RideUseCase {
	// setup Stripe's key
	stripe.Key = cfg.StripeKey

	// Replace the original Stripe API Backend with its mock
	stripemock.Init()

	return &rideUseCase{
		cfg:   cfg,
		store: ds,
	}
}

func (r *rideUseCase) Create(
	ctx context.Context,
	ik *entity.IdempotencyKey,
	rd *entity.Ride,
) (*entity.IdempotencyKey, error) {
	var key *entity.IdempotencyKey
	var ride *entity.Ride
	var err error

	key, err = r.getIdempotencyKey(ctx, ik)
	if err != nil {
		return nil, err
	}

	defer func() {
		// If we're leaving under an error condition, try to unlock the idempotency
		// key right away so that another request can try again.
		if err != nil && key != nil {
			r.unlockIdempotencyKey(ctx, key)
		}
	}()

	for {
		switch key.RecoveryPoint {
		case idempotency.RecoveryPointStarted:
			ride, err = r.createRide(ctx, key, rd)
			if err != nil {
				return nil, err
			}

		case idempotency.RecoveryPointCreated:
			err = r.createCharge(ctx, key, ride)
			if err != nil {
				return nil, err
			}

		case idempotency.RecoveryPointCharged:
			err = r.sendReceipt(ctx, key)
			if err != nil {
				return nil, err
			}

		case idempotency.RecoveryPointFinished:
			return key, nil

		default:
			return nil, entity.ErrIdemKeyUnknownRecoveryPoint
		}
	}
}

func (r *rideUseCase) getIdempotencyKey(
	ctx context.Context,
	ik *entity.IdempotencyKey,
) (*entity.IdempotencyKey, error) {
	var err error

	// may be created on this request or retrieved if it already exists
	var key *entity.IdempotencyKey

	// Our first atomic phase to create or update an idempotency key.
	//
	// A key concept here is that if two requests try to insert or update within
	// close proximity, one of the two will be aborted by Postgres because we're
	// using a transaction with SERIALIZABLE isolation level. It may not look
	// it, but this code is safe from races.
	err = r.store.Atomic(ctx, func(ds rocketride.Datastore) error {
		key, err = ds.GetIdempotencyKey(ctx, ik.IdempotencyKey, ik.UserID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				now := time.Now().UTC()
				ik.LastRunAt = now
				ik.LockedAt = &now
				ik.RecoveryPoint = idempotency.RecoveryPointStarted
				key, err = ds.CreateIdempotencyKey(ctx, ik)
			}
			return err
		}

		// Unmarshal the JSON returned from datastore, so we're able to
		// properly compare it against the request.
		rd1, rd2 := entity.Ride{}, entity.Ride{}
		if err := json.Unmarshal(ik.RequestParams, &rd1); err != nil {
			return entity.ErrIdemKeyParamsMismatch
		}

		if err := json.Unmarshal(key.RequestParams, &rd2); err != nil {
			return entity.ErrIdemKeyParamsMismatch
		}

		// Programs sending multiple requests with different parameters but the
		// same idempotency key is a bug.
		if !reflect.DeepEqual(rd1, rd2) {
			return entity.ErrIdemKeyParamsMismatch
		}

		// Only acquire a lock if the key is unlocked or its lock has expired
		// because it was long enough ago.
		timeout := time.Duration(r.cfg.IdemKeyTimeout) * time.Second
		if key.LockedAt != nil && key.LockedAt.After(time.Now().UTC().Add(-1*timeout)) {
			return entity.ErrIdemKeyRequestInProgress
		}

		// Lock the key and update latest run unless the request is already
		// finished.
		if key.RecoveryPoint != idempotency.RecoveryPointFinished {
			now := time.Now().UTC()
			key.LastRunAt = now
			key.LockedAt = &now
			key, err = ds.UpdateIdempotencyKey(ctx, key)
			return err
		}
		return nil
	})

	return key, err
}

func (r *rideUseCase) createRide(
	ctx context.Context,
	ik *entity.IdempotencyKey,
	rd *entity.Ride,
) (*entity.Ride, error) {
	var err error
	var ride *entity.Ride

	oip := originip.FromCtx(ctx)

	err = r.store.Atomic(ctx, func(ds rocketride.Datastore) error {
		rd.IdempotencyKeyID = &ik.ID
		rd.UserID = ik.UserID
		ride, err = ds.CreateRide(ctx, rd)
		if err != nil {
			return err
		}

		// in the same transaction insert an audit record for what happened
		ar := &entity.AuditRecord{
			Action:       audit.ActionCreateRide,
			CreatedAt:    time.Now().UTC(),
			Data:         ik.RequestParams,
			OriginIP:     oip.IP,
			ResourceID:   rd.ID,
			ResourceType: audit.ResourceTypeRide,
			UserID:       ik.UserID,
		}
		_, err = ds.CreateAuditRecord(ctx, ar)
		if err != nil {
			return err
		}

		ik.RecoveryPoint = idempotency.RecoveryPointCreated
		_, err = ds.UpdateIdempotencyKey(ctx, ik)
		return err
	})

	return ride, err
}

func (r *rideUseCase) createCharge(ctx context.Context, ik *entity.IdempotencyKey, rd *entity.Ride) error {
	var err error
	var ride *entity.Ride

	err = r.store.Atomic(ctx, func(ds rocketride.Datastore) error {
		handleStripeErr := func(resCode *idempotency.ResponseCode, resBody *idempotency.ResponseBody) {
			ik.LockedAt = nil
			ik.ResponseCode = resCode
			ik.ResponseBody = resBody
			// short-circuit to the final state
			ik.RecoveryPoint = idempotency.RecoveryPointFinished
			_, err = ds.UpdateIdempotencyKey(ctx, ik)
			if err != nil {
				log.Errorf("error updating idem key after stripe error: %v", err)
			}
		}

		// check if we're restarting from a Recovery Point and retrieve a ride from db
		if rd == nil {
			ride, err = ds.GetRideByIdempotencyKeyID(ctx, ik.ID)
			if err != nil {
				return err
			}
		} else {
			ride = rd
		}

		// Pass through our own unique ID rather than the value transmitted
		// to us so that we can guarantee uniqueness to Stripe across all
		// Rocket Rides accounts.
		stripeIK := fmt.Sprintf("go-rocket-ride-%v", ik.ID)
		customerID := ik.User.StripeCustomerID

		// Rocket Rides is still a new service, so during our prototype phase
		// we're going to give $20 fixed-cost rides to everyone, regardless of
		// distance. We'll implement a better algorithm later to better
		// represent the cost in time and jetfuel on the part of our pilots.
		params := &stripe.ChargeParams{
			Params:      stripe.Params{IdempotencyKey: &stripeIK},
			Amount:      stripe.Int64(2000),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Customer:    &customerID,
			Description: stripe.String(fmt.Sprintf("Charge for ride %v", ride.ID)),
		}

		c, err := charge.New(params)
		if err != nil {
			if stripeErr, ok := err.(*stripe.Error); ok {
				var resCode idempotency.ResponseCode
				var resBody idempotency.ResponseBody
				defer handleStripeErr(&resCode, &resBody)

				if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
					resCode = idempotency.ResponseCodeErrPayment
					resBody = idempotency.ResponseBody{Message: entity.ErrPaymentProvider.Error()}

					log.Errorf("stripe card error: %v", cardErr.Error())
					return entity.ErrPaymentProvider
				}

				resCode = idempotency.ResponseCodeErrPaymentGeneric
				resBody = idempotency.ResponseBody{Message: entity.ErrPaymentProviderGeneric.Error()}

				log.Errorf("stripe api error: %v", stripeErr.Error())
				return entity.ErrPaymentProviderGeneric
			}

			log.Errorf("stripe api request error: %v", err)
			return err
		}

		ride.StripeChargeID = &c.ID
		log.Debugf("stripe charge id: %v", c.ID)
		_, err = ds.UpdateRide(ctx, ride)
		if err != nil {
			return err
		}

		ik.RecoveryPoint = idempotency.RecoveryPointCharged
		_, err = ds.UpdateIdempotencyKey(ctx, ik)
		return err
	})
	return err
}

func (r *rideUseCase) sendReceipt(ctx context.Context, ik *entity.IdempotencyKey) error {
	// Send a receipt asynchronously by adding an entry to the staged_jobs
	// table. By funneling the job through Postgres, we make this
	// operation transaction-safe.
	err := r.store.Atomic(ctx, func(ds rocketride.Datastore) error {
		jobArgs := stagedjob.JobArgReceipt{
			Amount:   int64(20),
			Currency: "usd",
			UserID:   ik.UserID,
		}

		args, err := json.Marshal(jobArgs)
		if err != nil {
			return err
		}

		sj := &entity.StagedJob{
			JobName: stagedjob.JobNameSendReceipt,
			JobArgs: args,
		}
		_, err = ds.CreateStagedJob(ctx, sj)
		if err != nil {
			return err
		}

		resCode := idempotency.ResponseCodeOK
		resBody := idempotency.ResponseBody{Message: "OK"}

		ik.LockedAt = nil
		ik.ResponseCode = &resCode
		ik.ResponseBody = &resBody
		ik.RecoveryPoint = idempotency.RecoveryPointFinished
		_, err = ds.UpdateIdempotencyKey(ctx, ik)
		return err
	})
	return err
}

func (r *rideUseCase) unlockIdempotencyKey(ctx context.Context, ik *entity.IdempotencyKey) {
	ik.LockedAt = nil
	_, err := r.store.UpdateIdempotencyKey(ctx, ik)
	if err != nil {
		log.Errorf("unlock idem key error: %v", err)
	}
}

var _ rocketride.RideUseCase = (*rideUseCase)(nil)
