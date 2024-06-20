package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"

	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/service/user"
	"github.com/antoneka/auth/internal/storage/postgres"
	storeMocks "github.com/antoneka/auth/internal/storage/postgres/mocks"
	"github.com/antoneka/auth/pkg/client/db"
	transacMocks "github.com/antoneka/auth/pkg/client/db/mocks"
	"github.com/antoneka/auth/pkg/client/db/pg"
	"github.com/antoneka/auth/pkg/client/db/transaction"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type userStorageMockFunc func(mc *minimock.Controller) postgres.UserStorage
	type logStorageMockFunc func(mc *minimock.Controller) postgres.LogStorage
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		storeErr = fmt.Errorf("storage error")

		req = id
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
				mock.DeleteMock.Expect(ctx, req).Return(nil)

				return mock
			},
			logStorageMock: func(mc *minimock.Controller) postgres.LogStorage {
				mock := storeMocks.NewLogStorageMock(mc)
				mock.LogMock.Expect(ctx, &model.LogUser{
					UserID: id,
					Action: model.LogActionDeleteUser,
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
				mock.DeleteMock.Expect(ctx, req).Return(storeErr)

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

			err := service.Delete(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
