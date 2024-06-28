package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/antoneka/auth/internal/handler/grpc/user"
	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/service"
	serviceMocks "github.com/antoneka/auth/internal/service/mocks"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 6)

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Id:       id,
			Name:     &name,
			Email:    &email,
			Password: &password,
		}

		userModel = &model.User{
			ID: id,
			UserInfo: model.UserInfo{
				Name:     name,
				Email:    email,
				Password: password,
			},
		}
	)
	// t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, userModel).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, userModel).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			resResponse, err := api.Update(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.want, resResponse)
		})
	}
}
