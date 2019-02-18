package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/aneri/new_chat/api/dal"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
)

var addMessageResolver map[int]chan model.ChatConversation // To add message in chatconversation table
var updateMessageResolver map[int]chan model.ChatConversation

func init() {
	addMessageResolver = map[int]chan model.ChatConversation{}
	updateMessageResolver = map[int]chan model.ChatConversation{}
}

// Retrieve chat conversation by chatRoom Id
func (r *queryResolver) ChatconversationByChatRoomID(ctx context.Context, chatRoomID int, memberID int) ([]model.ChatConversation, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	var chattconversationarr []model.ChatConversation
	isMemberExist, err := CheckMemberExistence(ctx, chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isDeleteFlagUpdate, err := ChehckDeleteFlag(ctx, chatRoomID, memberID)
		if err != nil {
			er.DebugPrintf(err)
			return []model.ChatConversation{}, er.InternalServerError
		}

		if isDeleteFlagUpdate != nil {
			rows, err := crConn.Db.Query("SELECT id, chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at, updatedat FROM chatconversation LEFT JOIN members ON members.deleted_at <= chatconversation.created_at WHERE chatconversation.chatroom_id = $1 AND chatconversation.chatroom_id = members.chatroom_id GROUP BY chatconversation.chatroom_id,sender_id, message,message_type, message_parent_id, message_status, createat, updateat, chatconversation.id ORDER BY chatconversation.id", chatRoomID)
			if err != nil {
				er.DebugPrintf(err)
				return []model.ChatConversation{}, er.InternalServerError
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&chatconversation.MessageID, &chatconversation.ChatRoomID, &chatconversation.SenderID, &chatconversation.Message, &chatconversation.MessageType, &chatconversation.MessageStatus, &chatconversation.MessageParentID, &chatconversation.CreatedAt, &chatconversation.UpdatedAt)
				if err != nil {
					er.DebugPrintf(err)
					return []model.ChatConversation{}, er.InternalServerError
				}
				chattconversationarr = append(chattconversationarr, chatconversation)
			}
		} else {
			rows, err := crConn.Db.Query("SELECT id, chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at, updatedat FROM chatconversation LEFT JOIN members ON members.created_at <= chatconversation.created_at WHERE chatconversation.chatroom_id = $1 AND chatconversation.chatroom_id = members.chatroom_id GROUP BY chatconversation.chatroom_id,sender_id, message,message_type, message_parent_id, message_status, createat, updateat, chatconversation.id ORDER BY chatconversation.id", chatRoomID)
			if err != nil {
				er.DebugPrintf(err)
				return []model.ChatConversation{}, er.InternalServerError
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&chatconversation.MessageID, &chatconversation.ChatRoomID, &chatconversation.SenderID, &chatconversation.Message, &chatconversation.MessageType, &chatconversation.MessageParentID, &chatconversation.MessageStatus, &chatconversation.CreatedAt, &chatconversation.UpdatedAt)
				if err != nil {
					er.DebugPrintf(err)
					return []model.ChatConversation{}, er.InternalServerError
				}
				chattconversationarr = append(chattconversationarr, chatconversation)
			}
		}
	}
	return chattconversationarr, nil
}

func (r *queryResolver) MemberListByChatRoomID(ctx context.Context, chatRoomID int, memberID int) ([]model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	var members []model.Member
	isMemberExist, err := CheckMemberExistence(ctx, chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.Member{}, er.InternalServerError
	}
	if isMemberExist {
		rows, err := crConn.Db.Query("SELECT members.id, chatroom_id, member_id, joined_at FROM members, users WHERE members.member_id = users.id AND chatroom_id = $1 ORDER BY users.name ASC", chatRoomID)
		if err != nil {
			er.DebugPrintf(err)
			return []model.Member{}, er.InternalServerError
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&member.ID, &member.ChatRoomID, &member.MemberID, &member.JoinAt)
			if err != nil {
				er.DebugPrintf(err)
				return []model.Member{}, er.InternalServerError
			}
			members = append(members, member)
		}
	}
	return members, nil
}

func (r *queryResolver) ChatRoomListByMemberID(ctx context.Context, memberId int) ([]model.ChatRoom, error) {
	//crConn := ctx.Value("crconn").(*dal.DbConnection)
	//var chatroom model.ChatRoom
	//var chatrooms []model.ChatRoom
	//rows, err := crConn.Db.
	panic("not implemented")
}

