package repository

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
