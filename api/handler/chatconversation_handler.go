package handler

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"
	er "github.com/aneri/new_chat/error"
	"github.com/aneri/new_chat/model"
)

// Retrieve chat conversation by chatRoom Id
func (r *queryResolver) ChatconversationByChatRoomID(ctx context.Context, chatRoomID int, memberID int) ([]model.ChatConversation, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	var chattconversationarr []model.ChatConversation
	var rows *sql.Rows
	isMemberExist, err := checkMemberExistence(ctx, chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isDeleted, err := checkIsDeleted(ctx, chatRoomID, memberID)
		if err != nil {
			er.DebugPrintf(err)
			return []model.ChatConversation{}, er.InternalServerError
		}

		if isDeleted  {
				rows, err = crConn.Db.Query("SELECT chatconversation.id, chatconversation.chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at, updatedat FROM chatconversation LEFT JOIN members ON members.deleted_at <= chatconversation.created_at WHERE chatconversation.chatroom_id = $1 AND chatconversation.chatroom_id = members.chatroom_id GROUP BY chatconversation.chatroom_id,sender_id, message,message_type, message_parent_id, message_status, created_at, updatedat, chatconversation.id ORDER BY chatconversation.id ASC", chatRoomID)
			if err != nil {
				er.DebugPrintf(err)
				return []model.ChatConversation{}, er.InternalServerError
			}
		} else {
			rows, err = crConn.Db.Query("SELECT chatconversation.id, chatconversation.chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at, updatedat FROM chatconversation LEFT JOIN members ON members.joined_at <= chatconversation.created_at WHERE chatconversation.chatroom_id = $1 AND chatconversation.chatroom_id = members.chatroom_id GROUP BY chatconversation.chatroom_id,sender_id, message,message_type, message_parent_id, message_status, created_at, updatedat, chatconversation.id ORDER BY chatconversation.id ASC", chatRoomID)
			if err != nil {
				er.DebugPrintf(err)
				return []model.ChatConversation{}, er.InternalServerError
			}
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
	}
	return chattconversationarr, nil
}

func (r *queryResolver) MemberListByChatRoomID(ctx context.Context, chatRoomID int, memberID int) ([]model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	var members []model.Member
	isMemberExist, err := checkMemberExistence(ctx, chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.Member{}, er.InternalServerError
	}
	if isMemberExist {
		rows, err := crConn.Db.Query("SELECT members.id, chatroom_id, member_id, joined_at FROM members, users WHERE members.member_id = users.id AND chatroom_id = $1 ORDER BY users.first_name ASC", chatRoomID)
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

func (r *queryResolver) ChatRoomListByMemberID(ctx context.Context, memberID int) ([]model.ChatRoomList, error) {
	chatroomlist := g.ChatRoomList[memberID]
	if chatroomlist == nil {
		chatroomlist = &model.Member{MemberID: memberID, ChatRoomListObservers: map[int]chan []model.ChatRoomList{}}
		g.ChatRoomList[memberID] = chatroomlist
	}

	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroomLists []model.ChatRoomList
	var chatroomList model.ChatRoomList
	var row *sql.Row
	rows, err := crConn.Db.Query("select distinct(chatrooms.id), chatrooms.chatroom_type from chatconversation join members on members.chatroom_id = chatconversation.chatroom_id join chatrooms on chatrooms.id = members.chatroom_id where members.member_id = $1 order by chatconversation.created_at desc", memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatRoomList{}, er.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&chatroomList.ChatRoomID, &chatroomList.ChatRoomType)
		if err != nil {
			er.DebugPrintf(err)
			return []model.ChatRoomList{}, er.InternalServerError
		}

		switch chatroomList.ChatRoomType {

		case "PRIVATE":
			row = crConn.Db.QueryRow("SELECT chatrooms.id,username AS name, chatrooms.chatroom_type, chatrooms.created_at FROM users JOIN members ON members.member_id = users.id JOIN chatrooms ON chatrooms.id = members.chatroom_id WHERE chatrooms.id = $1 AND members.member_id != $2", chatroomList.ChatRoomID, memberID)
		default:
			row = crConn.Db.QueryRow("SELECT chatrooms.id, chatrooms.chatroom_name AS name, chatrooms.chatroom_type, chatrooms.created_at FROM chatrooms WHERE  chatrooms.id = $1", chatroomList.ChatRoomID)
		}
		err = row.Scan(&chatroomList.ChatRoomID, &chatroomList.Name, &chatroomList.ChatRoomType, &chatroomList.CreatedAt)
		if err != nil {
			er.DebugPrintf(err)
			return []model.ChatRoomList{}, er.InternalServerError
		}
		chatroomLists = append(chatroomLists, chatroomList)
	}
	chatroomlist.ChatRoomLists = append(chatroomlist.ChatRoomLists, chatroomLists)
	for _, observer := range chatroomlist.ChatRoomListObservers {
		observer <- chatroomLists
	}
	return chatroomLists, nil
}
func (r *subscriptionResolver) ChatRoomListByMember(ctx context.Context, memberID int) (<-chan []model.ChatRoomList, error) {
	chatroomlist := g.ChatRoomList[memberID]
	if chatroomlist == nil {
		chatroomlist = &model.Member{MemberID: memberID, ChatRoomListObservers: map[int]chan []model.ChatRoomList{}}
		g.ChatRoomList[memberID] = chatroomlist
	}
	id, _ := strconv.Atoi(helper.RandString(20))
	chatroomListEvenet := make(chan []model.ChatRoomList, 1)
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(chatroomlist.ChatRoomListObservers, id)
		r.mu.Unlock()
	}()
	chatroomlist.ChatRoomListObservers[id] = chatroomListEvenet
	return chatroomListEvenet, nil
}

// Create New Message
func (r *mutationResolver) NewMessage(ctx context.Context, input model.NewMessage) (model.ChatConversation, error) {
	addChatConvo := g.AddMessages[input.ChatRoomID]
	if addChatConvo == nil {
		addChatConvo = &model.ChatRoom{ChatRoomID: input.ChatRoomID, AddMessageObservers: map[int]chan model.ChatConversation{}}
		g.AddMessages[input.ChatRoomID] = addChatConvo
	}

	crConn := ctx.Value("crConn").(*dal.DbConnection)
	chatconversation := model.ChatConversation{
		ChatRoomID:      input.ChatRoomID,
		SenderID:        input.SenderID,
		Message:         input.Message,
		MessageType:     input.MessageType,
		MessageParentID: input.MessageParentID,
		MessageStatus:   input.MessageStatus,
	}
	isMemberExist, err := checkMemberExistence(ctx, input.ChatRoomID, input.SenderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		row := crConn.Db.QueryRow("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at", input.ChatRoomID, input.SenderID, input.Message, input.MessageType, input.MessageStatus, input.MessageParentID, time.Now())
		err := row.Scan(&chatconversation.MessageID, &chatconversation.CreatedAt)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
	}
	addChatConvo.ChatConversations = append(addChatConvo.ChatConversations, chatconversation)
	for _, observer := range addChatConvo.AddMessageObservers {
		observer <- chatconversation
	}
	return chatconversation, nil
}

// Live updates of new messages
func (r *subscriptionResolver) MessagePost(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	addChatConvo := g.AddMessages[chatRoomID]
	if addChatConvo == nil {
		addChatConvo = &model.ChatRoom{ChatRoomID: chatRoomID, AddMessageObservers: map[int]chan model.ChatConversation{}}
		g.AddMessages[chatRoomID] = addChatConvo
	}

	id := helper.Random(1, 99999999999999)
	addMessageEvent := make(chan model.ChatConversation, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(addChatConvo.AddMessageObservers, id)
		r.mu.Unlock()
	}()

	addChatConvo.AddMessageObservers[id] = addMessageEvent

	return addMessageEvent, nil
}

func (r *mutationResolver) UpdateMessage(ctx context.Context, input *model.UpdateMessage) (model.ChatConversation, error) {
	updateChatConvo := g.UpdateMessage[input.ChatRoomID]
	if updateChatConvo == nil {
		updateChatConvo = &model.ChatRoom{ChatRoomID: input.ChatRoomID, UpdateMessageObservers: map[int]chan model.ChatConversation{}}
		g.UpdateMessage[input.ChatRoomID] = updateChatConvo
	}

	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	chatconversation = model.ChatConversation{
		MessageID:  input.MessageID,
		ChatRoomID: input.ChatRoomID,
		SenderID:   input.SenderID,
		Message:    *input.Message,
	}
	isMemberExist, err := checkMemberExistence(ctx, input.ChatRoomID, input.SenderID)
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
		}
	}
	updateChatConvo.ChatConversations = append(updateChatConvo.ChatConversations, chatconversation)
	for _, observer := range updateChatConvo.UpdateMessageObservers {
		observer <- chatconversation
	}
	return chatconversation, nil
}

