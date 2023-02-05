package repository

import "time"

type Message struct {
	Id            int
	SendUserId    int
	ReceiveUserId int
	Time          time.Time
	Content       string
}

func CreateMessage(sendUserId, receiveUserId int, content string) (int, error) {
	message := Message{
		SendUserId:    sendUserId,
		ReceiveUserId: receiveUserId,
		Content:       content,
		Time:          time.Now(),
	}
	tx := ORM.Create(&message)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return message.Id, nil
}

func GetMessageListFromSIdToRId(sendUserId, receiveUserId int) ([]Message, error) {
	messageList := make([]Message, 0)
	tx := ORM.Where("send_user_id=? and receive_user_id=?", sendUserId, receiveUserId).Find(&messageList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return messageList, nil
}
