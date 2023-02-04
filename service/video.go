package service

import (
	"douyin/conf"
	"douyin/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
	"time"
)

type Video struct {
	Id            int    `json:"id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}
type VideoList []Video

type FeedRequest struct {
	LatestTime    int64 `query:"latest_time"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	CurrentUserId int   `context:"current_user_id"`
}

type NextTime int64

// Feed 视频流接口
// 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func Feed(request FeedRequest) (NextTime, VideoList) {
	var _time time.Time
	if request.LatestTime <= 0 || request.LatestTime > time.Now().Unix() {
		_time = time.Now()
	} else {
		_time = time.Unix(request.LatestTime, 0)
	}

	// 1. 查询按投稿时间倒序的视频列表
	var videos = repository.GetVideoListOrderTime(_time, 2)

	videoList := toVideoList(request.CurrentUserId, videos)

	// 4. 最后一个视频的投稿时间
	var nextTime int64 = 0
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].Time.Unix()
	}

	return NextTime(nextTime), videoList
}

func toVideoList(currentUserId int, videos []repository.VideoWithAuthor) VideoList {
	var videoList = make(VideoList, len(videos))
	for i, video := range videos {
		// 1. 当前用户是否关注了该视频的作者, 是否点赞了该视频
		isFollowAuthor, isFavorite := false, false
		if currentUserId != 0 {
			isFollowAuthor = repository.IsFollow(currentUserId, video.AuthorId)
			isFavorite = repository.IsFavorite(currentUserId, video.Id)
		}

		// 2. 重构返回数据格式
		videoList[i] = Video{
			Id: video.Id,
			Author: User{
				Id:            video.AuthorId,
				Name:          video.Username,
				FollowCount:   video.FollowCount,
				FollowerCount: video.FollowerCount,
				IsFollow:      isFollowAuthor,
			},
			PlayUrl:       conf.Hostname + conf.DataUrl + video.Data,
			CoverUrl:      conf.Hostname + conf.DataUrl + video.Cover,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		}
	}
	return videoList
}

type PublishListRequest struct {
	UserId        int `query:"user_id"`
	CurrentUserId int `context:"current_user_id"`
}

// PublishList 发布列表
// 用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(request PublishListRequest) (VideoList, error) {
	// 1. 查询作者信息
	author, err := UserInfo(UserInfoRequest{
		UserId:        request.UserId,
		CurrentUserId: request.CurrentUserId,
	})
	if err != nil {
		return nil, err
	}
	// 2. 查询该作者的全部视频
	var videos = repository.GetVideoListByAuthorId(author.Id)

	// 3. 重构返回数据格式
	var videoList = make(VideoList, len(videos))
	for i, video := range videos {
		isFavorite := false
		if request.CurrentUserId != 0 {
			isFavorite = repository.IsFavorite(request.CurrentUserId, video.Id)
		}

		videoList[i] = Video{
			Id:            video.Id,
			Author:        author,
			PlayUrl:       conf.Hostname + conf.DataUrl + video.Data,
			CoverUrl:      conf.Hostname + conf.DataUrl + video.Cover,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		}
	}
	return videoList, nil
}

type PublishRequest struct {
	Title         string `form:"title"`
	CurrentUserId int    `context:"current_user_id"`
}

// Publish 投稿接口
// 登录用户选择视频上传
func Publish(c *gin.Context, request PublishRequest) error {
	if request.CurrentUserId == 0 {
		return errors.New("请先登录再投稿视频")
	}

	// 1. 保存视频文件
	file, err := c.FormFile("data")
	if err != nil {
		return errors.New("视频上传出错")
	}

	filename := time.Now().Format("20060102") + "_" +
		strconv.FormatInt(time.Now().Unix(), 10) + "_" +
		file.Filename
	err = c.SaveUploadedFile(file, filepath.Join(conf.DataPath, filename))
	if err != nil {
		return errors.New("视频上传出错")
	}

	// 2. 保存视频信息到数据库
	repository.InsertVideo(repository.Video{
		AuthorId: request.CurrentUserId,
		Title:    request.Title,
		Data:     filename,
		Cover:    "", // TODO
	})

	return nil
}
