package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
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
func (r *mutationResolver) NewChatRoom(ctx context.Context, input model.NewChatRoom, receiverID *int) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	chatroom = model.ChatRoom{
		CreatorID:    input.CreatorID,
		ChatRoomName: input.ChatRoomName,
		ChatRoomType: input.ChatRoomType,
	}
	if input.ChatRoomType.String() == "PRIVATE" {
		// HashKey Creation for private chatroom -- start
		var memberIdsArray []int
		memberIdsArray = append(memberIdsArray, input.CreatorID, *receiverID)
		sort.Ints(memberIdsArray)
		hashKey := helper.HashKeycreation(memberIdsArray)
		// HashKey Creation for private chatroom -- end
		isChatRoomExist, err := CheckHashKeyExistence(ctx, hashKey)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		tx, err := crConn.Db.Begin()
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		if isChatRoomExist.ChatRoomID == 0{
			row := tx.QueryRow("INSERT INTO chatrooms (creator_id, chatroom_name, chatroom_type, created_at, hashkey) VALUES ($1, $2, $3, $4, $5) RETURNING id", input.CreatorID, input.ChatRoomName, input.ChatRoomType, time.Now(), hashKey)
			err := row.Scan(&chatroom.ChatRoomID)
			if err != nil && err == sql.ErrNoRows {
				er.DebugPrintf(err)
				tx.Rollback()
				return model.ChatRoom{}, er.InternalServerError
			}
			// Insert member in chatroom
			_, err = tx.Exec("INSERT INTO members (chatroom_id, member_id, joined_at) VALUES ($1, $2, $3), ($1, $4, $3)", chatroom.ChatRoomID, input.CreatorID, time.Now(), receiverID)
			if err != nil {
				er.DebugPrintf(err)
				tx.Rollback()
				return model.ChatRoom{}, er.InternalServerError
			}
			tx.Commit()
		}else {
			chatroom = model.ChatRoom{
				ChatRoomID: isChatRoomExist.ChatRoomID,
				ChatRoomName: input.ChatRoomName,
				ChatRoomType: input.ChatRoomType,
				CreatorID: input.CreatorID,
			}
			return chatroom, nil
		}
	} else {
		row := crConn.Db.QueryRow("INSERT INTO members (chatroom_id, member_id, joined_at) SELECT *FROM [INSERT INTO chatrooms (creator_id, chatroom_name, chatroom_type, created_at) VALUES ($1, $2, $3, $4) RETURNING id, creator_id, created_at] RETURNING chatroom_id", input.CreatorID, input.ChatRoomName, input.ChatRoomType, time.Now())
		err := row.Scan(&chatroom.ChatRoomID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
	}
	return chatroom, nil
}

// Delete chat by particular member
func (r *mutationResolver) DeleteChat(ctx context.Context, input model.DeleteChat) (model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	_, err := crConn.Db.Exec("UPDATE members SET deleted_at = $1 WHERE chatroom_id = $2 AND member_id = $3", time.Now(), input.ChatRoomID, input.MemberID)
	if err != nil {
		er.DebugPrintf(err)
		return model.Member{}, er.InternalServerError
	}
	return model.Member{}, nil
}
func (r *subscriptionResolver) ChatDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error){
	panic("not implemented")
}
// Update chatroom detail
func (r *mutationResolver) UpdateChatRoomDetail(ctx context.Context, input model.UpdateChatRoomDetail) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	chatroom = model.ChatRoom{
		ChatRoomID:   input.ChatRoomID,
		ChatRoomName: input.ChatRoomName,
		UpdateByID:   input.UpdateByID,
	}
	isMemberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, *input.UpdateByID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	if isMemberExist {
		row := crConn.Db.QueryRow("UPDATE chatrooms SET (chatroom_name, updated_by, updated_at) = ($1, $2, $3) WHERE id = $4 RETURNING updated_at", input.ChatRoomName, input.UpdateByID, time.Now(), input.ChatRoomID)
		err := row.Scan(&chatroom.UpdatedAt)
		if err != nil {
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
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	if isCreator {
		totalChatRoomMember, err := CountChatRoomMember(ctx, input.ChatRoomID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		if totalChatRoomMember <= 1 {
			_, err := crConn.Db.Exec("DELETE FROM chatrooms WHERE id = $1", input.ChatRoomID)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatRoom{}, er.InternalServerError
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
	rows, _ := crConn.Db.Query("SELECT id, chatroom_id, member_id, joined_at, deleted_at FROM members WHERE chatroom_id = $1", &obj.ChatRoomID)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&member.ID, &member.ChatRoomID, &member.MemberID, &member.JoinAt, &member.DeleteAt)
		if err != nil {
			er.DebugPrintf(err)
			return []model.Member{}, er.InternalServerError
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *chatRoomResolver) Creator(ctx context.Context, obj *model.ChatRoom) (model.User, error) {
	var creator model.User
	fmt.Println(obj.ChatRoomID)
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	row := crConn.Db.QueryRow("SELECT users.id, username, first_name, last_name, email, contact, bio, profile_picture, users.created_at, users.updated_at FROM users, chatrooms WHERE chatrooms.id = $1 AND chatrooms.creator_id = users.id", obj.ChatRoomID)
	err := row.Scan(&creator.ID, &creator.Username, &creator.FirstName, &creator.LastName, &creator.Email, &creator.Contact, &creator.Bio, &creator.ProfilePicture, &creator.CreatedAt, &creator.UpdatedAt)
	if err != nil && err != sql.ErrNoRows{
		er.DebugPrintf(err)
		return model.User{}, er.InternalServerError
	}
	return creator, nil
}

func (r *chatRoomResolver) UpdateBy(ctx context.Context, obj *model.ChatRoom) (*model.User, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var updateById int = *obj.UpdateByID
	var user model.User
	rows := crConn.Db.QueryRow("SELECT id, username, first_name, last_name, email, contact, bio, profile_picture, users.created_at FROM users WHERE id = $1", updateById)

	err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Contact, &user.Bio, &user.ProfilePicture, &user.CreatedAt)
	if err != nil && err != sql.ErrNoRows{
		er.DebugPrintf(err)
		return nil, er.InternalServerError
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
func CheckHashKeyExistence(ctx context.Context, hashKey string) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	row := crConn.Db.QueryRow("SELECT id, creator_id, chatroom_name, chatroom_type FROM chatrooms WHERE hashkey = $1", hashKey)
	err := row.Scan(&chatroom.ChatRoomID, &chatroom.CreatorID, &chatroom.ChatRoomName, &chatroom.ChatRoomType)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	return chatroom, nil
}

func CheckCreator(ctx context.Context, chatRoomID int, creatorID int) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var isCreator bool
	row := crConn.Db.QueryRow("SELECT true FROM chatrooms WHERE id = $1 AND creator_id = $2", chatRoomID, creatorID)
	err := row.Scan(&isCreator)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isCreator, nil
}

func CountChatRoomMember(ctx context.Context, chatRoomID int) (int, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var totalChatRoomMember int
	row := crConn.Db.QueryRow("SELECT count(member_id) FROM members WHERE chatroom_id = $1", chatRoomID)
	err := row.Scan(totalChatRoomMember)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return 0, er.InternalServerError
	}
	return totalChatRoomMember, nil
}
