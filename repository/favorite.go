package repository

import "time"

type Favorite struct {
	Id      int
	UserId  int
	VideoId int
	Time    time.Time
}

func IsFavorite(userId, videoId int) bool {
	var f Favorite
	ORM.Select("id").
		Where("user_id = ? and video_id = ?", userId, videoId).
		First(&f)
	return f.Id != 0
}

func GetVideoIdsByUserId(userId int) []int {
	var videoIds []int
	ORM.Model(&Favorite{}).
		Select("video_id").
		Where("user_id = ?", userId).
		Find(&videoIds)
	return videoIds
}
