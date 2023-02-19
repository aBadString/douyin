package service

import (
	"douyin/base"
	"douyin/conf"
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
	return videoWithAuthorToVideoList(request.CurrentUserId, videos)
}

func videoWithAuthorToVideoList(currentUserId int, videos []repository.VideoWithAuthor) VideoList {
	var videoList = make(VideoList, len(videos))
	for i, video := range videos {
		// 1. 当前用户是否关注了该视频的作者, 是否点赞了该视频
		isFavorite := false
		if currentUserId != 0 {
			isFavorite = repository.IsFavorite(currentUserId, video.Id)
		}

		// 2. 重构返回数据格式
		videoList[i] = Video{
			Id:            video.Id,
			PlayUrl:       conf.Properties.Hostname + conf.Properties.DataUrl + video.Data,
			CoverUrl:      conf.Properties.Hostname + conf.Properties.DataUrl + video.Cover,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		}
	}
	return videoList
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
