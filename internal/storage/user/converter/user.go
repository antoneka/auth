package converter

import (
	"github.com/antoneka/auth/internal/model"
	modelStore "github.com/antoneka/auth/internal/storage/user/model"
)

// StorageToServiceUser ...
func StorageToServiceUser(user *modelStore.User) *model.User {
	return &model.User{
		ID:        user.ID,
		UserInfo:  StorageToServiceUserInfo(user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// StorageToServiceUserInfo ...
func StorageToServiceUserInfo(userInfo modelStore.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
		Role:     model.Role(userInfo.Role),
	}
}