func (r *subscriptionResolver) MessageUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	updateChatConvo := g.UpdateMessage[chatRoomID]
	if updateChatConvo == nil {
		updateChatConvo = &model.ChatRoom{ChatRoomID: chatRoomID, UpdateMessageObservers: map[int]chan model.ChatConversation{}}
		g.UpdateMessage[chatRoomID] = updateChatConvo
	}

	id := helper.Random(1, 100000000000000)
	updateMessageEvent := make(chan model.ChatConversation, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(updateChatConvo.UpdateMessageObservers, id)
		r.mu.Unlock()
	}()

	updateChatConvo.UpdateMessageObservers[id] = updateMessageEvent

	return updateMessageEvent, nil
}

func (r *mutationResolver) UpdateMessageStatus(ctx context.Context, input model.UpdateMessageStatus) (model.ChatConversation, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) MessageStatusUpdate(ctx context.Context, messageID int, chatRoomID int) (<-chan model.ChatConversation, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, input *model.DeleteMessage) (model.ChatConversation, error) {
	deleteChatConvo := g.DeleteMessage[input.ChatRoomID]
	if deleteChatConvo == nil {
		deleteChatConvo = &model.ChatRoom{ChatRoomID: input.ChatRoomID, DeleteMessageObservers: map[int]chan model.ChatConversation{}}
		g.DeleteMessage[input.ChatRoomID] = deleteChatConvo
	}
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	isMemberExist, err := checkMemberExistence(ctx, input.ChatRoomID, input.DeleteByID)
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
			row := crConn.Db.QueryRow("DELETE FROM chatconversation WHERE id = $1 RETURNING id, chatroom_id, message", input.MessageID)
			err := row.Scan(&chatconversation.MessageID, &chatconversation.ChatRoomID, &chatconversation.Message)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatConversation{}, er.InternalServerError
			}
		}
	}
	for _, observer := range deleteChatConvo.DeleteMessageObservers {
		observer <- chatconversation
	}
	return chatconversation, nil
}

