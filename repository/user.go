package repository

import (
	"fmt"
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

func UpdateUserCount(currentUserId, toUserId, mode int) error {

	var rowsAffected int64
	var exp1, exp2 string
	if mode == 1 {
		exp1 = "follow_count+ ?"
		exp2 = "follower_count+ ?"
	} else {
		exp1 = "follow_count- ?"
		exp2 = "follower_count- ?"
	}
	return ORM.Transaction(func(tx *gorm.DB) error {
		rowsAffected = tx.Model(&User{Id: currentUserId}).
			Update("follow_count", gorm.Expr(exp1, 1)).RowsAffected
		if rowsAffected == 0 {
			return fmt.Errorf("invalid userId:%v", currentUserId)
		}
		rowsAffected = tx.Model(&User{Id: toUserId}).
			Update("follower_count", gorm.Expr(exp2, 1)).RowsAffected
		if rowsAffected == 0 {
			return fmt.Errorf("invalid userId:%v", toUserId)
		}
		return nil
	})
}
