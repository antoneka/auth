package converter

import (
	"github.com/antoneka/auth/internal/model"
	modelStore "github.com/antoneka/auth/internal/storage/postgres/user/model"
)

// StorageToServiceUser converts a user model from the storage layer to the service layer.
func StorageToServiceUser(user *modelStore.User) *model.User {
	return &model.User{
		ID:        user.ID,
		UserInfo:  StorageToServiceUserInfo(user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// StorageToServiceUserInfo converts a user information model from the storage layer to the service layer.
func StorageToServiceUserInfo(userInfo modelStore.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
		Role:     model.Role(userInfo.Role),
	}
}

// ServiceUserToStorage converts a user model from the service layer to the storage layer.
func ServiceUserToStorage(user *model.User) modelStore.User {
	return modelStore.User{
		ID:        user.ID,
		UserInfo:  ServiceUserInfoToStorage(&user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ServiceUserInfoToStorage converts a user information model from the service layer to the storage.
func ServiceUserInfoToStorage(userInfo *model.UserInfo) modelStore.UserInfo {
	return modelStore.UserInfo{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
		Role:     string(userInfo.Role),
	}
}
