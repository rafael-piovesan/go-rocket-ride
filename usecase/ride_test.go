//go:build unit
// +build unit

package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	rocketride "github.com/rafael-piovesan/go-rocket-ride"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/idempotency"
	"github.com/rafael-piovesan/go-rocket-ride/entity/originip"
	"github.com/rafael-piovesan/go-rocket-ride/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetIdempotencyKey(t *testing.T) {
	ctx := context.Background()

	mockDS := &mocks.Datastore{}

	var mockAtomic *mock.Call
	mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(rocketride.Datastore) error)
			if !ok {
				panic("argument mismatch")
			}

			// Call the actual func argument 'fn' passed in to
			// 'Atomic(context.Context, func(rocketride.Datastore) error) error'
			// as expected from its second parameter and, while doing so, inject the
			// mocked Datastore instance 'mockDS' so we're able to test the other calls
			// made to it inside the 'Atomic' block.
			mockAtomic.Return(fn(mockDS))
		})

	uc := rideUseCase{store: mockDS}

	t.Run("Error on GetIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
		}

		retErr := errors.New("err GetIdempotencyKey")

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(nil, retErr)

		_, err := uc.getIdempotencyKey(ctx, ik)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
	})

	t.Run("Error on CreateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
		}

		retErr := errors.New("err CreateIdempotencyKey")

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(nil, entity.ErrNotFound)

		mockDS.On("CreateIdempotencyKey", ctx, ik).
			Once().
			Return(nil, retErr)

		_, err := uc.getIdempotencyKey(ctx, ik)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "CreateIdempotencyKey", ctx, ik)
	})

	t.Run("Success on CreateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
		}

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(nil, entity.ErrNotFound)

		mockDS.On("CreateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		res, err := uc.getIdempotencyKey(ctx, ik)

		assert.NoError(t, err)
		assert.Equal(t, ik, res)
		assert.Equal(t, idempotency.RecoveryPointStarted, ik.RecoveryPoint)
		assert.True(t, time.Now().After(ik.LastRunAt))
		if assert.NotNil(t, ik.LockedAt) {
			assert.True(t, time.Now().After(*ik.LockedAt))
		}
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "CreateIdempotencyKey", ctx, ik)
	})

	t.Run("Request parameters mismatch", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("bar"),
		}

		retIK := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("foo"),
		}

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(retIK, nil)

		_, err := uc.getIdempotencyKey(ctx, ik)

		assert.Equal(t, entity.ErrIdemKeyParamsMismatch, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
	})

	t.Run("Request in progress", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("foo"),
		}

		now := time.Now().UTC()
		retIK := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("foo"),
			LockedAt:       &now,
		}

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(retIK, nil)

		_, err := uc.getIdempotencyKey(ctx, ik)

		assert.Equal(t, entity.ErrIdemKeyRequestInProgress, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
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
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("foo"),
			RecoveryPoint:  rps[ix],
		}

		retErr := errors.New("err UpdateIdempotencyKey")

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(nil, retErr)

		_, err := uc.getIdempotencyKey(ctx, ik)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
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
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("foo"),
			RecoveryPoint:  rps[ix],
		}

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		res, err := uc.getIdempotencyKey(ctx, ik)

		assert.NoError(t, err)
		assert.Equal(t, ik, res)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("No-op on RecoveryPointFinished", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RequestParams:  []byte("foo"),
			RecoveryPoint:  idempotency.RecoveryPointFinished,
		}

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		res, err := uc.getIdempotencyKey(ctx, ik)

		assert.NoError(t, err)
		assert.Equal(t, ik, res)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
	})
}

