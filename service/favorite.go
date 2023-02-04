package service

import "douyin/repository"

type FavoriteListRequest struct {
	UserId        int `query:"user_id"`
	CurrentUserId int `context:"current_user_id"`
}

func FavoriteList(request FavoriteListRequest) VideoList {
	videoIds := repository.GetVideoIdsByUserId(request.UserId)
	videos := repository.GetVideoListIn(videoIds)
	return toVideoList(request.CurrentUserId, videos)
}
