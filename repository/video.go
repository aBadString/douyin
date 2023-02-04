package repository

import "time"

type Video struct {
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
}

func GetVideoListOrderTime(time time.Time, count int) []VideoWithAuthor {
	var v []VideoWithAuthor
	ORM.Raw(
		"select "+
			"time, video.id as id, title, data, cover, favorite_count, comment_count,"+
			"user.id as author_id, username, follow_count, follower_count "+
			"from video join user on video.author_id = user.id "+
			"where time < ? "+
			"order by time desc "+
			"limit ?",
		time, count,
	).Scan(&v)
	return v
}
