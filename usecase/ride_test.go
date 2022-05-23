//go:build unit
// +build unit

package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/idempotency"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity/originip"
	mocks "github.com/rafael-piovesan/go-rocket-ride/v2/mocks/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/config"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v72"
	"gopkg.in/h2non/gock.v1"
)

const (
	stripeURL = "http://stripeapi"
)

type testMocks struct {
	store   *mocks.Store
	idemKey *mocks.IdempotencyKey
	ride    *mocks.Ride
	audit   *mocks.AuditRecord
	user    *mocks.User
	job     *mocks.StagedJob
}

func init() {
	maxRetries := int64(0)
	stripeMockBackend := stripe.GetBackendWithConfig(
		stripe.APIBackend,
		&stripe.BackendConfig{
			URL:               stripe.String(stripeURL),
			LeveledLogger:     stripe.DefaultLeveledLogger,
			MaxNetworkRetries: &maxRetries,
		},
	)
	stripe.SetBackend(stripe.APIBackend, stripeMockBackend)
	stripe.SetBackend(stripe.UploadsBackend, stripeMockBackend)
}

func getMocks() testMocks {
	return getMocksWithTimes(1)
}

func getMocksWithTimes(n int) testMocks {
	m := testMocks{
		store:   &mocks.Store{},
		idemKey: &mocks.IdempotencyKey{},
		ride:    &mocks.Ride{},
		audit:   &mocks.AuditRecord{},
		user:    &mocks.User{},
		job:     &mocks.StagedJob{},
	}

	mockAS := &mocks.AtomicStore{}
	mockAS.On("AuditRecords").Return(m.audit)
	mockAS.On("IdempotencyKeys").Return(m.idemKey)
	mockAS.On("Rides").Return(m.ride)
	mockAS.On("StagedJobs").Return(m.job)
	mockAS.On("Users").Return(m.user)

	var mockAtomic *mock.Call
	mockAtomic = m.store.On("Atomic", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(datastore.AtomicBlock)
			if !ok {
				panic("argument mismatch")
			}

			// Call the actual func argument 'fn' passed in to
			// 'Atomic(context.Context, datastore.AtomicBlock) error'
			// as expected from its second parameter and, while doing so, inject the
			// mocked Datastore instance 'mockDS' so we're able to test the other calls
			// made to it inside the 'Atomic' block.
			mockAtomic.Return(fn(mockAS))
		})

	if n >= 0 {
		mockAtomic.Times(n)
	}

	return m
}

