package service

import (
	"douyin/base"
	"douyin/conf"
	"douyin/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	username := c.Query("username")

	// 1. 密码加盐哈希
	password, err := generateHashFromPassword(c.Query("password"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "注册失败",
		})
		return
	}

	// 2. 创建用户, 并返回 user_id
	userId := repository.InsertUser(username, password)
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "注册失败, 用户已经存在",
		})
		return
	}

	token, err := generateToken(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "Token 生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     userId,
		Token:      token,
	})
}

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

	userId := checkPassword(username, c.Query("password"))
	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusUnauthorized,
			"status_msg":  "用户名或密码错误",
		})
		return
	}

	token, err := generateToken(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusInternalServerError,
			"status_msg":  "Token 生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     userId,
		Token:      token,
	})
}

func GetUserIdByToken(token string) (int, error) {
	userId := verifyToken(token)
	if userId == 0 {
		return 0, base.NewError(401, "token 被篡改或者已经失效")
	}
	return userId, nil
}

// checkPassword 验证密码, 若正确返回 userId, 否则返回 0
func checkPassword(username, password string) int {
	// 1. 获取哈希后的密码, 顺便判断 username 是否存在
	u := repository.GetUsernamePasswordByUsername(username)
	if u.Id == 0 || u.Password == "" {
		return 0
	}

	// 2. 比较密码
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	)

	if err == nil {
		return u.Id
	} else {
		return 0
	}
}

// encodePassword 密码加盐后哈希
// 相同的 password 加盐后每次得到的哈希结果都不同
func generateHashFromPassword(password string) (string, error) {
	encodedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(encodedPassword), nil
}

type claims struct {
	jwt.StandardClaims
	UserId int
}

func generateToken(userId int) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    conf.Properties.Hostname,
			ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		},
		UserId: userId,
	}).SignedString([]byte(conf.Properties.SecretKey))

	if err != nil {
		return "", err
	}
	return token, nil
}

func verifyToken(token string) int {
	_token, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Properties.SecretKey), nil
	})
	if err != nil {
		return 0
	}

	_claims, ok := _token.Claims.(*claims)
	if !ok || !_token.Valid {
		return 0
	}

	return _claims.UserId
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
