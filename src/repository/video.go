package repository

import "time"

type Video struct {
	Time          time.Time
	Id            int
	AuthorId      int
	Title         string
	Data          string
	Cover         string
	FavoriteCount int
	CommentCount  int
}

// GetVideoListByAuthorId 获取某个用户投稿的所有视频
func GetVideoListByAuthorId(authorId int) []Video {
	var videos []Video
	ORM.Select("id, title, data, cover, favorite_count, comment_count").
		Where("author_id = ?", authorId).
		Order("time desc").
		Find(&videos)
	return videos
}

func InsertVideo(video Video) {
	ORM.Save(&video)
}

type VideoWithAuthor struct {
	Time          time.Time
	Id            int
	Title         string
	Data          string
	Cover         string
	FavoriteCount int
	CommentCount  int
	AuthorId      int
	Username      string
	FollowCount   int
	FollowerCount int
	Avatar        string
}

func GetVideoListOrderTime(time time.Time, count int) []Video {
	var v []Video
	ORM.Raw(
		"select "+
			"time, id, title, data, cover, favorite_count, comment_count,"+
			"author_id "+
			"from video "+
			"where time < ? "+
			"order by time desc "+
			"limit ?",
		time, count,
	).Scan(&v)
	return v
}

func GetVideoListIn(videoIds []int) []VideoWithAuthor {
	var v []VideoWithAuthor
	ORM.Raw(
		"select "+
			"video.id as id, title, data, cover, favorite_count, comment_count,"+
			"user.id as author_id, username, follow_count, follower_count, avatar "+
			"from video join user on video.author_id = user.id "+
			"where video.id in ? ",
		videoIds,
	).Scan(&v)
	return v
}

func ExistVideoById(videoId int) bool {
	var v Video
	ORM.Select("id").
		Where("id = ?", videoId).
		First(&v)
	return v.Id != 0
}

func CountVideoByAuthorId(authorId int) int {
	var n int64
	ORM.Model(&Video{}).Where("author_id = ?", authorId).
		Count(&n)
	return int(n)
}
