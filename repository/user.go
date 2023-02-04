package repository

import (
	"gorm.io/gorm"
)

type User struct {
	Id            int
	Username      string
	Password      string
	FollowCount   int
	FollowerCount int
}

// GetUserIdByUsernameAndPassword 验证用户名和密码, 如果正确返回 user_id
func GetUserIdByUsernameAndPassword(username, password string) int {
	var user User
	ORM.Select("id").
		Where("username = ? and password = ?", username, password).
		First(&user)
	return user.Id
}

func GetUserById(userId int) User {
	var user User
	ORM.Select("id, username, follow_count, follower_count").
		Where("id = ?", userId).
		First(&user)
	return user
}

func UpdateUserCount(currentUserId, toUserId, mode int) {
	if mode == 1 {
		ORM.Model(&User{Id: currentUserId}).
			Update("follow_count", gorm.Expr("follow_count+?", 1))
		ORM.Model(&User{Id: toUserId}).
			Update("follower_count", gorm.Expr("follower_count+?", 1))
	}
	if mode == 2 {
		ORM.Model(&User{Id: currentUserId}).
			Update("follow_count", gorm.Expr("follow_count-?", 1))
		ORM.Model(&User{Id: toUserId}).
			Update("follower_count", gorm.Expr("follower_count-?", 1))
	}
}
