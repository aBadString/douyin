package service

import (
	"douyin/base"
	"douyin/repository"
	"sort"
)

type Message struct {
	Id         int    `json:"id"`
	FromUserId int    `json:"from_user_id"`
	ToUserId   int    `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type MessageRequest struct {
	CurrentUserId int    `context:"current_user_id"`
	ToUserId      int    `query:"to_user_id"`
	ActionType    int    `query:"action_type"`
	Content       string `query:"content"`
}

type MessageListRequest struct {
	CurrentUserId int   `context:"current_user_id"`
	ToUserId      int   `query:"to_user_id"`
	PreMsgTime    int64 `query:"pre_msg_time"`
}

type MessageList []*Message

func MessageAction(msg MessageRequest) error {
	if msg.CurrentUserId == 0 {
		return base.NewUnauthorizedError()
	}
	_, err := repository.CreateMessage(msg.CurrentUserId, msg.ToUserId, msg.Content)
	return err
}

/*
app前端暂时没有成功显示
*/
func MessageChat(ml MessageListRequest) (MessageList, error) {
	if ml.CurrentUserId == 0 {
		return nil, base.NewUnauthorizedError()
	}

	msgList1, err := getSingleMessage(ml.CurrentUserId, ml.ToUserId, ml.PreMsgTime)
	if err != nil {
		return nil, err
	}
	msgList2, err := getSingleMessage(ml.ToUserId, ml.CurrentUserId, ml.PreMsgTime)
	if err != nil {
		return nil, err
	}

	var msgList = append(msgList1, msgList2...)
	sort.Slice(msgList, func(i, j int) bool {
		return msgList[i].CreateTime < msgList[j].CreateTime
	})

	return msgList, nil
}

func getSingleMessage(from, to int, preMsgTime int64) (MessageList, error) {
	msgList, err := repository.GetMessageListFromSIdToRId(from, to, preMsgTime)
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
			CreateTime: message.Time.UnixMilli(),
		}
	}
	return msgListResp, nil
}