// Create New Message
func (r *mutationResolver) NewMessage(ctx context.Context, input model.NewMessage, senderID int) (model.ChatConversation, error) {
	var chatconversation model.ChatConversation
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	isMemberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, senderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		row := crConn.Db.QueryRow("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, message_parent_id, createat) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", input.ChatRoomID, input.SenderID, input.Message, input.MessageType, input.MessageStatus, input.MessageParentID, time.Now())
		err := row.Scan(&chatconversation.MessageID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
		chatconversation = model.ChatConversation{
			ChatRoomID:      input.ChatRoomID,
			SenderID:        input.SenderID,
			Message:         input.Message,
			MessageType:     input.MessageType,
			MessageParentID: input.MessageParentID,
			MessageStatus:   input.MessageStatus,
		}
	}

	// add new chatconversation in observer
	channelMsg := addMessageResolver[input.ChatRoomID]
	if channelMsg != nil {
		channelMsg <- chatconversation
	}
	return chatconversation, nil
}

// Live updates of new messages
func (r *subscriptionResolver) MessagePost(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {

	chatevent := make(chan model.ChatConversation, 1)
	addMessageResolver[chatRoomID] = chatevent
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(addMessageResolver, chatRoomID)
		r.mu.Unlock()
	}()
	return chatevent, nil
}

func (r *mutationResolver) UpdateMessage(ctx context.Context, input *model.UpdateMessage) (model.ChatConversation, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	isMemberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, input.SenderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isMessageOwner, err := CheckMemberOwner(ctx, input.MessageID, input.SenderID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
		if isMessageOwner {
			row := crConn.Db.QueryRow("UPDATE chatconversation SET (message, updatedat) = ($1, $2) where id = $3 RETURNING message_status, message_parent_id, created_at, updatedat", input.Message, time.Now(), input.MessageID)
			err := row.Scan(&chatconversation.MessageStatus, &chatconversation.MessageParentID, &chatconversation.CreatedAt, &chatconversation.UpdatedAt)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatConversation{}, er.InternalServerError
			}
			chatconversation = model.ChatConversation{
				MessageID:  input.MessageID,
				ChatRoomID: input.ChatRoomID,
				SenderID:   input.SenderID,
				Message:    *input.Message,
			}
		}
	}
	channelMsg := updateMessageResolver[input.ChatRoomID]
	if channelMsg != nil {
		channelMsg <- chatconversation
	}
	return chatconversation, nil
}

func (r *subscriptionResolver) MessageUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	chatevent := make(chan model.ChatConversation, 1)
	updateMessageResolver[chatRoomID] = chatevent
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(updateMessageResolver, chatRoomID)
		r.mu.Unlock()
	}()
	return chatevent, nil
}

func (r *mutationResolver) UpdateMessageStatus(ctx context.Context, input model.UpdateMessageStatus) (model.ChatConversation, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) MessageStatusUpdate(ctx context.Context, messageID int, chatRoomID int) (<-chan model.ChatConversation, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, input *model.DeleteMessage) (model.ChatConversation, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation

	isMemberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, input.DeleteByID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isMessageOwner, err := CheckMemberOwner(ctx, input.MessageID, input.DeleteByID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
		if isMessageOwner {
			_, err := crConn.Db.Exec("DELETE FROM chatconversation WHERE id = $1", input.MessageID)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatConversation{}, er.InternalServerError
			}
		}
	}
	return chatconversation, nil
}

func (r *subscriptionResolver) MessageDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	panic("not implemented")
}

func (r *chatConversationResolver) Sender(ctx context.Context, obj *model.ChatConversation) (model.User, error) {
	crconn := ctx.Value("crConn").(*dal.DbConnection)
	var sender model.User
	rows, err := crconn.Db.Query("SELECT id, username, first_name, last_name, email, contact, bio, profile_picture, created_at, updated_at FROM users WHERE id = $1", obj.SenderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.User{}, er.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&sender.ID, &sender.Username, &sender.FirstName, &sender.LastName, &sender.Email, &sender.Contact, &sender.Bio, &sender.ProfilePicture, &sender.CreatedAt, &sender.UpdatedAt)
		if err != nil {
			er.DebugPrintf(err)
			return model.User{}, er.InternalServerError
		}
	}
	return sender, nil
}

func ChehckDeleteFlag(ctx context.Context, chatRoomID int, memberID int) (*time.Time, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	row := crConn.Db.QueryRow("SELECT deleted_at FROM members WHERE chatroom_id = $1 AND member_id = $2", chatRoomID, memberID)
	err := row.Scan(&member.DeleteAt)
	if err != nil {
		er.DebugPrintf(err)
		return nil, er.InternalServerError
	}
	return member.DeleteAt, nil
}
func CheckMemberOwner(ctx context.Context, messageID int, senderID int) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var isMessageOwner bool
	row := crConn.Db.QueryRow("SELECT true FROM chatconversation WHERE id = $1 And sender_id = $2", messageID, senderID)
	err := row.Scan(isMessageOwner)
	if err != nil && err == sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isMessageOwner, nil
}
