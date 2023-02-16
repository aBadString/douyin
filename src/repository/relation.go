package repository

import (
	"gorm.io/gorm"
	"time"
)

type Relation struct {
	Id             int
	UserId         int
	FollowedUserId int
	Time           time.Time
}

/*
发现前端一个问题：第一次进入关注或者粉丝界面会发起follow/list请求；点击粉丝会发起follower/list请求
这个时候你进行关注或者取关动作，“关注”he“粉丝”两个界面互不感知，
	也就是其中一个窗口的更新没有实时更新到另一个窗口；需要重新从外面进来才能更新
*/

// IsFollow 判断 userId 是否关注了 followedUserId
func IsFollow(userId, followedUserId int) bool {
	var r Relation
	ORM.Select("id").
		Where("user_id = ? and followed_user_id = ?", userId, followedUserId).
		First(&r)
	return r.Id != 0
}

func CreateRelation(userId, followedUserId int) int {
	if IsFollow(userId, followedUserId) {
		return 0
	}
	tx := ORM.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var r *gorm.DB
	f := Relation{UserId: userId, FollowedUserId: followedUserId, Time: time.Now()}
	r = tx.Create(&f)
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return 0
	}

	r = tx.Model(&User{Id: userId}).Update("follow_count", gorm.Expr("follow_count+ ?", 1))
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return 0
	}

	r = tx.Model(&User{Id: followedUserId}).Update("follower_count", gorm.Expr("follower_count+ ?", 1))
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return 0
	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		return 0
	}
	return f.Id

}

func CancelRelation(userId, followedUserId int) bool {
	if !IsFollow(userId, followedUserId) {
		return false
	}
	tx := ORM.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var r *gorm.DB
	r = tx.Where("user_id=? and followed_user_id=?", userId, followedUserId).Delete(&Relation{})
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return false
	}
	r = tx.Model(&User{Id: userId}).Update("follow_count", gorm.Expr("follow_count- ?", 1))
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	r = tx.Model(&User{Id: followedUserId}).Update("follower_count", gorm.Expr("follower_count- ?", 1))
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return false
	}
	if tx.Commit().Error != nil {
		tx.Rollback()
		return false
	}
	return true
}

func GetRelationListByUserid(userId int) ([]Relation, error) {
	relationList := make([]Relation, 0)
	tx := ORM.Where("user_id=?", userId).Find(&relationList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return relationList, nil
}
func GetRelationListByFollowerUserid(followerUserId int) ([]Relation, error) {
	relationList := make([]Relation, 0)
	tx := ORM.Where("followed_user_id=?", followerUserId).Find(&relationList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return relationList, nil
}
