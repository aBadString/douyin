package service

import (
	"errors"
)

type UserInfoRequest struct {
	UserId int    `query:"user_id"`
	Token  string `query:"token"`
}

type User struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func UserInfo(userRequest UserInfoRequest) (User, error) {
	if userRequest.UserId < 0 {
		return User{}, errors.New("用户不存在")
	}

	return User{
		Id:   userRequest.UserId,
		Name: "松松",
	}, nil
}