func TestCreateRide(t *testing.T) {
	oip := &originip.OriginIP{IP: gofakeit.IPv4Address()}
	ctx := originip.NewContext(context.Background(), oip)

	mockDS := &mocks.Datastore{}

	var mockAtomic *mock.Call
	mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(rocketride.Datastore) error)
			if !ok {
				panic("argument mismatch")
			}

			// Call the actual func argument 'fn' passed in to
			// 'Atomic(context.Context, func(rocketride.Datastore) error) error'
			// as expected from its second parameter and, while doing so, inject the
			// mocked Datastore instance 'mockDS' so we're able to test the other calls
			// made to it inside the 'Atomic' block.
			mockAtomic.Return(fn(mockDS))
		})

	uc := rideUseCase{store: mockDS}

	t.Run("Error on CreateRide", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err CreateRide")

		mockDS.On("CreateRide", ctx, rd).
			Once().
			Return(nil, retErr)

		_, err := uc.createRide(ctx, ik, rd)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "CreateRide", ctx, rd)
	})

	t.Run("Error on CreateAuditRecord", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err CreateAuditRecord")

		mockDS.On("CreateRide", ctx, rd).
			Once().
			Return(rd, nil)

		mockDS.On("CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(nil, retErr)

		_, err := uc.createRide(ctx, ik, rd)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "CreateRide", ctx, rd)
		mockDS.AssertCalled(t, "CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord"))
	})

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err UpdateIdempotencyKey")

		mockDS.On("CreateRide", ctx, rd).
			Once().
			Return(rd, nil)

		mockDS.On("CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(&entity.AuditRecord{}, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(nil, retErr)

		_, err := uc.createRide(ctx, ik, rd)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "CreateRide", ctx, rd)
		mockDS.AssertCalled(t, "CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord"))
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("Success on createRide", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		mockDS.On("CreateRide", ctx, rd).
			Once().
			Return(rd, nil)

		mockDS.On("CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(&entity.AuditRecord{}, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		resRd, err := uc.createRide(ctx, ik, rd)

		assert.NoError(t, err)
		assert.Equal(t, rd, resRd)
		assert.Equal(t, idempotency.RecoveryPointCreated, ik.RecoveryPoint)
		mockDS.AssertCalled(t, "CreateRide", ctx, rd)
		mockDS.AssertCalled(t, "CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord"))
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})
}

func TestCreateCharge(t *testing.T) {
	ctx := context.Background()

	mockDS := &mocks.Datastore{}

	var mockAtomic *mock.Call
	mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(rocketride.Datastore) error)
			if !ok {
				panic("argument mismatch")
			}

			// Call the actual func argument 'fn' passed in to
			// 'Atomic(context.Context, func(rocketride.Datastore) error) error'
			// as expected from its second parameter and, while doing so, inject the
			// mocked Datastore instance 'mockDS' so we're able to test the other calls
			// made to it inside the 'Atomic' block.
			mockAtomic.Return(fn(mockDS))
		})

	uc := rideUseCase{store: mockDS}

	t.Run("Error on GetRideByIdempotencyKeyID", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		retErr := errors.New("err GetRideByIdempotencyKeyID")

		mockDS.On("GetRideByIdempotencyKeyID", ctx, keyID).
			Once().
			Return(nil, retErr)

		err := uc.createCharge(ctx, ik, nil)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetRideByIdempotencyKeyID", ctx, keyID)
	})

	t.Run("Error on UpdateRide", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err UpdateRide")

		mockDS.On("GetRideByIdempotencyKeyID", ctx, keyID).
			Once().
			Return(rd, nil)

		mockDS.On("UpdateRide", ctx, rd).
			Once().
			Return(nil, retErr)

		err := uc.createCharge(ctx, ik, nil)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetRideByIdempotencyKeyID", ctx, keyID)
		mockDS.AssertCalled(t, "UpdateRide", ctx, rd)
	})

	// TODO: include tests for Stripe call

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		retErr := errors.New("err UpdateIdempotencyKey")

		mockDS.On("GetRideByIdempotencyKeyID", ctx, keyID).
			Once().
			Return(rd, nil)

		mockDS.On("UpdateRide", ctx, rd).
			Once().
			Return(rd, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(nil, retErr)

		err := uc.createCharge(ctx, ik, nil)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetRideByIdempotencyKeyID", ctx, keyID)
		mockDS.AssertCalled(t, "UpdateRide", ctx, rd)
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("Success on createCharge", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
		}

		rd := &entity.Ride{}

		mockDS.On("GetRideByIdempotencyKeyID", ctx, keyID).
			Once().
			Return(rd, nil)

		mockDS.On("UpdateRide", ctx, rd).
			Once().
			Return(rd, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		err := uc.createCharge(ctx, ik, nil)

		assert.NoError(t, err)
		assert.Equal(t, idempotency.RecoveryPointCharged, ik.RecoveryPoint)
		mockDS.AssertCalled(t, "GetRideByIdempotencyKeyID", ctx, keyID)
		mockDS.AssertCalled(t, "UpdateRide", ctx, rd)
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})
}

