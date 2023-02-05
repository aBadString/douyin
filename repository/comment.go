package repository

import "time"

type Comment struct {
	Id          int
	UserId      int
	VideoId     int
	CommentText string
	Time        time.Time
}

type CommentWithUser struct {
	Id            int
	CommentText   string
	Time          time.Time
	UserId        int
	Username      string
	FollowCount   int
	FollowerCount int
}

var CommentWithUserViewSql = "select " +
	"comment.id, comment_text, time, " +
	"user_id, username, follow_count, follower_count " +
	"from comment join user on user.id = comment.user_id "

func GetCommentListByVideoId(videoId int) []CommentWithUser {
	var c []CommentWithUser
	ORM.Raw(CommentWithUserViewSql+
		"where video_id = ? "+
		"order by time desc",
		videoId,
	).Scan(&c)
	return c
}

func GetCommentWithUserById(commentId int) CommentWithUser {
	var c CommentWithUser
	ORM.Raw(CommentWithUserViewSql+
		"where comment.id = ? "+
		"limit 1",
		commentId,
	).Scan(&c)
	return c
}

func GetCommentById(commentId int) Comment {
	var c Comment
	ORM.Raw("select "+
		"id, comment_text, time, user_id, video_id "+
		"from comment "+
		"where comment.id = ? "+
		"limit 1",
		commentId,
	).Scan(&c)
	return c
}

func InsertComment(comment Comment) (commentId int) {
	tx := ORM.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0
	}

	tx.Select("user_id", "video_id", "comment_text").Create(&comment)
	if tx.Error != nil {
		tx.Rollback()
		return 0
	}

	tx.Exec("update video "+
		"set comment_count = comment_count + 1 "+
		"where id = ?", comment.VideoId)
	if tx.Error != nil {
		tx.Rollback()
		return 0
	}

	tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		return 0
	}

	return comment.Id
}

func DeleteCommentById(comment Comment) bool {
	tx := ORM.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false
	}

	tx.Delete(&Comment{}, comment.Id)
	if tx.Error != nil {
		tx.Rollback()
		return false
	}

	tx.Exec("update video "+
		"set comment_count = comment_count - 1 "+
		"where id = ?", comment.VideoId)
	if tx.Error != nil {
		tx.Rollback()
		return false
	}

	tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		return false
	}

	return true
}