func TestGetIdempotencyKey(t *testing.T) {
	ctx := context.Background()

	mockCfg := config.Config{IdemKeyTimeout: 5}

	gofakeit.Seed(time.Now().UnixNano())
	jsonRide, err := json.Marshal(entity.Ride{
		OriginLat: gofakeit.Float64(),
		OriginLon: gofakeit.Float64(),
		TargetLat: gofakeit.Float64(),
		TargetLon: gofakeit.Float64(),
	})
	require.NoError(t, err)

	t.Run("Error on GetIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
		}

		retErr := errors.New("err GetIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(entity.IdempotencyKey{}, retErr)

		_, err := uc.getIdempotencyKey(ctx, &ik)

		assert.Equal(t, retErr, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Error on CreateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
		}

		retErr := errors.New("err CreateIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(entity.IdempotencyKey{}, data.ErrRecordNotFound)

		m.idemKey.On("Save", ctx, &ik).
			Once().
			Return(retErr)

		_, err := uc.getIdempotencyKey(ctx, &ik)

		assert.Equal(t, retErr, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("Success on CreateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(entity.IdempotencyKey{}, data.ErrRecordNotFound)

		m.idemKey.On("Save", ctx, &ik).
			Once().
			Return(nil)

		idk, err := uc.getIdempotencyKey(ctx, &ik)

		assert.NoError(t, err)
		assert.Equal(t, idempotency.RecoveryPointStarted, idk.RecoveryPoint)
		assert.GreaterOrEqual(t, time.Now().UTC(), idk.LastRunAt)
		if assert.NotNil(t, idk.LockedAt) {
			assert.GreaterOrEqual(t, time.Now().UTC(), *idk.LockedAt)
		}
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("Request parameters mismatch", func(t *testing.T) {
		gofakeit.Seed(time.Now().UnixNano())
		jsonRide2, err := json.Marshal(entity.Ride{
			OriginLat: gofakeit.Float64(),
			OriginLon: gofakeit.Float64(),
			TargetLat: gofakeit.Float64(),
			TargetLon: gofakeit.Float64(),
		})
		require.NoError(t, err)

		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
		}

		retIK := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide2,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(retIK, nil)

		_, err = uc.getIdempotencyKey(ctx, &ik)

		assert.Equal(t, entity.ErrIdemKeyParamsMismatch, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Request in progress", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
		}

		now := time.Now().UTC()
		retIK := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			LockedAt:       &now,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(retIK, nil)

		_, err := uc.getIdempotencyKey(ctx, &ik)

		assert.Equal(t, entity.ErrIdemKeyRequestInProgress, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		// list non-terminal Recovery Point (i.e., all but 'FINISHED')
		rps := []idempotency.RecoveryPoint{
			idempotency.RecoveryPointStarted,
			idempotency.RecoveryPointCreated,
			idempotency.RecoveryPointCharged,
		}
		// randomly pick a Recovery Point
		ix := gofakeit.Number(0, len(rps)-1)

		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  rps[ix],
		}

		retErr := errors.New("err UpdateIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		m.idemKey.On("Update", ctx, mock.Anything).
			Once().
			Return(retErr)

		_, err := uc.getIdempotencyKey(ctx, &ik)

		assert.Equal(t, retErr, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Success on UpdateIdempotencyKey", func(t *testing.T) {
		// list non-terminal Recovery Point (i.e., all but 'FINISHED')
		rps := []idempotency.RecoveryPoint{
			idempotency.RecoveryPointStarted,
			idempotency.RecoveryPointCreated,
			idempotency.RecoveryPointCharged,
		}
		// randomly pick a Recovery Point
		ix := gofakeit.Number(0, len(rps)-1)

		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  rps[ix],
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		m.idemKey.On("Update", ctx, mock.Anything).
			Once().
			Return(nil)

		_, err := uc.getIdempotencyKey(ctx, &ik)

		assert.NoError(t, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("No-op on RecoveryPointFinished", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  idempotency.RecoveryPointFinished,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		_, err := uc.getIdempotencyKey(ctx, &ik)

		assert.NoError(t, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})
}

func TestCreateRide(t *testing.T) {
	oip := &originip.OriginIP{IP: gofakeit.IPv4Address()}
	ctx := originip.NewContext(context.Background(), oip)

	mockCfg := config.Config{IdemKeyTimeout: 5}

	t.Run("Error on CreateRide", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err CreateRide")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("Save", ctx, rd).
			Once().
			Return(retErr)

		err := uc.createRide(ctx, &ik, rd)

		assert.Equal(t, retErr, err)
		m.ride.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("Error on CreateAuditRecord", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err CreateAuditRecord")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("Save", ctx, rd).
			Once().
			Return(nil)

		m.audit.On("Save", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(retErr)

		err := uc.createRide(ctx, &ik, rd)

		assert.Equal(t, retErr, err)
		m.ride.AssertNumberOfCalls(t, "Save", 1)
		m.audit.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err UpdateIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("Save", ctx, rd).
			Once().
			Return(nil)

		m.audit.On("Save", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(nil)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(retErr)

		err := uc.createRide(ctx, &ik, rd)

		assert.Equal(t, retErr, err)
		m.ride.AssertNumberOfCalls(t, "Save", 1)
		m.audit.AssertNumberOfCalls(t, "Save", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Success on createRide", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("Save", ctx, rd).
			Once().
			Return(nil)

		m.audit.On("Save", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(nil)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(nil)

		err := uc.createRide(ctx, &ik, rd)

		assert.NoError(t, err)
		assert.Equal(t, idempotency.RecoveryPointCreated, ik.RecoveryPoint)
		m.ride.AssertNumberOfCalls(t, "Save", 1)
		m.audit.AssertNumberOfCalls(t, "Save", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestCreateCharge(t *testing.T) {
	defer gock.Off()
	ctx := context.Background()

	mockCfg := config.Config{IdemKeyTimeout: 5}

	t.Run("Error on GetRideByIdempotencyKeyID", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		retErr := errors.New("err GetRideByIdempotencyKeyID")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(entity.Ride{}, retErr)

		err := uc.createCharge(ctx, &ik, nil)

		assert.Equal(t, retErr, err)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Stripe card error", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(entity.Ride{}, nil)

		gock.New(stripeURL).
			Post("/v1/charges").
			Reply(402).
			BodyString(`{
				"error": {
					"type":"card_error",
					"code": "balance_insufficient",
					"message":"card is suspicious"
				}
			}`)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(nil)

		err := uc.createCharge(ctx, &ik, nil)

		assert.Equal(t, entity.ErrPaymentProvider, err)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Stripe generic error", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(entity.Ride{}, nil)

		gock.New(stripeURL).
			Post("/v1/charges").
			Reply(503).
			BodyString(`{
				"error": {
					"type":"api_error",
					"message":"system is down"
				}
			}`)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(nil)

		err := uc.createCharge(ctx, &ik, nil)

		assert.Equal(t, entity.ErrPaymentProviderGeneric, err)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Stripe unknown error", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(entity.Ride{}, nil)

		gock.New(stripeURL).
			Post("/v1/charges").
			ReplyError(errors.New("unknown error"))

		err := uc.createCharge(ctx, &ik, nil)

		assert.Error(t, err)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Error on UpdateRide", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
		}

		rd := entity.Ride{StripeChargeID: new(string)}

		retErr := errors.New("err UpdateRide")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(rd, nil)

		gock.New(stripeURL).
			Post("/v1/charges").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		m.ride.On("Update", ctx, &rd).
			Once().
			Return(retErr)

		err := uc.createCharge(ctx, &ik, nil)

		assert.Equal(t, retErr, err)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
		m.ride.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
		}

		rd := entity.Ride{StripeChargeID: new(string)}

		retErr := errors.New("err UpdateIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(rd, nil)

		gock.New(stripeURL).
			Post("/v1/charges").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		m.ride.On("Update", ctx, &rd).
			Once().
			Return(nil)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(retErr)

		err := uc.createCharge(ctx, &ik, nil)

		assert.Equal(t, retErr, err)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
		m.ride.AssertNumberOfCalls(t, "Update", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Success on createCharge", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
		}

		rd := entity.Ride{StripeChargeID: new(string)}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.ride.On("FindOne", ctx, mock.Anything).
			Once().
			Return(rd, nil)

		gock.New(stripeURL).
			Post("/v1/charges").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		m.ride.On("Update", ctx, &rd).
			Once().
			Return(nil)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(nil)

		err := uc.createCharge(ctx, &ik, nil)

		assert.NoError(t, err)
		assert.Equal(t, idempotency.RecoveryPointCharged, ik.RecoveryPoint)
		m.ride.AssertNumberOfCalls(t, "FindOne", 1)
		m.ride.AssertNumberOfCalls(t, "Update", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestSendReceipt(t *testing.T) {
	ctx := context.Background()

	mockCfg := config.Config{IdemKeyTimeout: 5}

	t.Run("Error on CreateStagedJob", func(t *testing.T) {
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			UserID: userID,
		}

		retErr := errors.New("err CreateStagedJob")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.job.On("Save", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(retErr)

		err := uc.sendReceipt(ctx, &ik)

		assert.Equal(t, retErr, err)
		m.job.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			UserID: userID,
		}

		retErr := errors.New("err UpdateIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.job.On("Save", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(nil)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(retErr)

		err := uc.sendReceipt(ctx, &ik)

		assert.Error(t, retErr, err)
		m.job.AssertNumberOfCalls(t, "Save", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Success on CreateStagedJob", func(t *testing.T) {
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			UserID: userID,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.job.On("Save", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(nil)

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(nil)

		err := uc.sendReceipt(ctx, &ik)

		assert.NoError(t, err)
		assert.Equal(t, idempotency.RecoveryPointFinished, ik.RecoveryPoint)
		assert.Nil(t, ik.LockedAt)
		assert.Equal(t, idempotency.ResponseCodeOK, *ik.ResponseCode)
		assert.Equal(t, idempotency.ResponseBody{Message: "OK"}, *ik.ResponseBody)
		m.job.AssertNumberOfCalls(t, "Save", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestUnlockIdempotencyKey(t *testing.T) {
	ctx := context.Background()

	mockCfg := config.Config{IdemKeyTimeout: 5}

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		ik := entity.IdempotencyKey{}

		retErr := errors.New("err UpdateIdempotencyKey")

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(retErr)

		uc.unlockIdempotencyKey(ctx, &ik)

		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Success on UpdateIdempotencyKey", func(t *testing.T) {
		ik := entity.IdempotencyKey{}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("Update", ctx, &ik).
			Once().
			Return(nil)

		uc.unlockIdempotencyKey(ctx, &ik)

		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestCreate(t *testing.T) {
	defer gock.Off()
	ctx := context.Background()

	mockCfg := config.Config{IdemKeyTimeout: 5}

	jsonRide, err := json.Marshal(entity.Ride{
		OriginLat: gofakeit.Float64(),
		OriginLon: gofakeit.Float64(),
		TargetLat: gofakeit.Float64(),
		TargetLon: gofakeit.Float64(),
	})
	require.NoError(t, err)

	t.Run("Error on createRide", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  idempotency.RecoveryPointStarted,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).
			Twice().
			Return(nil)

		// Get Idempotency Key
		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		// Create Ride
		retErr := errors.New("error createRide")
		m.store.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Return(retErr)

		err := uc.Create(ctx, &ik, &entity.Ride{})

		assert.Equal(t, retErr, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 2)
	})

	t.Run("Error on createCharge", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  idempotency.RecoveryPointCreated,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).
			Twice().
			Return(nil)

		// Get Idempotency Key
		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		// Create Charge
		retErr := errors.New("error createCharge")
		m.store.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Return(retErr)

		err := uc.Create(ctx, &ik, &entity.Ride{})

		assert.Equal(t, retErr, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Error on sendReceipt", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  idempotency.RecoveryPointCharged,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).
			Twice().
			Return(nil)

		// Get Idempotency Key
		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		// Send Receipt
		retErr := errors.New("error sendReceipt")
		m.store.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Return(retErr)

		err := uc.Create(ctx, &ik, &entity.Ride{})

		assert.Equal(t, retErr, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("No-op on finished recovery point", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  idempotency.RecoveryPointFinished,
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		// Get Idempotency Key
		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		err := uc.Create(ctx, &ik, &entity.Ride{})

		assert.NoError(t, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
	})

	t.Run("Error on unknown recovery point", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  jsonRide,
			RecoveryPoint:  "unknown",
		}

		m := getMocks()
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		m.idemKey.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).
			Once().
			Return(nil)

		err := uc.Create(ctx, &ik, &entity.Ride{})

		assert.Equal(t, entity.ErrIdemKeyUnknownRecoveryPoint, err)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.idemKey.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("Success on Create", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		user := &entity.User{
			ID:               gofakeit.Int64(),
			Email:            gofakeit.Email(),
			StripeCustomerID: gofakeit.UUID(),
		}
		ik := entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         user.ID,
			User:           user,
			RequestParams:  jsonRide,
			RecoveryPoint:  idempotency.RecoveryPointStarted,
		}

		rd := &entity.Ride{StripeChargeID: new(string)}

		m := getMocksWithTimes(0)
		uc := rideUC{cfg: mockCfg, store: m.store, ikStore: m.idemKey}

		m.idemKey.On("Update", ctx, mock.AnythingOfType("*entity.IdempotencyKey")).
			Times(4).
			Return(nil)

		// Get Idempotency Key
		m.idemKey.On("FindOne", ctx, mock.Anything, mock.Anything).
			Once().
			Return(ik, nil)

		// Create Ride
		m.ride.On("Save", ctx, rd).
			Once().
			Return(nil)

		m.audit.On("Save", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(nil)

		// Create Charge
		gock.New(stripeURL).
			Post("/v1/charges").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		m.ride.On("Update", ctx, rd).
			Once().
			Return(nil)

		// Send Receipt
		m.job.On("Save", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(nil)

		err := uc.Create(ctx, &ik, rd)

		assert.NoError(t, err)
		m.idemKey.AssertNumberOfCalls(t, "Update", 4)
		m.idemKey.AssertNumberOfCalls(t, "FindOne", 1)
		m.ride.AssertNumberOfCalls(t, "Save", 1)
		m.audit.AssertNumberOfCalls(t, "Save", 1)
		m.ride.AssertNumberOfCalls(t, "Update", 1)
		m.job.AssertNumberOfCalls(t, "Save", 1)
	})
}