func (r *subscriptionResolver) MessageDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	deleteChatConvo := g.DeleteMessage[chatRoomID]
	if deleteChatConvo == nil {
		deleteChatConvo = &model.ChatRoom{ChatRoomID: chatRoomID, DeleteMessageObservers: map[int]chan model.ChatConversation{}}
		g.DeleteMessage[chatRoomID] = deleteChatConvo
	}

	id := helper.Random(1, 11111111111111)
	deleteMessageEvent := make(chan model.ChatConversation, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(deleteChatConvo.DeleteMessageObservers, id)
		r.mu.Unlock()
	}()

	deleteChatConvo.DeleteMessageObservers[id] = deleteMessageEvent

	return deleteMessageEvent, nil
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

func checkIsDeleted(ctx context.Context, chatRoomID int, memberID int) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	row := crConn.Db.QueryRow("SELECT deleted_at FROM members WHERE chatroom_id = $1 AND member_id = $2", chatRoomID, memberID)
	err := row.Scan(&member.DeleteAt)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	if member.DeleteAt == nil {
		return false, nil
	}
	return true, nil
}
func CheckMemberOwner(ctx context.Context, messageID int, senderID int) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var isMessageOwner bool
	row := crConn.Db.QueryRow("SELECT true FROM chatconversation WHERE id = $1 And sender_id = $2", messageID, senderID)
	err := row.Scan(&isMessageOwner)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isMessageOwner, nil
}
