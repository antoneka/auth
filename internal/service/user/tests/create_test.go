package tests

import (
	"context"
	"fmt"
	"testing"

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

func TestCreate(t *testing.T) {
	t.Parallel()

	type userStorageMockFunc func(mc *minimock.Controller) postgres.UserStorage
	type logStorageMockFunc func(mc *minimock.Controller) postgres.LogStorage
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 6)

		storeErr = fmt.Errorf("storage error")

		req = &model.UserInfo{
			Name:     name,
			Email:    email,
			Password: password,
		}
	)
	// t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            int64
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
			want: id,
			err:  nil,
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
				mock.CreateMock.Expect(ctx, req).Return(id, nil)

				return mock
			},
			logStorageMock: func(mc *minimock.Controller) postgres.LogStorage {
				mock := storeMocks.NewLogStorageMock(mc)
				mock.LogMock.Expect(ctx, &model.LogUser{
					UserID: id,
					Action: model.LogActionCreateUser,
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
			want: 0,
			err:  storeErr,
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
				mock.CreateMock.Expect(ctx, req).Return(0, storeErr)

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

			resResponse, err := service.Create(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.want, resResponse)
		})
	}
}
