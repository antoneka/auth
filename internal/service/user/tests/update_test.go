package tests

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"

	"github.com/antoneka/platform-common/pkg/db"
	transacMocks "github.com/antoneka/platform-common/pkg/db/mocks"
	"github.com/antoneka/platform-common/pkg/db/pg"
	"github.com/antoneka/platform-common/pkg/db/transaction"

	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/service/user"
	"github.com/antoneka/auth/internal/storage/postgres"
	storeMocks "github.com/antoneka/auth/internal/storage/postgres/mocks"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userStorageMockFunc func(mc *minimock.Controller) postgres.UserStorage
	type logStorageMockFunc func(mc *minimock.Controller) postgres.LogStorage
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id          = gofakeit.Int64()
		oldName     = gofakeit.Name()
		newName     = gofakeit.Name()
		oldEmail    = gofakeit.Email()
		newEmail    = gofakeit.Email()
		oldPassword = gofakeit.Password(true, true, true, true, true, 6)
		newPassword = gofakeit.Password(true, true, true, true, true, 6)
		oldRole     = strconv.Itoa(int(gofakeit.Int32()%2 + 1))
		newRole     = strconv.Itoa(int(gofakeit.Int32()%2 + 1))

		currentUser = &model.User{
			ID: id,
			UserInfo: model.UserInfo{
				Name:     oldName,
				Email:    oldEmail,
				Password: oldPassword,
				Role:     model.Role(oldRole),
			},
		}

		req = &model.User{
			ID: id,
			UserInfo: model.UserInfo{
				Name:     newName,
				Email:    newEmail,
				Password: newPassword,
				Role:     model.Role(newRole),
			},
		}

		storeErr = fmt.Errorf("storage error")
	)
	// t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		err             error
		txManagerMock   txManagerMockFunc
		userStorageMock userStorageMockFunc
		logStorageMock  logStorageMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				txMock := transacMocks.NewTxMock(mc)
				ctxWithTx := pg.MakeContextTx(ctx, txMock)
				txMock.CommitMock.Expect(ctxWithTx).Return(nil)

				txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
				transactorMock := transacMocks.NewTransactorMock(mc)
				transactorMock.BeginTxMock.Expect(ctx, txOpts).Return(txMock, nil)

				txManagerMock := transaction.NewTransactionManager(transactorMock)

				return txManagerMock
			},
			userStorageMock: func(mc *minimock.Controller) postgres.UserStorage {
				mock := storeMocks.NewUserStorageMock(mc)
				mock.GetMock.Expect(ctx, id).Return(currentUser, nil)

				newUser := req
				newUser.UpdatedAt = sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				}
				mock.UpdateMock.Expect(ctx, newUser).Return(nil)

				return mock
			},
			logStorageMock: func(mc *minimock.Controller) postgres.LogStorage {
				mock := storeMocks.NewLogStorageMock(mc)
				mock.LogMock.Expect(ctx, &model.LogUser{
					UserID: id,
					Action: model.LogActionUpdateUser,
				}).Return(nil)

				return mock
			},
		},
		{
			name: "storage error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: storeErr,
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				txMock := transacMocks.NewTxMock(mc)
				ctxWithTx := pg.MakeContextTx(ctx, txMock)
				txMock.RollbackMock.Expect(ctxWithTx).Return(nil)

				txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
				transactorMock := transacMocks.NewTransactorMock(mc)
				transactorMock.BeginTxMock.Expect(ctx, txOpts).Return(txMock, nil)

				txManagerMock := transaction.NewTransactionManager(transactorMock)

				return txManagerMock
			},
			userStorageMock: func(mc *minimock.Controller) postgres.UserStorage {
				mock := storeMocks.NewUserStorageMock(mc)
				mock.GetMock.Expect(ctx, id).Return(currentUser, nil)

				newUser := req
				newUser.UpdatedAt = sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				}
				mock.UpdateMock.Expect(ctx, newUser).Return(storeErr)

				return mock
			},
			logStorageMock: func(mc *minimock.Controller) postgres.LogStorage {
				mock := storeMocks.NewLogStorageMock(mc)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userStorageMock := tt.userStorageMock(mc)
			logStorageMock := tt.logStorageMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := user.NewService(
				userStorageMock,
				logStorageMock,
				txManagerMock,
			)

			err := service.Update(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
