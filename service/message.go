package service

import (
	"douyin/repository"
	"errors"
	"time"
)

type Message struct {
	Id         int    `json:"id"`
	FromUserId int    `json:"from_user_id"`
	ToUserId   int    `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}
type MessageRequest struct {
	CurrentUserId int    `context:"current_user_id"`
	ToUserId      int    `query:"to_user_id"`
	ActionType    int    `query:"action_type"`
	Content       string `query:"content"`
}

type MessageListRequest struct {
	CurrentUserId int `context:"current_user_id"`
	ToUserId      int `query:"to_user_id"`
}

type MessageList []*Message

func MessageAction(msg MessageRequest) error {
	if msg.CurrentUserId == 0 {
		return errors.New("请先登录")
	}
	_, err := repository.CreateMessage(msg.CurrentUserId, msg.ToUserId, msg.Content)
	return err
}

/*
app前端暂时没有成功显示
*/
func MessageChat(ml MessageListRequest) (MessageList, error) {

	if ml.CurrentUserId == 0 {
		return nil, errors.New("请先登录")
	}
	msgList, err := repository.GetMessageListFromSIdToRId(ml.CurrentUserId, ml.ToUserId)
	if err != nil {
		return nil, err
	}
	msgListResp := make([]*Message, len(msgList))
	for i, message := range msgList {
		msgListResp[i] = &Message{
			Id:         message.Id,
			FromUserId: message.SendUserId,
			ToUserId:   message.ReceiveUserId,
			Content:    message.Content,
			CreateTime: time.Unix(message.Time.Unix(), 0).Format("2006-01-02 15:04:05"),
		}
	}
	return msgListResp, nil
}
