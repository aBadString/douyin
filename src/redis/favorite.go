package redis

import "time"

type Favorite struct {
	UserId        int
	FavoriteCount int // 点赞过的视频总数, 包括 Videos 列表和数据库中的
	Videos        []FavoriteVideo
}
type FavoriteVideo struct {
	VideoId int
	Time    time.Time
}

type RedisFavoriteOperator struct {
}

func (*RedisFavoriteOperator) IsFavorite(userId, videoId int) bool {
	return false
}

func (*RedisFavoriteOperator) GetVideoIdsByUserId(userId int) []int {
	return nil
}

func (*RedisFavoriteOperator) CreateFavorite(userId, videoId int) bool {
	return false
}
func (*RedisFavoriteOperator) CancelFavorite(userId, videoId int) bool {
	return false
}

func (*RedisFavoriteOperator) CountFavoriteByUserId(userId int) int {
	return 0
}

func (*RedisFavoriteOperator) SetFavoriteCountWithUserId(id int, count int) {

}