func TestSendReceipt(t *testing.T) {
	ctx := context.Background()

	mockDS := &mocks.Datastore{}

	var mockAtomic *mock.Call
	mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn, ok := args.Get(1).(func(rocketride.Datastore) error)
			if !ok {
				panic("argument mismatch")
			}

			// Call the actual func argument 'fn' passed in to
			// 'Atomic(context.Context, func(rocketride.Datastore) error) error'
			// as expected from its second parameter and, while doing so, inject the
			// mocked Datastore instance 'mockDS' so we're able to test the other calls
			// made to it inside the 'Atomic' block.
			mockAtomic.Return(fn(mockDS))
		})

	uc := rideUseCase{store: mockDS}

	t.Run("Error on CreateStagedJob", func(t *testing.T) {
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			UserID: userID,
		}

		retErr := errors.New("err CreateStagedJob")

		mockDS.On("CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(nil, retErr)

		err := uc.sendReceipt(ctx, ik)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob"))
	})

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			UserID: userID,
		}

		retErr := errors.New("err UpdateIdempotencyKey")

		mockDS.On("CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(&entity.StagedJob{}, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(nil, retErr)

		err := uc.sendReceipt(ctx, ik)

		assert.Error(t, retErr, err)
		mockDS.AssertCalled(t, "CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob"))
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("Success on CreateStagedJob", func(t *testing.T) {
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			UserID: userID,
		}

		sj := &entity.StagedJob{}

		mockDS.On("CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(sj, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		err := uc.sendReceipt(ctx, ik)

		assert.NoError(t, err)
		assert.Equal(t, idempotency.RecoveryPointFinished, ik.RecoveryPoint)
		assert.Equal(t, idempotency.ResponseCodeOK, *ik.ResponseCode)
		assert.Equal(t, idempotency.ResponseBody{Message: "OK"}, *ik.ResponseBody)
		mockDS.AssertCalled(t, "CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob"))
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})
}

