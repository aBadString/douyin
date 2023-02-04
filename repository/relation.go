package repository

import (
	"fmt"
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

func CreateRelation(userId, followedUserId int) (int, error) {

	if !IsFollow(userId, followedUserId) {
		r := &Relation{UserId: userId, FollowedUserId: followedUserId, Time: time.Now()}
		tx := ORM.Create(r)
		if tx.Error != nil {
			return 0, tx.Error
		}
		return r.Id, nil
	}
	return 0, fmt.Errorf("repeatedly follow user %d", userId)
}

func CancelRelation(userId, followedUserId int) error {
	if IsFollow(userId, followedUserId) {
		tx := ORM.Where("user_id = ? and followed_user_id = ?", userId, followedUserId).Delete(&Relation{})
		if tx.Error != nil {
			return tx.Error
		}
		return nil
	}
	return fmt.Errorf("has even not followed  user %d", userId)
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
