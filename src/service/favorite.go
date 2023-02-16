package service

import (
	"douyin/base"
	"douyin/repository"
)

type FavoriteListRequest struct {
	UserId        int `query:"user_id"`
	CurrentUserId int `context:"current_user_id"`
}

type FavoriteActionRequest struct {
	CurrentUserId int `context:"current_user_id"`
	VideoId       int `query:"video_id"`
	ActionType    int `query:"action_type"`
}

func FavoriteList(request FavoriteListRequest) VideoList {
	videoIds := repository.GetVideoIdsByUserId(request.UserId)
	videos := repository.GetVideoListIn(videoIds)
	return toVideoList(request.CurrentUserId, videos)
}

func FavoriteAction(request FavoriteActionRequest) error {
	if request.CurrentUserId == 0 {
		return base.NewUnauthorizedError()
	}
	switch request.ActionType {
	case 1:
		if repository.CreateFavorite(request.CurrentUserId, request.VideoId) == 0 {
			return base.NewServerError("点赞失败")
		}
	case 2:
		if !repository.CancelFavorite(request.CurrentUserId, request.VideoId) {
			return base.NewServerError("取消点赞失败")
		}
	default:
		return base.NewServerError("非法的action_type")
	}
	return nil
}
