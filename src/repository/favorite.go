package repository

import (
	"gorm.io/gorm"
	"time"
)

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
		Order("time desc").
		Find(&videoIds)
	return videoIds
}

func CreateFavorite(userId, videoId int) int {

	if IsFavorite(userId, videoId) {
		return 0
	}
	tx := ORM.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var r *gorm.DB
	fav := Favorite{UserId: userId, VideoId: videoId, Time: time.Now()}
	r = tx.Create(&fav)
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return 0
	}
	r = tx.Model(&Video{Id: videoId}).Update("favorite_count", gorm.Expr("favorite_count+?", 1))
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return 0
	}
	if tx.Commit().Error != nil {
		tx.Rollback()
		return 0
	}
	return fav.Id
}
func CancelFavorite(userId, videoId int) bool {
	if !IsFavorite(userId, videoId) {
		return false
	}
	tx := ORM.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	var r *gorm.DB
	r = tx.Where("user_id=? and video_id=?", userId, videoId).Delete(&Favorite{})
	if r.Error != nil || r.RowsAffected == 0 {
		tx.Rollback()
		return false
	}
	tx.Model(&Video{Id: videoId}).Update("favorite_count", gorm.Expr("favorite_count-?", 1))
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

func CountFavoriteByUserId(userId int) int {
	var n int64
	ORM.Model(&Favorite{}).Where("user_id = ?", userId).
		Count(&n)
	return int(n)
}
