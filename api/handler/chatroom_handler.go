package handler

import (
	"context"
	"database/sql"
	"sort"
	"strconv"
	"time"

	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
)


// Retrieve all Chatrooms either it is private chat or group chat
func (r *queryResolver) ChatRooms(ctx context.Context) ([]model.ChatRoom, error) {
	chatRoomData, err := chatRoomData(ctx)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatRoom{}, er.InternalServerError
	}
	return chatRoomData, nil
}

// Create New Chat room
func (r *mutationResolver) NewPrivateChatRoom(ctx context.Context, input model.NewPrivateChatRoom) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var err error

	// HashKey Creation for private chatroom -- start
	var memberIdsArray []int
	memberIdsArray = append(memberIdsArray, input.CreatorID, input.ReceiverID)
	sort.Ints(memberIdsArray)
	hashKey := helper.HashKeycreation(memberIdsArray)
	// HashKey Creation for private chatroom -- end

	chatRoom, err := checkHashKeyExistence(ctx, hashKey)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}

	if chatRoom.ChatRoomID == 0 {
		chatRoom.ChatRoomID, err = createChatRoom(crConn, input.CreatorID, nil, input.ChatRoomType, &hashKey)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		_, err = crConn.Db.Exec("INSERT INTO members (chatroom_id, member_id, joined_at) VALUES ($1, $2, $3), ($1, $4, $3) ON CONFLICT (chatroom_id, member_id) DO NOTHING", chatRoom.ChatRoomID, input.CreatorID, time.Now().UTC(), input.ReceiverID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
	}
	return chatRoom, nil
}

func (r *mutationResolver) NewGroupchatRoom(ctx context.Context, input model.NewGroupChatRoom) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	var err error
	chatroom = model.ChatRoom{
		CreatorID:    input.CreatorID,
		ChatRoomName: &input.ChatRoomName,
		ChatRoomType: input.ChatRoomType,
	}
	chatroom.ChatRoomID, err = createChatRoom(crConn, chatroom.CreatorID, chatroom.ChatRoomName, chatroom.ChatRoomType, nil)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	input.ReceiverID = append(input.ReceiverID, input.CreatorID)
	for _, memberID := range input.ReceiverID {
		_, err := crConn.Db.Exec("INSERT INTO members (chatroom_id, member_id, joined_at) VALUES ($1, $2, $3) ON CONFLICT (chatroom_id, member_id) DO NOTHING", chatroom.ChatRoomID, memberID, time.Now().UTC())
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		msg := strconv.Itoa(memberID) + " have added in this group"
		_, err = crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", chatroom.ChatRoomID, memberID, msg, model.MessageTypeText, model.StateAdd, time.Now().UTC())
		if err != nil {
			er.DebugPrintf(err)
			return  model.ChatRoom{}, er.InternalServerError
		}
		_, err = chatRoomListByMemberID(ctx, memberID)
		if err != nil {
			er.DebugPrintf(err)
			return  model.ChatRoom{}, er.InternalServerError
		}
	}
	return chatroom, nil
}

// Delete chat by particular member
//func (r *mutationResolver) DeleteChat(ctx context.Context, input model.DeleteChat) (bool, error) {
//	crConn := ctx.Value("crConn").(*dal.DbConnection)
//	_, err := crConn.Db.Exec("UPDATE members SET deleted_at = $1 WHERE chatroom_id = $2 AND member_id = $3", time.Now().UTC(), input.ChatRoomID, input.MemberID)
//	if err != nil {
//		er.DebugPrintf(err)
//		return false, er.InternalServerError
//	}
//	_, err = chatRoomListByMemberID(ctx, input.MemberID)
//	if err != nil {
//		er.DebugPrintf(err)
//		return  false, er.InternalServerError
//	}
//	return true, nil
//}
//func (r *subscriptionResolver) ChatDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
//	panic("not implemented")
//}

// Update chatroom detail
func (r *mutationResolver) UpdateChatRoomDetail(ctx context.Context, input model.UpdateChatRoomDetail) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	chatroom = model.ChatRoom{
		ChatRoomID:   input.ChatRoomID,
		ChatRoomName: input.ChatRoomName,
		UpdateByID:   input.UpdateByID,
	}
	isMemberExist, err := checkMemberExistence(ctx, input.ChatRoomID, *input.UpdateByID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	if isMemberExist {
		row := crConn.Db.QueryRow("UPDATE chatrooms SET (chatroom_name, updated_by, updated_at) = ($1, $2, $3) WHERE id = $4 RETURNING updated_at", input.ChatRoomName, input.UpdateByID, time.Now().UTC(), input.ChatRoomID)
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
	err = fetchMemberIDsAndUpdateCharoomList(ctx, crConn, input.ChatRoomID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	return chatroom, nil
}
//
//func (r *subscriptionResolver) ChatRoomDetailUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
//	panic("not implemented)
//}

// Delete group chatroom only by creator i.e. admin
func (r *mutationResolver) DeleteChatRoom(ctx context.Context, input model.DeleteChat) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	_, err := crConn.Db.Exec("UPDATE members SET deleted_at = $1 WHERE chatroom_id = $2 AND member_id = $3", time.Now().UTC(), input.ChatRoomID, input.MemberID)
	if err != nil {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	_, err = chatRoomListByMemberID(ctx, input.MemberID)
	if err != nil {
		er.DebugPrintf(err)
		return  false, er.InternalServerError
	}
	return true, nil
}

//func (r *subscriptionResolver) ChatRoomDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
//	panic("not implemented")
//}

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
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	row := crConn.Db.QueryRow("SELECT users.id, username, first_name, last_name, email, contact, bio, profile_picture, users.created_at, users.updated_at FROM users, chatrooms WHERE chatrooms.id = $1 AND chatrooms.creator_id = users.id", obj.ChatRoomID)
	err := row.Scan(&creator.ID, &creator.Username, &creator.FirstName, &creator.LastName, &creator.Email, &creator.Contact, &creator.Bio, &creator.ProfilePicture, &creator.CreatedAt, &creator.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
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
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return nil, er.InternalServerError
	}
	return &user, nil
}

func chatRoomData(ctx context.Context) ([]model.ChatRoom, error) {
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
func checkHashKeyExistence(ctx context.Context, hashKey string) (model.ChatRoom, error) {
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

func checkCreator(ctx context.Context, chatRoomID int, creatorID int) (bool, error) {
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

func countChatRoomMember(ctx context.Context, chatRoomID int) (int, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var totalChatRoomMember int
	row := crConn.Db.QueryRow("SELECT count(member_id) FROM members WHERE chatroom_id = $1", chatRoomID)
	err := row.Scan(&totalChatRoomMember)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return 0, er.InternalServerError
	}
	return totalChatRoomMember, nil
}

func createChatRoom(crConn *dal.DbConnection, creator_id int, chatroom_name *string, chatroom_type model.ChatRoomType, hashkey *string) (int, error) {
	var chatRoomID int
	row := crConn.Db.QueryRow("INSERT INTO chatrooms (creator_id, chatroom_name, chatroom_type, created_at, hashkey) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (creator_id, chatroom_name) DO NOTHING RETURNING id", creator_id, chatroom_name, chatroom_type, time.Now(), hashkey)
	err := row.Scan(&chatRoomID)
	if err != nil {
		return 0, err
	}
	return chatRoomID, nil
}
