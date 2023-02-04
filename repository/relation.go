package repository

import "time"

type Relation struct {
	Id             int
	UserId         int
	FollowedUserId int
	Time           time.Time
}

// IsFollow 判断 userId 是否关注了 followedUserId
func IsFollow(userId, followedUserId int) bool {
	var c int64
	ORM.Count(&c).
		Where("user_id = ? and followed_user_id = ?", userId, followedUserId).
		Limit(1)
	return c != 0
}
