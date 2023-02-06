package service

import (
	"douyin/base"
	"douyin/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int    `json:"user_id"`
	Token      string `json:"token"`
}

// Login 用户登录
// 通过用户名和密码进行登录，登录成功后返回用户 id 和权限 token
func Login(c *gin.Context) {
	username := c.Query("username")
	// TODO: 密码加盐后哈希
	password := ""

	// 1. 验证密码, 返回 userId
	userId := repository.GetUserIdByUsernameAndPassword(username, password)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code": http.StatusUnauthorized,
			"status_msg":  "User doesn't exist",
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     userId,
		// TODO: 生成 Token 算法
		Token: strconv.Itoa(userId),
	})
}

type UserInfoRequest struct {
	UserId        int `query:"user_id"`
	CurrentUserId int `context:"current_user_id"`
}

type User struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// UserInfo 用户信息
// 获取用户的 id、昵称，如果实现社交部分的功能，还会返回关注数和粉丝数
func UserInfo(userRequest UserInfoRequest) (User, error) {
	var user = repository.GetUserById(userRequest.UserId)
	if user.Id == 0 {
		return User{}, base.NewNotFoundError("用户不存在")
	}

	var isFollow = false
	if userRequest.CurrentUserId != 0 {
		isFollow = repository.IsFollow(userRequest.CurrentUserId, userRequest.UserId)
	}

	return User{
		Id:            user.Id,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}, nil
}

func GetUserIdByToken(token string) (int, error) {
	// TODO: 校验 Token 算法
	userId, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, base.NewError(401, "token 被篡改或者已经失效")
	}
	return int(userId), nil
}
