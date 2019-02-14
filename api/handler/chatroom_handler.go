package handler

import (
	"context"
	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"
	"github.com/aneri/new_chat/model"
	"log"
	"sort"
)

// Retrieve all Chatrooms either it is private chat or group chat
func (r *queryResolver) ChatRooms(ctx context.Context) ([]model.ChatRoom, error) {
	chatRoomData, err := ChatRoomData(ctx)
	if err != nil {
		log.Println("Error to read chatrom data", err)
	}
	return chatRoomData, nil
}

// Create New Chat room
func (r *mutationResolver) NewChatRoom(ctx context.Context, input model.NewChatRoom, receiver *int) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	if input.ChatRoomType.String() != "GROUP"{
		// Check hash exist or not
		var memberIdsArray []int
		memberIdsArray = append(memberIdsArray, input.CreatorID, *receiver)
		sort.Ints(memberIdsArray)
		hashKey := helper.HashKeycreation(memberIdsArray)
		chatRoomId, _ := CheckPrivateChatRoomExistByMemberIds(ctx, hashKey)
		if chatRoomId == 0{
			_, err := crConn.Db.Exec("INSERT INTO chatroom_test (creator_id, chatroom_type, createat, hashkey) VALUES ($1, $2, NOW(), $3)", input.CreatorID, input.ChatRoomType, hashKey)
			if err != nil{
				log.Println("Error to create new private chatroom", err)
			}
		}
		chatRoomId,_ = ChatRoomIdByHashKey(ctx, hashKey)
		chatroom = model.ChatRoom{
			ChatRoomID: chatRoomId,
			CreatorID: input.CreatorID,
			ChatRoomName: input.ChatRoomName,
			ChatRoomType: input.ChatRoomType,
		}
	}else if(input.ChatRoomType.String() == "GROUP"){
		_, err := crConn.Db.Exec("INSERT INTO chatroom_test (creator_id, chatroom_name, chatroom_type, createat) VALUES ($1, $2, $3, NOW())", input.CreatorID, input.ChatRoomName, input.ChatRoomType)
		if err != nil{
			log.Println("Error to create new group chatroom", err)
		}
		chatroomId,_ := chatRoomIdForGroupChatRoom(ctx, input.CreatorID, *input.ChatRoomName)
		chatroom = model.ChatRoom{
			ChatRoomID: chatroomId,
			CreatorID: input.CreatorID,
			ChatRoomName: input.ChatRoomName,
			ChatRoomType: input.ChatRoomType,
		}
	}
	return  chatroom, nil
}

func ChatRoomData(ctx context.Context) ([]model.ChatRoom, error) {
	var chatroom model.ChatRoom
	var chatrooms []model.ChatRoom
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	rows, _ := crConn.Db.Query("SELECT id, creator_id, chatroom_name, chatroom_type, createat, updateby, updateat, deleteat FROM chatroom_test")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&chatroom.ChatRoomID, &chatroom.CreatorID, &chatroom.ChatRoomName, &chatroom.ChatRoomType, &chatroom.CreatedAt, &chatroom.UpdateBy, &chatroom.UpdatedAt, &chatroom.DeleteAt)
		if err != nil {
			log.Println("Error to scan chat room data", err)
		}
		chatrooms = append(chatrooms, chatroom)
	}
	return chatrooms, nil
}
func (r *chatRoomResolver) Members(ctx context.Context, obj *model.ChatRoom) ([]model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	var members []model.Member
	rows, _ := crConn.Db.Query("SELECT id, chatroom_id, member_id, joinat from members_test where chatroom_id = $1", obj.ChatRoomID)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&member.ID, &member.ChatRoomID, &member.MemberID, &member.JoinAt)
		if err != nil {
			log.Println("Error to retrieve member detail as per chatroom_id", err)
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *chatRoomResolver) Creator(ctx context.Context, obj *model.ChatRoom) (model.User, error) {
	var creator model.User
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	row, _ := crConn.Db.Query("SELECT name, email, contact, profile_picture, bio, user_test.createdat from user_test, chatroom_test WHERE chatroom_test.id = $1 and chatroom_test.creator_id = user_test.id", obj.ChatRoomID)
	defer row.Close()
	for row.Next() {
		err := row.Scan(&creator.Name, &creator.Email, &creator.Contact, &creator.ProfilePicture, &creator.Bio, &creator.CreatedAt)
		if err != nil {
			log.Println("Error to scan user data as per memberid at line 24 of members_handler", err)
		}
	}
	return creator, nil
}

func CheckPrivateChatRoomExistByMemberIds(ctx context.Context, hashKey string)(int, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatRoomId int
	row := crConn.Db.QueryRow("SELECT id FROM chatroom_test WHERE hashkey = $1", hashKey)
	err := row.Scan(&chatRoomId)
	if err != nil{
		log.Println("Error to count chatroom id as per hash key at time chat room creation", err)
	}
	return chatRoomId, nil
}
func chatRoomIdForGroupChatRoom(ctx context.Context, creatorId int, chatroom string)(int, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatRoomId int
	row := crConn.Db.QueryRow("SELECT id from chatroom_test WHERE creator_id = $1 and chatroom_name = $2", creatorId, chatroom)
	err := row.Scan(&chatRoomId)
	if err != nil{
		log.Println("Error to fetch chatroom id of group chat room", err)
	}
	return chatRoomId, nil
}
func ChatRoomIdByHashKey(ctx context.Context, hashKey string)(int, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatRoomId int
	row := crConn.Db.QueryRow("SELECT id from chatroom_test WHERE hashkey = $1", hashKey)
	err := row.Scan(&chatRoomId)
	if err != nil{
		log.Println("Error to find chatroom id by hashkey", err)
	}
	return chatRoomId, nil
}