package service

import (
	"douyin/base"
	"douyin/conf"
	"douyin/redis"
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

type FavoriteOperator interface {
	IsFavorite(userId, videoId int) bool
	GetVideoIdsByUserId(userId int) []int
	CreateFavorite(userId, videoId int) bool
	CancelFavorite(userId, videoId int) bool
	CountFavoriteByUserId(userId int) int
}

var dbFavoriteOperator = &repository.DbFavoriteOperator{}
var redisFavoriteOperator = &redis.RedisFavoriteOperator{}

func FavoriteList(request FavoriteListRequest) VideoList {
	videoIdsFromRedis := redisFavoriteOperator.GetVideoIdsByUserId(request.UserId)
	videoIdsFromDb := dbFavoriteOperator.GetVideoIdsByUserId(request.UserId)
	videos := repository.GetVideoListIn(append(videoIdsFromDb, videoIdsFromRedis...))
	return videoWithAuthorToVideoList(request.CurrentUserId, videos)
}

func videoWithAuthorToVideoList(currentUserId int, videos []repository.VideoWithAuthor) VideoList {
	var videoList = make(VideoList, len(videos))
	for i, video := range videos {
		// 1. 当前用户是否关注了该视频的作者, 是否点赞了该视频
		isFavorite := false
		if currentUserId != 0 {
			isFavorite = IsFavorite(currentUserId, video.Id)
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
			// 幂等
			if IsFavorite(r.CurrentUserId, r.VideoId) {
				return nil
			}
			if !dbFavoriteOperator.CreateFavorite(r.CurrentUserId, r.VideoId) {
				return base.NewServerError("点赞失败")
			}
			return nil
		})

	case 2:
		return singleflight.DefaultGroup.Do(strconv.Itoa(r.CurrentUserId)+"cancel_favorite"+strconv.Itoa(r.VideoId), func() error {
			// 幂等
			if !IsFavorite(r.CurrentUserId, r.VideoId) {
				return nil
			}
			if !dbFavoriteOperator.CancelFavorite(r.CurrentUserId, r.VideoId) {
				return base.NewServerError("取消点赞失败")
			}
			return nil
		})
	default:
		return base.NewServerError("非法的action_type")
	}
}

func IsFavorite(userId, videoId int) bool {
	if redisFavoriteOperator.IsFavorite(userId, videoId) {
		return true
	}
	return dbFavoriteOperator.IsFavorite(userId, videoId)
}

// CountFavoriteByUserId 某个用户点赞过的视频的数量
func CountFavoriteByUserId(userId int) int {
	// 先从 redis 取
	favoriteCount := redisFavoriteOperator.CountFavoriteByUserId(userId)
	if favoriteCount != 0 {
		return favoriteCount
	}

	// 没有再从数据库取, 并更新到 Redis
	favoriteCount = dbFavoriteOperator.CountFavoriteByUserId(userId)
	redisFavoriteOperator.SetFavoriteCountWithUserId(userId, favoriteCount)
	return favoriteCount
}
