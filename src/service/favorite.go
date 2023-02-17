package service

import (
	"douyin/base"
	"douyin/repository"
	"douyin/singleflight"
	"strconv"
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

func FavoriteAction(r FavoriteActionRequest) error {
	if r.CurrentUserId == 0 {
		return base.NewUnauthorizedError()
	}
	switch r.ActionType {
	case 1:
		return singleflight.DefaultGroup.Do(strconv.Itoa(r.CurrentUserId)+"favorite"+strconv.Itoa(r.VideoId), func() error {
			if repository.CreateFavorite(r.CurrentUserId, r.VideoId) == 0 {
				return base.NewServerError("点赞失败")
			}
			return nil
		})

	case 2:
		return singleflight.DefaultGroup.Do(strconv.Itoa(r.CurrentUserId)+"cancel_favorite"+strconv.Itoa(r.VideoId), func() error {
			if !repository.CancelFavorite(r.CurrentUserId, r.VideoId) {
				return base.NewServerError("取消点赞失败")
			}
			return nil
		})
	default:
		return base.NewServerError("非法的action_type")
	}
}
