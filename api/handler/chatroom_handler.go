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
		chatRoomId, _ := ChatRoomIdByHashKey(ctx, hashKey)
		if chatRoomId == 0{
			_, err := crConn.Db.Exec("INSERT INTO chatroom_test (creator_id, chatroom_type, createat, hashkey) VALUES ($1, $2, NOW(), $3)", input.CreatorID, input.ChatRoomType, hashKey)
			if err != nil{
				log.Println("Error to create new private chatroom", err)
			}
			chatRoomId,_ = ChatRoomIdByHashKey(ctx, hashKey)
		}
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

// Delete chat by particular member
func (r *mutationResolver)DeleteChat(ctx context.Context, input model.DeleteChat) (model.Member, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	_, err := crConn.Db.Exec("UPDATE members_test SET deleteat = NOW() WHERE chatroom_id = $1 AND member_id = $2", input.ChatRoomID, input.MemberID)
	if err != nil{
		log.Println("Error to update delete chat by member", err)
	}
	return member, nil
}

// Update chatroom detail
func (r *mutationResolver)UpdateChatRoomDetail(ctx context.Context, input model.UpdateChatRoomDetail) (model.ChatRoom, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	memberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, *input.UpdateByID)
	if err != nil{
		log.Println("Member not exist", err)
	}
	if memberExist.MemberID != 0{
		_, err := crConn.Db.Exec("UPDATE chatroom_test SET (chatroom_name, updateby, updateat) = ($1, $2, now()) where id = $3", input.ChatRoomName, input.UpdateByID, input.ChatRoomID)
		if err != nil{
			log.Println("Error to update chatroom details", err)
		}
		chatroom = model.ChatRoom{
			ChatRoomID: input.ChatRoomID,
			ChatRoomName: input.ChatRoomName,
			UpdateByID: input.UpdateByID,
		}
	}
	return chatroom, nil
}

func(r *subscriptionResolver)ChatRoomDetailUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error){
	panic("not implemented")
}

// Delete group chatroom only by creator i.e. admin
func (r *mutationResolver)DeleteChatRoom(ctx context.Context, input model.DeleteChatRoom) (model.ChatRoom, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	creatorData, err := CreatorDataByChatRoomID(ctx, input.ChatRoomID)
	if err != nil{
		log.Println("Error to check that creator exist or not")
	}
	if creatorData == input.CreaorID{
		countChatRoomMemebr,err  := ChatRoomTotalMemberByChatRoomId(ctx, input.ChatRoomID)
		if err != nil{
			log.Println("Error to count chatroom member", err)
		}
		if countChatRoomMemebr == 0{
			_, err := crConn.Db.Exec("DELETE FROM chatroom_test WHERE id = $1", input.ChatRoomID)
			if err != nil{
				log.Println("Error to delete chat room", err)
			}
		}
	}
	return  chatroom, nil
}

func(r *subscriptionResolver)ChatRoomDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error){
	panic("not implemented")
}

func (r *chatRoomResolver) Members(ctx context.Context, obj *model.ChatRoom) ([]model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	var members []model.Member
	rows, _ := crConn.Db.Query("SELECT id, chatroom_id, member_id, joinat from members_test where chatroom_id = $1", &obj.ChatRoomID)
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
	row, _ := crConn.Db.Query("SELECT name, email, contact, profile_picture, bio, user_test.createdat FROM user_test, chatroom_test WHERE chatroom_test.id = $1 and chatroom_test.creator_id = user_test.id", obj.ChatRoomID)
	defer row.Close()
	for row.Next() {
		err := row.Scan(&creator.Name, &creator.Email, &creator.Contact, &creator.ProfilePicture, &creator.Bio, &creator.CreatedAt)
		if err != nil {
			log.Println("Error to scan user data as per memberid at line 24 of members_handler", err)
		}
	}
	return creator, nil
}

func (r*chatRoomResolver)UpdateBy(ctx context.Context, obj *model.ChatRoom) (*model.User, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var updateById int = *obj.UpdateByID
	var user model.User
	rows, _ := crConn.Db.Query("SELECT id, name, email, contact, profile_picture, bio, createdat FROM user_test WHERE id = $1", updateById)
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Contact, &user.ProfilePicture, &user.Bio, &user.CreatedAt)
		if err != nil{
			log.Println("Error to scan user details which update chat room details", err)
		}
	}
	return &user, nil
}

func ChatRoomData(ctx context.Context) ([]model.ChatRoom, error) {
	var chatroom model.ChatRoom
	var chatrooms []model.ChatRoom
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	rows, _ := crConn.Db.Query("SELECT id, creator_id, chatroom_name, chatroom_type, createat, updateby, updateat, deleteat FROM chatroom_test")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&chatroom.ChatRoomID, &chatroom.CreatorID, &chatroom.ChatRoomName, &chatroom.ChatRoomType, &chatroom.CreatedAt, &chatroom.UpdateByID, &chatroom.UpdatedAt, &chatroom.DeleteAt)
		if err != nil {
			log.Println("Error to scan chat room data", err)
		}
		chatrooms = append(chatrooms, chatroom)
	}
	return chatrooms, nil
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
func CreatorDataByChatRoomID(ctx context.Context, chatRoomId int)(int, error){
	crConn := ctx.Value("CrConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	row := crConn.Db.QueryRow("SELECT creator_id FROM chatroom_test WHERE id = $1", chatRoomId)
	err := row.Scan(&chatroom.CreatorID)
	if err != nil{
		log.Println("Error to find creator id", err)
	}
	return  chatroom.CreatorID, nil
}