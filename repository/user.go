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

func UpdateUserCount(currentUserId, toUserId, mode int) error {

	var err error
	return ORM.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if mode == 1 {
			err = tx.Model(&User{Id: currentUserId}).
				Update("follow_count", gorm.Expr("follow_count+?", 1)).Error
			if err != nil {
				return err
			}
			err = tx.Model(&User{Id: toUserId}).
				Update("follower_count", gorm.Expr("follower_count+?", 1)).Error
			if err != nil {
				return err
			}
			return nil

		} else {
			err = tx.Model(&User{Id: currentUserId}).
				Update("follow_count", gorm.Expr("follow_count-?", 1)).Error
			if err != nil {
				return err
			}
			err = tx.Model(&User{Id: toUserId}).
				Update("follower_count", gorm.Expr("follower_count-?", 1)).Error
			if err != nil {
				return err
			}
			return nil
		}
	})
}
