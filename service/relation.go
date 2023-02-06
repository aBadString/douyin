package service

import (
	"douyin/base"
	"douyin/repository"
)

type ActionRequest struct {
	CurrentUserId int `context:"current_user_id"`
	ToUserId      int `query:"to_user_id"`
	ActionType    int `query:"action_type"`
}
type ListRequest struct {
	CurrentUserId int `context:"current_user_id"`
	UserId        int `query:"user_id"`
}

// userList response
type UserList []*User

func RelationAction(r ActionRequest) error {

	if r.CurrentUserId == 0 {
		return base.NewUnauthorizedError()
	}
	//判断操作类型：1 关注； 2 取关； 其他报错
	if r.ActionType == 1 {
		_, err := repository.CreateRelation(r.CurrentUserId, r.ToUserId)
		if err != nil {
			return err
		}
		return repository.UpdateUserCount(r.CurrentUserId, r.ToUserId, 1)
	}
	if r.ActionType == 2 {
		err := repository.CancelRelation(r.CurrentUserId, r.ToUserId)
		if err != nil {
			return err
		}
		return repository.UpdateUserCount(r.CurrentUserId, r.ToUserId, 2)
	}
	return base.NewBadRequestError("invalid action_type")
}

func FollowList(lr ListRequest) (UserList, error) {

	//拿到userId的关注列表
	followList, errRel := repository.GetRelationListByUserid(lr.UserId)
	if errRel != nil {
		return nil, errRel
	}

	//查出关注列表里的用户信息，并判断currentUser是否有关注
	userList := make([]*User, 0)
	for _, rela := range followList {
		user := makeUser(rela.FollowedUserId, lr.CurrentUserId)
		if user != nil {
			userList = append(userList, user)
		}
	}
	return userList, nil
}

func FollowerList(lr ListRequest) (UserList, error) {

	//拿到userId的粉丝列表
	fansList, errRel := repository.GetRelationListByFollowerUserid(lr.UserId)
	if errRel != nil {
		return nil, errRel
	}

	//查出粉丝列表里的用户信息，并判断currentUser是否有关注
	userList := make([]*User, 0)
	for _, rela := range fansList {
		user := makeUser(rela.UserId, lr.CurrentUserId)
		if user != nil {
			userList = append(userList, user)
		}
	}
	return userList, nil
}
func FriendList(lr ListRequest) (UserList, error) {

	//先拿到userId的关注列表
	followList, errFol := repository.GetRelationListByUserid(lr.UserId)
	if errFol != nil {
		return nil, errFol
	}

	//遍历userId的关注列表，查看userId关注的用户有没有也关注userId
	//如果是朋友关系，则继续判断currentUser有没有关注该朋友
	friendList := make([]*User, 0)
	for _, rela := range followList {
		if repository.IsFollow(rela.FollowedUserId, lr.UserId) {
			user := makeUser(rela.FollowedUserId, lr.CurrentUserId)
			if user != nil {
				friendList = append(friendList, user)
			}
		}
	}
	return friendList, nil
}

func makeUser(userId, currentUserId int) *User {
	user := repository.GetUserById(userId)
	if user.Id == 0 {
		return nil
	}
	return &User{
		Id:            user.Id,
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      repository.IsFollow(currentUserId, user.Id),
	}
}
