package service

import (
	"douyin/repository"
	"errors"
)

type CommentItem struct {
	Id         int    `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type CommentList []CommentItem
type Comment *CommentItem

type CommentListRequest struct {
	VideoId       int `query:"video_id"`
	CurrentUserId int `context:"current_user_id"`
}

// VideoCommentList 评论列表
// 查看视频的所有评论，按发布时间倒序
func VideoCommentList(request CommentListRequest) CommentList {
	var comments = repository.GetCommentListByVideoId(request.VideoId)
	var commentList = make(CommentList, len(comments))
	for i, comment := range comments {
		commentList[i] = toComment(request.CurrentUserId, request.VideoId, comment)
	}
	return commentList
}

func toComment(currentUserId, videoId int, comment repository.CommentWithUser) CommentItem {
	isFollow := false
	if currentUserId != 0 {
		isFollow = repository.IsFollow(currentUserId, videoId)
	}

	return CommentItem{
		Id: comment.Id,
		User: User{
			Id:            comment.UserId,
			Name:          comment.Username,
			FollowCount:   comment.FollowCount,
			FollowerCount: comment.FollowerCount,
			IsFollow:      isFollow,
		},
		Content:    comment.CommentText,
		CreateDate: comment.Time.Format("2006-01-02"),
	}
}

type CommentActionRequest struct {
	VideoId       int    `query:"video_id"`
	CurrentUserId int    `context:"current_user_id"`
	ActionType    int    `query:"action_type"`  // 1-发布评论，2-删除评论
	CommentText   string `query:"comment_text"` // 用户填写的评论内容，在action_type=1的时候使用
	CommentId     int    `query:"comment_id"`   // 要删除的评论id，在action_type=2的时候使用
}

// CommentAction 评论操作
// 登录用户对视频进行评论
func CommentAction(request CommentActionRequest) (Comment, error) {
	if request.CurrentUserId == 0 {
		return nil, errors.New("请先登录")
	}

	switch request.ActionType {
	case 1:
		return PublishComment(request)
	case 2:
		return DeleteComment(request)
	default:
		return nil, errors.New("action_type 值错误, 1-发布评论, 2-删除评论")
	}
}

// PublishComment 发布评论
func PublishComment(request CommentActionRequest) (Comment, error) {
	commentId := repository.InsertComment(repository.Comment{
		UserId:      request.CurrentUserId,
		VideoId:     request.VideoId,
		CommentText: request.CommentText,
	})
	if commentId == 0 {
		return nil, errors.New("发表评论失败")
	}

	var commentWithUser = repository.GetCommentWithUserById(commentId)
	comment := toComment(request.CurrentUserId, request.VideoId, commentWithUser)
	return &comment, nil
}

// DeleteComment 删除评论
func DeleteComment(request CommentActionRequest) (Comment, error) {
	var comment = repository.GetCommentById(request.CommentId)
	if comment.Id == 0 {
		return nil, errors.New("评论不存在")
	}

	if request.CurrentUserId != comment.UserId {
		return nil, errors.New("这不是您的评论")
	}

	isDelete := repository.DeleteCommentById(comment)
	if !isDelete {
		return nil, errors.New("删除评论失败")
	}

	return &CommentItem{Id: comment.Id}, nil
}
