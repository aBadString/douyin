package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	Id            int
	Username      string
	FollowCount   int
	FollowerCount int
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

type UsernamePassword struct {
	Id       int
	Username string
	Password string
	Salt     string
}

func GetUsernamePasswordByUsername(username string) UsernamePassword {
	var user UsernamePassword
	ORM.Select("id, username, password, salt").
		Where("username = ?", username).
		First(&user)
	return user
}

func InsertUser(username, password, salt string) int {
	tx := ORM.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return 0
	}

	user := User{
		Username: username,
	}
	tx.Create(&user)
	if tx.Error != nil || user.Id == 0 {
		tx.Rollback()
		return 0
	}

	usernamePassword := UsernamePassword{
		Id:       user.Id,
		Username: username,
		Password: password,
		Salt:     salt,
	}
	tx.Create(&usernamePassword)
	if tx.Error != nil {
		tx.Rollback()
		return 0
	}

	tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		return 0
	}

	return user.Id
}
