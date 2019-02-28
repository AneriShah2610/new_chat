package handler

import (
	"context"
	"database/sql"
	"github.com/aneri/new_chat/api/dal"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

func (r *mutationResolver) NewChatRoomMembers(ctx context.Context, input model.NewChatRoomMembers) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	// Fetch chatRoomType
	chatRoomType, err := checkChatRoomTypeByChatID(crConn, input.ChatRoomID)
	if err != nil && err == sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	if chatRoomType == "GROUP" {
		isCreator, err := checkMemberIsCreator(crConn, input.ChatRoomID, input.CreatorID)
		if err != nil {
			er.DebugPrintf(err)
			return false, er.InternalServerError
		}
		if isCreator{
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
	}
	err = fetchMemberIDsAndUpdateCharoomList(crConn, input.ChatRoomID)
	if err != nil {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return true, nil
}

// Leave chatroom only for group chat
func (r *mutationResolver) LeaveChatRoom(ctx context.Context, input model.LeaveChatRoom) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	chatRoomType, err := checkChatRoomTypeByChatID(crConn, input.ChatRoomID)
	if err != nil && err == sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	if chatRoomType == "GROUP" {
		isCreator, err := checkMemberIsCreator(crConn, input.ChatRoomID, input.MemberID)
		if err != nil {
			er.DebugPrintf(err)
			return false, er.InternalServerError
		}
		if !isCreator {
			_, err = leaveMemberFromChatRoom(crConn, input.ChatRoomID, input.MemberID)
			if err != nil {
				er.DebugPrintf(err)
				return false, er.InternalServerError
			}
		} else {
			_, err := UpdateCreatorOfChatRoom(crConn, input.ChatRoomID, input.MemberID)
			if err != nil {
				er.DebugPrintf(err)
				return false, er.InternalServerError
			}
			_, err = leaveMemberFromChatRoom(crConn, input.ChatRoomID, input.MemberID)
			if err != nil {
				er.DebugPrintf(err)
				return false, er.InternalServerError
			}
		}
	}
	_, err = chatRoomListByMemberID(crConn, input.MemberID)
	if err != nil {
		return false, er.InternalServerError
	}
	err = fetchMemberIDsAndUpdateCharoomList(crConn, input.ChatRoomID)
	if err != nil {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return true, nil
}

//
//func (r *subscriptionResolver) ChatRoomLeave(ctx context.Context, chatRoomID int) (<-chan model.ChatRoom, error) {
//	panic("not implemented")
//}
func(r *queryResolver)MemberListWhichAreNoTMembersOfChatRoom(ctx context.Context, chatRoomID int, memberID int) ([]model.User, error){

	crConn :=  ctx.Value("crConn").(*dal.DbConnection)
	existingMemberIds, err := fetchExistingMemberIdsByChatRoom(crConn, chatRoomID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.User{}, er.InternalServerError
	}
	nonExistingUserOfChatRoom, err := fetchNonExistingMemberFromChatRoom(crConn, existingMemberIds)
	if err != nil {
		er.DebugPrintf(err)
		return []model.User{}, er.InternalServerError
	}
	return nonExistingUserOfChatRoom, nil
}

//ToDO: Remove MemberBy Creator
func (r *mutationResolver)RemoveMembersFromChatRoomByCreator(ctx context.Context, input *model.RemoveMembersFromChatRoom) (model.ChatRoom, error){
	panic("not implemented")
}

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

func checkChatRoomTypeByChatID(crConn *dal.DbConnection, chatRoomID int) (string, error) {
	var chatroom model.ChatRoom
	row := crConn.Db.QueryRow("SELECT chatroom_type FROM chatrooms WHERE id = $1", chatRoomID)
	err := row.Scan(&chatroom.ChatRoomType)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return "null", er.InternalServerError
	}
	return chatroom.ChatRoomType.String(), nil
}

func checkMemberExistence(crConn *dal.DbConnection, chatRoomID int, memberID int) (bool, error) {
	var isMemberExist bool
	row := crConn.Db.QueryRow("SELECT true FROM members WHERE chatroom_id = $1 and member_id = $2", chatRoomID, memberID)
	err := row.Scan(&isMemberExist)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isMemberExist, nil
}

func fetchMemberIDsAndUpdateCharoomList(crConn *dal.DbConnection, chatRoomID int) error {
	rows, err := crConn.Db.Query("SELECT member_id FROM members WHERE chatroom_id = $1", chatRoomID)
	defer rows.Close()
	for rows.Next() {
		var memberID int
		err = rows.Scan(&memberID)
		if err != nil {
			return err
		}
		_, err = chatRoomListByMemberID(crConn, memberID)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateCreatorOfChatRoom(crConn *dal.DbConnection, chatRoomID int, memberID int) (int, error) {
	var newCreatorID int
	row := crConn.Db.QueryRow("SELECT member_id FROM members WHERE chatroom_id = $1 AND members.member_id != $2  limit 1", chatRoomID, memberID)
	err := row.Scan(&newCreatorID)
	if err != nil && err != sql.ErrNoRows {
		//er.DebugPrintf(err)
		return 0, err
	}
	_, err = crConn.Db.Exec("UPDATE chatrooms SET creator_id = $1 WHERE id = $2", newCreatorID, chatRoomID)
	if err != nil && err != sql.ErrNoRows {
		//er.DebugPrintf(err)
		return 0, err
	}
	msg := strconv.Itoa(memberID) + " is now Group Admin"
	_, err = crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", chatRoomID, newCreatorID, msg, model.MessageTypeText, model.StateAdd, time.Now().UTC())
	if err != nil {
		er.DebugPrintf(err)
		return 0, er.InternalServerError
	}
	return newCreatorID, nil
}

func leaveMemberFromChatRoom(crConn *dal.DbConnection, chatRoomID int, memberID int) (bool, error) {
	_, err := crConn.Db.Exec("DELETE FROM members WHERE chatroom_id = $1 and member_id = $2", chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return false, err
	}
	msg := strconv.Itoa(memberID) + " leaved from this group"
	_, err = crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", chatRoomID, memberID, msg, model.MessageTypeText, model.StateAdd, time.Now().UTC())
	if err != nil {
		er.DebugPrintf(err)
		return false, err
	}
	return true, nil
}

func fetchExistingMemberIdsByChatRoom(crConn *dal.DbConnection, ChatRoomID int)([]int, error){
	var existingMemberID int
	var memberIDs []int
	rows, err := crConn.Db.Query("SELECT member_id FROM members WHERE members.chatroom_id = $1", ChatRoomID)
	if err != nil {
		return nil, er.InternalServerError
	}
	defer  rows.Close()
	for rows.Next(){
		err := rows.Scan(&existingMemberID)
		if err != nil {
			return nil, er.InternalServerError
		}
		memberIDs = append(memberIDs, existingMemberID)
	}
	return memberIDs, nil
}

func fetchNonExistingMemberFromChatRoom(crConn *dal.DbConnection, existingMemberIds []int)([]model.User, error){
	var nonExistingUser model.User
	var nonExistingUsersList []model.User
	sqlQuery, arguments, err := sqlx.In("SELECT id, username, first_name, last_name, email, contact, bio, profile_picture, created_at FROM users WHERE id NOT IN(?)", existingMemberIds)
	if err != nil {
		er.DebugPrintf(err)
		return []model.User{}, er.InternalServerError
	}
	sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	rows, err :=crConn.Db.Query(sqlQuery, arguments...)
	if err != nil {
		er.DebugPrintf(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&nonExistingUser.ID, &nonExistingUser.Username, &nonExistingUser.FirstName, &nonExistingUser.LastName, &nonExistingUser.Email, &nonExistingUser.Contact, &nonExistingUser.Bio, &nonExistingUser.ProfilePicture, &nonExistingUser.CreatedAt)
		if err != nil {
			er.DebugPrintf(err)
			return nil, err
		}
		nonExistingUsersList = append(nonExistingUsersList, nonExistingUser)
	}
	return nonExistingUsersList, nil
}