func TestUnlockIdempotencyKey(t *testing.T) {
	ctx := context.Background()

	mockDS := &mocks.Datastore{}

	uc := rideUseCase{store: mockDS}

	t.Run("Error on UpdateIdempotencyKey", func(t *testing.T) {
		ik := &entity.IdempotencyKey{}

		retErr := errors.New("err UpdateIdempotencyKey")

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(nil, retErr)

		uc.unlockIdempotencyKey(ctx, ik)

		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("Success on UpdateIdempotencyKey", func(t *testing.T) {
		ik := &entity.IdempotencyKey{}

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		uc.unlockIdempotencyKey(ctx, ik)

		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	mockDS := &mocks.Datastore{}

	uc := NewRideUseCase(mockDS)

	t.Run("Error on createRide", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RecoveryPoint:  idempotency.RecoveryPointStarted,
		}

		rd := &entity.Ride{}

		var mockAtomic *mock.Call
		mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(rocketride.Datastore) error)
				if !ok {
					panic("argument mismatch")
				}
				mockAtomic.Return(fn(mockDS))
			})

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Twice().
			Return(ik, nil)

		// Get Idempotency Key
		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		// Create Ride
		retErr := errors.New("error createRide")
		mockDS.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Return(retErr)

		_, err := uc.Create(ctx, ik, rd)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("Error on createCharge", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RecoveryPoint:  idempotency.RecoveryPointCreated,
		}

		rd := &entity.Ride{}

		var mockAtomic *mock.Call
		mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(rocketride.Datastore) error)
				if !ok {
					panic("argument mismatch")
				}
				mockAtomic.Return(fn(mockDS))
			})

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Twice().
			Return(ik, nil)

		// Get Idempotency Key
		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		// Create Charge
		retErr := errors.New("error createCharge")
		mockDS.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Return(retErr)

		_, err := uc.Create(ctx, ik, rd)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
	})

	t.Run("Error on sendReceipt", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
			RecoveryPoint:  idempotency.RecoveryPointCharged,
		}

		rd := &entity.Ride{}

		var mockAtomic *mock.Call
		mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(rocketride.Datastore) error)
				if !ok {
					panic("argument mismatch")
				}
				mockAtomic.Return(fn(mockDS))
			})

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Twice().
			Return(ik, nil)

		// Get Idempotency Key
		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		// Send Receipt
		retErr := errors.New("error sendReceipt")
		mockDS.On("Atomic", mock.Anything, mock.Anything).
			Once().
			Return(retErr)

		_, err := uc.Create(ctx, ik, rd)

		assert.Equal(t, retErr, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
	})

	t.Run("No-op on finished recovery point", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RecoveryPoint:  idempotency.RecoveryPointFinished,
		}

		rd := &entity.Ride{}

		var mockAtomic *mock.Call
		mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(rocketride.Datastore) error)
				if !ok {
					panic("argument mismatch")
				}
				mockAtomic.Return(fn(mockDS))
			})

		// Get Idempotency Key with Recovery Point Finished
		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		res, err := uc.Create(ctx, ik, rd)

		assert.NoError(t, err)
		assert.Equal(t, ik, res)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
	})

	t.Run("Error on unknown recovery point", func(t *testing.T) {
		key := gofakeit.UUID()
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			IdempotencyKey: key,
			UserID:         userID,
			RecoveryPoint:  "unkown",
		}

		rd := &entity.Ride{}

		var mockAtomic *mock.Call
		mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(rocketride.Datastore) error)
				if !ok {
					panic("argument mismatch")
				}
				mockAtomic.Return(fn(mockDS))
			})

		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Once().
			Return(ik, nil)

		_, err := uc.Create(ctx, ik, rd)

		assert.Equal(t, entity.ErrIdemKeyUnknownRecoveryPoint, err)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "UpdateIdempotencyKey", ctx, ik)
	})

	t.Run("Success on Create", func(t *testing.T) {
		key := gofakeit.UUID()
		keyID := int64(gofakeit.Number(1, 1000))
		userID := int64(gofakeit.Number(1, 1000))
		ik := &entity.IdempotencyKey{
			ID:             keyID,
			IdempotencyKey: key,
			UserID:         userID,
			RecoveryPoint:  idempotency.RecoveryPointStarted,
		}

		rd := &entity.Ride{}

		var mockAtomic *mock.Call
		mockAtomic = mockDS.On("Atomic", mock.Anything, mock.Anything).
			Times(4).
			Run(func(args mock.Arguments) {
				fn, ok := args.Get(1).(func(rocketride.Datastore) error)
				if !ok {
					panic("argument mismatch")
				}
				mockAtomic.Return(fn(mockDS))
			})

		mockDS.On("UpdateIdempotencyKey", ctx, ik).
			Times(5).
			Return(ik, nil)

		// Get Idempotency Key
		mockDS.On("GetIdempotencyKey", ctx, key, userID).
			Once().
			Return(ik, nil)

		// Create Ride
		mockDS.On("CreateRide", ctx, rd).
			Once().
			Return(rd, nil)

		mockDS.On("CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord")).
			Once().
			Return(&entity.AuditRecord{}, nil)

		// Create Charge
		mockDS.On("UpdateRide", ctx, rd).
			Once().
			Return(rd, nil)

		// Send Receipt
		mockDS.On("CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob")).
			Once().
			Return(&entity.StagedJob{}, nil)

		res, err := uc.Create(ctx, ik, rd)

		assert.NoError(t, err)
		assert.Equal(t, ik, res)
		mockDS.AssertCalled(t, "GetIdempotencyKey", ctx, key, userID)
		mockDS.AssertCalled(t, "CreateRide", ctx, rd)
		mockDS.AssertCalled(t, "CreateAuditRecord", ctx, mock.AnythingOfType("*entity.AuditRecord"))
		mockDS.AssertCalled(t, "UpdateRide", ctx, rd)
		mockDS.AssertCalled(t, "CreateStagedJob", ctx, mock.AnythingOfType("*entity.StagedJob"))
	})
}
