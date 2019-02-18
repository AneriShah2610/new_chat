package handler

import (
	"context"
	"database/sql"
	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
	"log"
	"sort"
	"time"
)

// Retrieve all Chatrooms either it is private chat or group chat
func (r *queryResolver) ChatRooms(ctx context.Context) ([]model.ChatRoom, error) {
	chatRoomData, err := ChatRoomData(ctx)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatRoom{}, er.InternalServerError
	}
	return chatRoomData, nil
}

// Create New Chat room
func (r *mutationResolver) NewChatRoom(ctx context.Context, input model.NewChatRoom, receiver *int) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	if input.ChatRoomType.String() == "PRIVATE" {
		// Check hash exist or not
		var memberIdsArray []int
		memberIdsArray = append(memberIdsArray, input.CreatorID, *receiver)
		sort.Ints(memberIdsArray)
		hashKey := helper.HashKeycreation(memberIdsArray)
		isChatRoomExist, _ := CheckHashKeyExistence(ctx, hashKey)
		if !isChatRoomExist {
			row := crConn.Db.QueryRow("INSERT INTO chatrooms (creator_id, chatroom_name, chatroom_type, created_at, hashkey) VALUES ($1, $2, $3, $4, $5) RETURNING id", input.CreatorID, input.ChatRoomName, input.ChatRoomType, time.Now(), hashKey)
		err := row.Scan(&chatroom.ChatRoomID)
			if err != nil{
				er.DebugPrintf(err)
				return model.ChatRoom{}, er.InternalServerError
			}
		}
		chatroom = model.ChatRoom{
			CreatorID:    input.CreatorID,
			ChatRoomName: input.ChatRoomName,
			ChatRoomType: input.ChatRoomType,
		}
	} else {
		row := crConn.Db.QueryRow("INSERT INTO chatrooms (creator_id, chatroom_name, chatroom_type, created_at) VALUES ($1, $2, $3, $4) RETURNING id", input.CreatorID, input.ChatRoomName, input.ChatRoomType, time.Now())
		err := row.Scan(&chatroom.ChatRoomID)
		if err != nil{
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		chatroom = model.ChatRoom{
			CreatorID:    input.CreatorID,
			ChatRoomName: input.ChatRoomName,
			ChatRoomType: input.ChatRoomType,
		}
	}
	return chatroom, nil
}

// Delete chat by particular member
func (r *mutationResolver) DeleteChat(ctx context.Context, input model.DeleteChat) (model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	_, err := crConn.Db.Exec("UPDATE members SET deleted_at = $1 WHERE chatroom_id = $2 AND member_id = $3", time.Now(), input.ChatRoomID, input.MemberID)
	if err != nil{
		er.DebugPrintf(err)
		return model.Member{}, er.InternalServerError
	}
	return model.Member{}, nil
}

// Update chatroom detail
func (r *mutationResolver) UpdateChatRoomDetail(ctx context.Context, input model.UpdateChatRoomDetail) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	isMemberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, *input.UpdateByID)
	if err != nil {
		log.Println("Member not exist", err)
	}
	if  isMemberExist{
		row := crConn.Db.QueryRow("UPDATE chatrooms SET (chatroom_name, updated_by, updated_at) = ($1, $2, $3) WHERE id = $4 RETURNING updated_at", input.ChatRoomName, input.UpdateByID, time.Now(), input.ChatRoomID)
		err :=  row.Scan(&chatroom.UpdatedAt)
		if err != nil{
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		chatroom = model.ChatRoom{
			ChatRoomID:   input.ChatRoomID,
			ChatRoomName: input.ChatRoomName,
			UpdateByID:   input.UpdateByID,
		}
	}
	return chatroom, nil
}

func (r *subscriptionResolver) ChatRoomDetailUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
	panic("not implemented")
}

// Delete group chatroom only by creator i.e. admin
func (r *mutationResolver) DeleteChatRoom(ctx context.Context, input model.DeleteChatRoom) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	isCreator, err := CheckCreator(ctx, input.ChatRoomID, input.CreaorID)
	if err != nil {
		log.Println("Error to check that creator exist or not")
	}
	if creatorData == input.CreaorID {
		countChatRoomMemebr, err := ChatRoomTotalMemberByChatRoomId(ctx, input.ChatRoomID)
		if err != nil {
			log.Println("Error to count chatroom member", err)
		}
		if countChatRoomMemebr == 0 {
			_, err := crConn.Db.Exec("DELETE FROM chatroom_test WHERE id = $1", input.ChatRoomID)
			if err != nil {
				log.Println("Error to delete chat room", err)
			}
		}
	}
	return chatroom, nil
}

func (r *subscriptionResolver) ChatRoomDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
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

func (r *chatRoomResolver) UpdateBy(ctx context.Context, obj *model.ChatRoom) (*model.User, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var updateById int = *obj.UpdateByID
	var user model.User
	rows, _ := crConn.Db.Query("SELECT id, name, email, contact, profile_picture, bio, createdat FROM user_test WHERE id = $1", updateById)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Contact, &user.ProfilePicture, &user.Bio, &user.CreatedAt)
		if err != nil {
			log.Println("Error to scan user details which update chat room details", err)
		}
	}
	return &user, nil
}

func ChatRoomData(ctx context.Context) ([]model.ChatRoom, error) {
	var chatroom model.ChatRoom
	var chatrooms []model.ChatRoom
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	rows, _ := crConn.Db.Query("SELECT id, creator_id, chatroom_name, chatroom_type, created_at, updated_by, updated_at, deleted_at FROM chatrooms")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&chatroom.ChatRoomID, &chatroom.CreatorID, &chatroom.ChatRoomName, &chatroom.ChatRoomType, &chatroom.CreatedAt, &chatroom.UpdateByID, &chatroom.UpdatedAt, &chatroom.DeleteAt)
		if err != nil {
			er.DebugPrintf(err)
			return []model.ChatRoom{}, er.InternalServerError
		}
		chatrooms = append(chatrooms, chatroom)
	}
	return chatrooms, nil
}
func CheckHashKeyExistence(ctx context.Context, hashKey string) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var isChatRoomExist bool
	row := crConn.Db.QueryRow("SELECT true from chatrooms WHERE hashkey = $1", hashKey)
	err := row.Scan(&isChatRoomExist)
	if err != nil && err == sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isChatRoomExist, nil
}

func CheckCreator(ctx context.Context, chatRoomID int, creatorID int) (bool, error) {
	crConn := ctx.Value("CrConn").(*dal.DbConnection)
	var isCreator bool
	row := crConn.Db.QueryRow("SELECT true FROM chatrooms WHERE id = $1 AND creator_id = $2", chatRoomID, creatorID)
	err := row.Scan(&isCreator)
	if err != nil && err == sql.ErrNoRows{
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isCreator, nil
}
// Todo: change in all file