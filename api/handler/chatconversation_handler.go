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
	isMemberExist, err := checkMemberExistence(crConn, chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isDeleted, err := checkIsDeleted(crConn, chatRoomID, memberID)
		if err != nil {
			er.DebugPrintf(err)
			return []model.ChatConversation{}, er.InternalServerError
		}

		if isDeleted {
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

func (r *queryResolver) MemberListByChatRoomID(ctx context.Context, chatRoomID int, memberID int) (model.MemberCountsWithMemberDetailsByChatRoom, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	var members []model.Member
	var membersCountWithMemberDetail model.MemberCountsWithMemberDetailsByChatRoom
	isMemberExist, err := checkMemberExistence(crConn, chatRoomID, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return model.MemberCountsWithMemberDetailsByChatRoom{}, er.InternalServerError
	}
	if isMemberExist {
		rows, err := crConn.Db.Query("SELECT members.id, chatroom_id, member_id, joined_at FROM members, users WHERE members.member_id = users.id AND chatroom_id = $1 ORDER BY users.first_name ASC", chatRoomID)
		if err != nil {
			er.DebugPrintf(err)
			return model.MemberCountsWithMemberDetailsByChatRoom{}, er.InternalServerError
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&member.ID, &member.ChatRoomID, &member.MemberID, &member.JoinAt)
			if err != nil {
				er.DebugPrintf(err)
				return model.MemberCountsWithMemberDetailsByChatRoom{}, er.InternalServerError
			}
			members = append(members, member)
		}
	}
	membersCountWithMemberDetail.MemberCount = len(members)
	membersCountWithMemberDetail.Members = members
	return membersCountWithMemberDetail, nil
}

func (r *queryResolver) ChatRoomListByMemberID(ctx context.Context, memberID int) ([]model.ChatRoomList, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	chatroomList, err := chatRoomListByMemberID(crConn, memberID)
	if err != nil {
		er.DebugPrintf(err)
		return []model.ChatRoomList{}, er.InternalServerError
	}
	return chatroomList, nil
}

func (r *subscriptionResolver) ChatRoomListByMember(ctx context.Context, memberID int) (<-chan []model.ChatRoomList, error) {
	r.mu.Lock()
	chatroomlist := g.ChatRoomList[memberID]
	if chatroomlist == nil {
		chatroomlist = &model.Member{MemberID: memberID, ChatRoomListObservers: map[int]chan []model.ChatRoomList{}}
		g.ChatRoomList[memberID] = chatroomlist
	}
	r.mu.Unlock()
	id, _ := strconv.Atoi(helper.RandString(20))
	chatroomListEvenet := make(chan []model.ChatRoomList, 1)
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(chatroomlist.ChatRoomListObservers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()
	chatroomlist.ChatRoomListObservers[id] = chatroomListEvenet
	r.mu.Unlock()
	return chatroomListEvenet, nil
}

// Create New Message
func (r *mutationResolver) NewMessage(ctx context.Context, input model.NewMessage) (model.ChatConversation, error) {
	r.mu.Lock()
	addChatConvo := g.AddMessages[input.ChatRoomID]
	if addChatConvo == nil {
		addChatConvo = &model.ChatRoom{ChatRoomID: input.ChatRoomID, AddMessageObservers: map[int]chan model.ChatConversation{}}
		g.AddMessages[input.ChatRoomID] = addChatConvo
	}
	r.mu.Unlock()
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	chatconversation := model.ChatConversation{
		ChatRoomID:      input.ChatRoomID,
		SenderID:        input.SenderID,
		Message:         input.Message,
		MessageType:     input.MessageType,
		MessageParentID: input.MessageParentID,
		MessageStatus:   input.MessageStatus,
	}
	isMemberExist, err := checkMemberExistence(crConn, input.ChatRoomID, input.SenderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		row := crConn.Db.QueryRow("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at", input.ChatRoomID, input.SenderID, input.Message, input.MessageType, input.MessageStatus, input.MessageParentID, time.Now().UTC())
		err := row.Scan(&chatconversation.MessageID, &chatconversation.CreatedAt)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
		_, err = updateDeleteFlagAndSetZero(crConn, input.ChatRoomID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
	}

	addChatConvo.ChatConversations = append(addChatConvo.ChatConversations, chatconversation)
	r.mu.Lock()
	for _, observer := range addChatConvo.AddMessageObservers {
		observer <- chatconversation
	}
	r.mu.Unlock()
	rows, err := crConn.Db.Query("SELECT member_id FROM members WHERE chatroom_id = $1 AND member_id != $2", input.ChatRoomID, input.SenderID)
	defer rows.Close()
	for rows.Next() {
		var receiverID int
		err = rows.Scan(&receiverID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
		_, err := chatRoomListByMemberID(crConn, receiverID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
	}
	_, err = chatRoomListByMemberID(crConn, input.SenderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	return chatconversation, nil
}

// Live updates of new messages
func (r *subscriptionResolver) MessagePost(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	r.mu.Lock()
	addChatConvo := g.AddMessages[chatRoomID]
	if addChatConvo == nil {
		addChatConvo = &model.ChatRoom{ChatRoomID: chatRoomID, AddMessageObservers: map[int]chan model.ChatConversation{}}
		g.AddMessages[chatRoomID] = addChatConvo
	}
	r.mu.Unlock()
	id := helper.Random(1, 99999999999999)
	addMessageEvent := make(chan model.ChatConversation, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(addChatConvo.AddMessageObservers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()
	addChatConvo.AddMessageObservers[id] = addMessageEvent
	r.mu.Unlock()
	return addMessageEvent, nil
}

func (r *mutationResolver) UpdateMessage(ctx context.Context, input *model.UpdateMessage) (model.ChatConversation, error) {
	r.mu.Lock()
	updateChatConvo := g.UpdateMessage[input.ChatRoomID]
	if updateChatConvo == nil {
		updateChatConvo = &model.ChatRoom{ChatRoomID: input.ChatRoomID, UpdateMessageObservers: map[int]chan model.ChatConversation{}}
		g.UpdateMessage[input.ChatRoomID] = updateChatConvo
	}
	r.mu.Unlock()
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	chatconversation = model.ChatConversation{
		MessageID:  input.MessageID,
		ChatRoomID: input.ChatRoomID,
		SenderID:   input.SenderID,
		Message:    *input.Message,
	}
	isMemberExist, err := checkMemberExistence(crConn, input.ChatRoomID, input.SenderID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isMessageOwner, err := CheckMemberOwner(crConn, input.MessageID, input.SenderID)
		if err != nil {
			er.DebugPrintf(err)
			return model.ChatConversation{}, er.InternalServerError
		}
		if isMessageOwner {
			row := crConn.Db.QueryRow("UPDATE chatconversation SET (message, updatedat) = ($1, $2) where id = $3 RETURNING message_status, message_parent_id, created_at, updatedat", input.Message, time.Now().UTC(), input.MessageID)
			err := row.Scan(&chatconversation.MessageStatus, &chatconversation.MessageParentID, &chatconversation.CreatedAt, &chatconversation.UpdatedAt)
			if err != nil {
				er.DebugPrintf(err)
				return model.ChatConversation{}, er.InternalServerError
			}
		}
	}
	updateChatConvo.ChatConversations = append(updateChatConvo.ChatConversations, chatconversation)
	r.mu.Lock()
	for _, observer := range updateChatConvo.UpdateMessageObservers {
		observer <- chatconversation
	}
	r.mu.Unlock()
	return chatconversation, nil
}

func (r *subscriptionResolver) MessageUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	r.mu.Lock()
	updateChatConvo := g.UpdateMessage[chatRoomID]
	if updateChatConvo == nil {
		updateChatConvo = &model.ChatRoom{ChatRoomID: chatRoomID, UpdateMessageObservers: map[int]chan model.ChatConversation{}}
		g.UpdateMessage[chatRoomID] = updateChatConvo
	}
	r.mu.Unlock()
	id := helper.Random(1, 100000000000000)
	updateMessageEvent := make(chan model.ChatConversation, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(updateChatConvo.UpdateMessageObservers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()
	updateChatConvo.UpdateMessageObservers[id] = updateMessageEvent
	r.mu.Unlock()
	return updateMessageEvent, nil
}

func (r *mutationResolver) UpdateMessageStatus(ctx context.Context, input model.UpdateMessageStatus) (model.ChatConversation, error) {
	panic("not implemented")
}

func (r *subscriptionResolver) MessageStatusUpdate(ctx context.Context, messageID int, chatRoomID int) (<-chan model.ChatConversation, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, input *model.DeleteMessage) (model.ChatConversation, error) {
	r.mu.Lock()
	deleteChatConvo := g.DeleteMessage[input.ChatRoomID]
	if deleteChatConvo == nil {
		deleteChatConvo = &model.ChatRoom{ChatRoomID: input.ChatRoomID, DeleteMessageObservers: map[int]chan model.ChatConversation{}}
		g.DeleteMessage[input.ChatRoomID] = deleteChatConvo
	}
	r.mu.Unlock()
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	isMemberExist, err := checkMemberExistence(crConn, input.ChatRoomID, input.DeleteByID)
	if err != nil {
		er.DebugPrintf(err)
		return model.ChatConversation{}, er.InternalServerError
	}
	if isMemberExist {
		isMessageOwner, err := CheckMemberOwner(crConn, input.MessageID, input.DeleteByID)
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
	r.mu.Lock()
	for _, observer := range deleteChatConvo.DeleteMessageObservers {
		observer <- chatconversation
	}
	r.mu.Unlock()
	return chatconversation, nil
}

func (r *subscriptionResolver) MessageDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	r.mu.Lock()
	deleteChatConvo := g.DeleteMessage[chatRoomID]
	if deleteChatConvo == nil {
		deleteChatConvo = &model.ChatRoom{ChatRoomID: chatRoomID, DeleteMessageObservers: map[int]chan model.ChatConversation{}}
		g.DeleteMessage[chatRoomID] = deleteChatConvo
	}
	r.mu.Unlock()
	id := helper.Random(1, 11111111111111)
	deleteMessageEvent := make(chan model.ChatConversation, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(deleteChatConvo.DeleteMessageObservers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()
	deleteChatConvo.DeleteMessageObservers[id] = deleteMessageEvent
	r.mu.Unlock()
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

func checkIsDeleted(crConn *dal.DbConnection, chatRoomID int, memberID int) (bool, error) {
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

func CheckMemberOwner(crConn *dal.DbConnection, messageID int, senderID int) (bool, error) {
	var isMessageOwner bool

	row := crConn.Db.QueryRow("SELECT true FROM chatconversation WHERE id = $1 And sender_id = $2", messageID, senderID)
	err := row.Scan(&isMessageOwner)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isMessageOwner, nil
}

func fetchChatRoomList(crConn *dal.DbConnection, memberID int) ([]model.ChatRoomList, error) {
	var chatroomLists []model.ChatRoomList

	var row *sql.Row
	//TODO: Need to change query (Solve issue in distinct with order by)
	rows, err := crConn.Db.Query("select distinct (chatrooms.id), chatrooms.chatroom_type, members.delete_flag, max(chatconversation.created_at) from chatconversation join members on members.chatroom_id = chatconversation.chatroom_id join chatrooms on chatrooms.id = members.chatroom_id where members.member_id = $1 group by chatrooms.id, chatrooms.chatroom_type, members.delete_flag order by max(chatconversation.created_at) desc", memberID)
	if err != nil {
		//er.DebugPrintf(err)
		return []model.ChatRoomList{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var deleteFlag int
		var chatroomList model.ChatRoomList
		var create_at time.Time
		err := rows.Scan(&chatroomList.ChatRoomID, &chatroomList.ChatRoomType, &deleteFlag, &create_at)
		if err != nil {
			//er.DebugPrintf(err)
			return []model.ChatRoomList{}, err
		}

		switch chatroomList.ChatRoomType {
		case "PRIVATE":
			if deleteFlag == 0 {
				row = crConn.Db.QueryRow("SELECT chatrooms.id,username AS name, chatrooms.chatroom_type, chatrooms.created_at FROM users JOIN members ON members.member_id = users.id JOIN chatrooms ON chatrooms.id = members.chatroom_id WHERE chatrooms.id = $1 AND members.member_id != $2", chatroomList.ChatRoomID, memberID)
				err = row.Scan(&chatroomList.ChatRoomID, &chatroomList.Name, &chatroomList.ChatRoomType, &chatroomList.CreatedAt)
				if err != nil {
					//er.DebugPrintf(err)
					return []model.ChatRoomList{}, err
				}
				chatroomLists = append(chatroomLists, chatroomList)
			}
		default:
			row = crConn.Db.QueryRow("SELECT chatrooms.id, chatrooms.chatroom_name AS name, chatrooms.chatroom_type, chatrooms.created_at, count(members.member_id) FROM chatrooms join members on members.chatroom_id = chatrooms.id WHERE  chatrooms.id = $1 group by chatrooms.id, chatrooms.chatroom_name, chatrooms.chatroom_type, chatrooms.created_at", chatroomList.ChatRoomID)
			err = row.Scan(&chatroomList.ChatRoomID, &chatroomList.Name, &chatroomList.ChatRoomType, &chatroomList.CreatedAt, &chatroomList.TotalMember)
			if err != nil {
				//er.DebugPrintf(err)
				return []model.ChatRoomList{}, err
			}
			chatroomLists = append(chatroomLists, chatroomList)
		}
	}
	return chatroomLists, nil
}

func chatRoomListByMemberID(crConn *dal.DbConnection, memberID int) ([]model.ChatRoomList, error) {
	var chatroomLists []model.ChatRoomList
	var err error
	//r.mu.Lock()
	chatroomlist := g.ChatRoomList[memberID]
	if chatroomlist == nil {
		chatroomlist = &model.Member{MemberID: memberID, ChatRoomListObservers: map[int]chan []model.ChatRoomList{}}
		g.ChatRoomList[memberID] = chatroomlist
	}
	//r.mu.Unlock()
	chatroomLists, err = fetchChatRoomList(crConn, memberID)
	if err != nil {
		//er.DebugPrintf(err)
		return []model.ChatRoomList{}, err
	}
	chatroomlist.ChatRoomLists = append(chatroomlist.ChatRoomLists, chatroomLists)
	//r.mu.Lock()
	for _, observer := range chatroomlist.ChatRoomListObservers {
		observer <- chatroomLists
	}
	//r.mu.Unlock()
	return chatroomLists, nil
}
