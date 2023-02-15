package repository

type User struct {
	Id            int
	Username      string
	FollowCount   int
	FollowerCount int
	Password      string
}

func GetUserById(userId int) User {
	var user User
	ORM.Select("id, username, follow_count, follower_count").
		Where("id = ?", userId).
		First(&user)
	return user
}

func GetUsernamePasswordByUsername(username string) User {
	var user User
	ORM.Select("id, username, password").
		Where("username = ?", username).
		First(&user)
	return user
}

func InsertUser(username, password string) int {
	user := User{
		Username: username,
		Password: password,
	}
	ORM.Create(&user)
	return user.Id
}
