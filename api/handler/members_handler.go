package handler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aneri/new_chat/api/dal"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
	"strconv"
	"time"
)

func (r *mutationResolver) NewChatRoomMembers(ctx context.Context, input model.NewChatRoomMembers) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	// Fetch chatRoomType
	chatRoomType, err := checkChatRoomTypeByChatID(ctx, input.ChatRoomID)
	if err != nil && err == sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	if chatRoomType == "GROUP" {
		for _, memberID := range input.MemberIDs {
			_, err := crConn.Db.Exec("INSERT INTO members (chatroom_id, member_id, joined_at) VALUES ($1, $2, $3) ON CONFLICT (chatroom_id, member_id) DO NOTHING RETURNING id, joined_at", input.ChatRoomID, memberID, time.Now().UTC())
			if err != nil {
				er.DebugPrintf(err)
				return false, er.InternalServerError
			}
			msg := strconv.Itoa(memberID) + " have added in this group"
			_, err = crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", input.ChatRoomID, memberID, msg, model.MessageTypeText, model.StateAdd, time.Now().UTC())
			if err != nil {
				er.DebugPrintf(err)
				return false, er.InternalServerError
			}
		}
	}
	err = fetchMemberIDsAndUpdateCharoomList(ctx, crConn, input.ChatRoomID)
	if err != nil {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return true, nil
}

//func (r *subscriptionResolver) AddNewMemberInChatRoom(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
//	panic("Not implemented")
//}

// Leave chatroom only for group chat
func (r *mutationResolver) LeaveChatRoom(ctx context.Context, input model.LeaveChatRoom) (model.ChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	chatRoomType, err := checkChatRoomTypeByChatID(ctx, input.ChatRoomID)
	if err != nil && err == sql.ErrNoRows {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	if chatRoomType == "GROUP" {
		isCreator, err := checkMemberIsCreator(ctx, input.ChatRoomID, input.MemberID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatRoom{}, er.InternalServerError
		}
		if isCreator {
			_, err = leaveMemberFromChatRoom(crConn, input.ChatRoomID, input.MemberID)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatRoom{}, er.InternalServerError
			}
		} else {
			updateCreator, err := UpdateCreatorOfChatRoom(crConn, model.LeaveChatRoom{})
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatRoom{}, er.InternalServerError
			}
			_, err = leaveMemberFromChatRoom(crConn, input.ChatRoomID, updateCreator)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatRoom{}, er.InternalServerError
			}
		}
	}
	err = fetchMemberIDsAndUpdateCharoomList(ctx, crConn, input.ChatRoomID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, er.InternalServerError
	}
	return model.ChatRoom{}, nil
}

//
//func (r *subscriptionResolver) ChatRoomLeave(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
//	panic("not implemented")
//}

func (r *memberResolver) Member(ctx context.Context, obj *model.Member) (model.User, error) {
	var memberInfo model.User
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	row, _ := crConn.Db.Query("SELECT users.id, username, first_name, last_name, email, contact, bio, profile_picture, users.created_at, users.updated_at FROM users, members WHERE users.id = $1 AND users.id = members.member_id", obj.MemberID)
	defer row.Close()
	for row.Next() {
		err := row.Scan(&memberInfo.ID, &memberInfo.Username, &memberInfo.FirstName, &memberInfo.LastName, &memberInfo.Email, &memberInfo.Contact, &memberInfo.Bio, &memberInfo.ProfilePicture, &memberInfo.CreatedAt, &memberInfo.UpdatedAt)
		if err != nil {
			er.DebugPrintf(err)
			return model.User{}, er.InternalServerError
		}
	}
	return memberInfo, nil
}

func checkChatRoomTypeByChatID(ctx context.Context, chatRoomID int) (string, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	row := crConn.Db.QueryRow("SELECT chatroom_type FROM chatrooms WHERE id = $1", chatRoomID)
	err := row.Scan(&chatroom.ChatRoomType)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return "null", er.InternalServerError
	}
	return chatroom.ChatRoomType.String(), nil
}

func checkMemberExistence(ctx context.Context, chatRoomID int, memberID int) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var isMemberExist bool
	row := crConn.Db.QueryRow("SELECT true FROM members WHERE chatroom_id = $1 and member_id = $2", chatRoomID, memberID)
	err := row.Scan(&isMemberExist)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isMemberExist, nil
}

func fetchMemberIDsAndUpdateCharoomList(ctx context.Context, crConn *dal.DbConnection, chatRoomID int) error {
	rows, err := crConn.Db.Query("SELECT member_id FROM members WHERE chatroom_id = $1", chatRoomID)
	defer rows.Close()
	for rows.Next() {
		var memberID int
		err = rows.Scan(&memberID)
		if err != nil {
			return err
		}
		fmt.Println(memberID)
		_, err = chatRoomListByMemberID(ctx, memberID)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateCreatorOfChatRoom(crConn *dal.DbConnection, room model.LeaveChatRoom) (int, error) {
	var memberID int
	row := crConn.Db.QueryRow("SELECT member_id FROM members WHERE chatroom_id = $1 AND members.member_id != $2  limit 1", room.ChatRoomID, room.MemberID)
	err := row.Scan(memberID)
	if err != nil && err != sql.ErrNoRows {
		//er.DebugPrintf(err)
		return 0, err
	}
	_, err = crConn.Db.Exec("UPDATE chatrooms SET (creator_id) = $1 WHERE id = $2", memberID, room.ChatRoomID)
	if err != nil && err != sql.ErrNoRows {
		//er.DebugPrintf(err)
		return 0, err
	}
	msg := strconv.Itoa(memberID) + " is now Group Admin"
	_, err = crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", room.ChatRoomID, memberID, msg, model.MessageTypeText, model.StateAdd, time.Now().UTC())
	if err != nil {
		er.DebugPrintf(err)
		return 0, er.InternalServerError
	}
	return memberID, nil
}

func leaveMemberFromChatRoom(crConn *dal.DbConnection, chatRoomID int, memberID int) (model.ChatRoom, error) {
	_, err := crConn.Db.Exec("DELETE FROM members WHERE chatroom_id = $1 and member_id = $2", chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, err
	}
	msg := strconv.Itoa(memberID) + " leaved from this group"
	_, err = crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", chatRoomID, memberID, msg, model.MessageTypeText, model.StateAdd, time.Now().UTC())
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatRoom{}, err
	}
	return model.ChatRoom{}, nil
}
