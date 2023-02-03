package repository

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
