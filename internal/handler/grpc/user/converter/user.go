package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/antoneka/auth/internal/model"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// ServiceToGetResponse converts a user model from the service layer to a gRPC response.
func ServiceToGetResponse(user *model.User) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp

	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        user.ID,
		Name:      user.UserInfo.Name,
		Email:     user.UserInfo.Email,
		Role:      ServiceToGRPCRole(user.UserInfo.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// UpdateRequestToService converts a gRPC update request to a service layer user model.
func UpdateRequestToService(updateRequest *desc.UpdateRequest) *model.User {
	userInfo := model.UserInfo{
		Name:     updateRequest.GetName(),
		Email:    updateRequest.GetEmail(),
		Password: updateRequest.GetPassword(),
		Role:     GRPCToServiceRole(updateRequest.GetRole()),
	}

	return &model.User{
		ID:       updateRequest.GetId(),
		UserInfo: userInfo,
	}
}

// CreateRequestToService converts a gRPC create request to a service layer user information model.
func CreateRequestToService(createRequest *desc.CreateRequest) *model.UserInfo {
	return &model.UserInfo{
		Name:     createRequest.GetName(),
		Email:    createRequest.GetEmail(),
		Password: createRequest.GetPassword(),
		Role:     GRPCToServiceRole(createRequest.GetRole()),
	}
}

// GRPCToServiceRole converts a gRPC role to a service layer role.
func GRPCToServiceRole(role desc.Role) model.Role {
	return model.Role(desc.Role_name[int32(role)])
}

// ServiceToGRPCRole converts a service layer role to a gRPC role.
func ServiceToGRPCRole(role model.Role) desc.Role {
	return desc.Role(desc.Role_value[string(role)])
}
