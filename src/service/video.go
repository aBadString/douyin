package service

import (
	"bytes"
	"douyin/base"
	"douyin/conf"
	"douyin/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/jpeg"
	"log"
	"os"
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
	var videos = repository.GetVideoListOrderTime(_time, 20)

	videoList := toVideoList(request.CurrentUserId, videos)

	// 4. 最后一个视频的投稿时间
	var nextTime int64 = 0
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].Time.Unix()
	}

	return NextTime(nextTime), videoList
}

func toVideoList(currentUserId int, videos []repository.Video) VideoList {
	userMap := make(map[int]User, len(videos))

	var videoList = make(VideoList, len(videos))
	for i, video := range videos {
		// 1. 当前用户是否关注了该视频的作者
		isFavorite := false
		if currentUserId != 0 {
			isFavorite = repository.IsFavorite(currentUserId, video.Id)
		}

		// 2. 获取视频作者信息
		user, exist := userMap[video.AuthorId]
		if !exist {
			user, _ = UserInfo(UserInfoRequest{
				UserId:        video.AuthorId,
				CurrentUserId: currentUserId,
			})
			userMap[user.Id] = user
		}

		// 3. 重构返回数据格式
		videoList[i] = Video{
			Id:            video.Id,
			Author:        user,
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

type PublishListRequest struct {
	UserId        int `query:"user_id"`
	CurrentUserId int `context:"current_user_id"`
}

// PublishList 发布列表
// 用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(request PublishListRequest) VideoList {
	// 1. 查询作者信息
	author, err := UserInfo(UserInfoRequest{
		UserId:        request.UserId,
		CurrentUserId: request.CurrentUserId,
	})
	if err != nil {
		return make(VideoList, 0)
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

type PublishRequest struct {
	Title         string `form:"title"`
	CurrentUserId int    `context:"current_user_id"`
}

// PublishVideo 投稿接口
// 登录用户选择视频上传
func PublishVideo(c *gin.Context, request PublishRequest) error {
	if request.CurrentUserId == 0 {
		return base.NewUnauthorizedError()
	}

	// 1. 保存视频文件
	file, err := c.FormFile("data")
	if err != nil {
		return base.NewServerError("视频上传出错")
	}

	filename := time.Now().Format("20060102") + "_" +
		strconv.FormatInt(time.Now().Unix(), 10) + "_" +
		file.Filename
	videoPath := filepath.Join(conf.Properties.DataPath, filename)
	err = c.SaveUploadedFile(file, videoPath)
	if err != nil {
		return base.NewServerError("视频上传出错")
	}

	// 2. 生成视频封面
	coverFilename := "default.jpg"
	coverPath := filepath.Join(conf.Properties.DataPath, filename+".jpg")
	hasCover := generateCover(videoPath, coverPath)
	if hasCover {
		coverFilename = filename + ".jpg"
	}

	// 3. 保存视频信息到数据库
	repository.InsertVideo(repository.InsertVideoModel{
		AuthorId: request.CurrentUserId,
		Title:    request.Title,
		Data:     filename,
		Cover:    coverFilename,
	})

	return nil
}

// generateCover 抽取视频第一帧做为封面
// https://github.com/u2takey/ffmpeg-go/blob/master/examples/readFrameAsJpeg.go
func generateCover(videoPath string, coverOutputPath string) bool {
	// 1. 使用 ffmpeg 提取指定帧作为图像文件
	imgBuf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(imgBuf, os.Stdout).
		Run()
	if err != nil {
		log.Println("生成封面图片失败", err)
		return false
	}

	// 2. 编码为图片
	img, _, err := image.Decode(imgBuf)
	if err != nil {
		return false
	}

	// 3. 保存图片
	outFile, err := os.Create(coverOutputPath)
	if err != nil {
		return false
	}
	defer func() { _ = outFile.Close() }()

	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return false
	}

	return true
}
