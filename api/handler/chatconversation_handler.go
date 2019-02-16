package handler

import (
	"context"
	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/model"
	"log"
)

// Retrieve chat conversation by chatRoom Id
func (r *queryResolver) ChatconversationByChatRoomID(ctx context.Context, chatRoomID int, memberID int)([]model.ChatConversation, error) {
	crConn :=ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	var chattconversationarr []model.ChatConversation
	checkChatRoomMember, err := CheckMemberExistence(ctx, chatRoomID, memberID)
	if err != nil{
		log.Println("Error", err)
	}
	if checkChatRoomMember.DeleteAt != nil{
		rows, err := crConn.Db.Query("SELECT chatconversation.id, chatconversation.chatroom_id, sender_id, message, message_type, message_parent_id, message_status, createat, updateat FROM chatconversation LEFT JOIN members_test ON members_test.deleteat <= chatconversation.createat WHERE chatconversation.chatroom_id = $1 AND deleteflag = 0 AND chatconversation.chatroom_id = members_test.chatroom_id GROUP BY chatconversation.chatroom_id,sender_id, message,message_type, message_parent_id, message_status, createat, updateat, chatconversation.id ORDER BY chatconversation.id", chatRoomID)
		if err != nil{
			log.Println("Error to fetch chatconversation by chatroom id when chatroom is delete by particular person", err)
		}
		defer rows.Close()
		for rows.Next(){
			err := rows.Scan(&chatconversation.MessageID, &chatconversation.ChatRoomID, &chatconversation.SenderID, &chatconversation.Message, &chatconversation.MessageType, &chatconversation.MessageParentID, &chatconversation.MessageStatus, &chatconversation.CreatedAt, &chatconversation.UpdatedAt)
			if err != nil{
				log.Println("Error to scan chat", err)
			}
			chattconversationarr = append(chattconversationarr, chatconversation)
		}
	}else if(checkChatRoomMember.DeleteAt == nil){
		rows, err := crConn.Db.Query("SELECT chatconversation.id, chatconversation.chatroom_id, sender_id, message, message_type, message_parent_id, message_status, createat, updateat FROM chatconversation LEFT JOIN members_test ON members_test.joinat <= chatconversation.createat WHERE chatconversation.chatroom_id = $1 AND deleteflag = 0 AND chatconversation.chatroom_id = members_test.chatroom_id GROUP BY chatconversation.chatroom_id,sender_id, message,message_type, message_parent_id, message_status, createat, updateat, chatconversation.id ORDER BY chatconversation.id", chatRoomID)
		if err != nil{
			log.Println("Error to fetch chatconversation", err)
		}
		defer rows.Close()
		for  rows.Next(){
			err := rows.Scan(&chatconversation.MessageID, &chatconversation.ChatRoomID, &chatconversation.SenderID, &chatconversation.Message, &chatconversation.MessageType, &chatconversation.MessageParentID, &chatconversation.MessageStatus, &chatconversation.CreatedAt, &chatconversation.UpdatedAt)
			if err != nil{
				log.Println("Error to scan chat", err)
			}
			chattconversationarr = append(chattconversationarr, chatconversation)
		}
	}
	return  chattconversationarr, nil
}

func (r *queryResolver)MemberListByChatRoomID(ctx context.Context, chatRoomID int, memberID int) ([]model.Member, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	var members []model.Member
	chehckMemberExist, err :=  CheckMemberExistence(ctx, chatRoomID, memberID)
	if  err != nil{
		log.Println("Error to fetch member id when fetching member list by chatroomid", err)
	}
	if chehckMemberExist.MemberID !=0 {
		rows, err := crConn.Db.Query("SELECT members_test.id, chatroom_id, member_id, joinat FROM members_test,user_test WHERE members_test.member_id = user_test.id AND chatroom_id = $1 ORDER BY user_test.name", chatRoomID)
		if err != nil{
			log.Println("Error to fetch member data by chatroonid", err)
		}
		defer rows.Close()
		for rows.Next(){
			err := rows.Scan(&member.ID, &member.ChatRoomID, &member.MemberID, &member.JoinAt)
			if err != nil{
				log.Println("Error at scanning member data", err)
			}
			members = append(members, member)
		}
	}
	return  members, nil
}

func (r *queryResolver)ChatRoomListByMemberID(ctx context.Context, memberId int) ([]model.ChatRoom, error){
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
	checkSenderExistence, err := CheckSenderExistence(ctx, input.ChatRoomID, senderID)
	if err != nil{
		log.Println("Error at checking that member is exist in chatroom", err)
	}
	if checkSenderExistence != 0{
		_, err := crConn.Db.Exec("INSERT INTO chatconversation (chatroom_id, sender_id, message, message_type, message_parent_id, message_status, createat) VALUES ($1, $2, $3, $4, $5, $6, NOW())", input.ChatRoomID, input.SenderID, input.Message, input.MessageType, input.MessageParentID, input.MessageStatus)
		if err != nil{
			log.Println("Error to insert new message in chatroom", err)
		}
		chatconversation = model.ChatConversation{
			ChatRoomID: input.ChatRoomID,
			SenderID: input.SenderID,
			Message: input.Message,
			MessageType: input.MessageType,
			MessageParentID: input.MessageParentID,
			MessageStatus: input.MessageStatus,
		}
	}
	return chatconversation, nil
}

// Live updates of new messages
func (r *subscriptionResolver) MessagePost(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
	panic("not implemented")
}

func (r *mutationResolver)UpdateMessage(ctx context.Context, input *model.UpdateMessage) (model.ChatConversation, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	checkMessageDetail, err := checkMessageDetail(ctx, input.MessageID)
	if err != nil{
		log.Println("Error to get message detail by messageId", err)
	}
	if checkMessageDetail.SenderID == input.SenderID{
		_, err := crConn.Db.Exec("UPDATE chatconversation SET (message, updateat) = ($1, NOW()) where id = $2", input.Message, input.MessageID)
		if err != nil{
			log.Println("Error to update message", err)
		}
		chatconversation = model.ChatConversation{
			MessageID: input.MessageID,
			ChatRoomID: checkMessageDetail.ChatRoomID,
			SenderID: input.SenderID,
			Message: *input.Message,
			MessageType: checkMessageDetail.MessageType,
			MessageParentID: checkMessageDetail.MessageParentID,
			MessageStatus: checkMessageDetail.MessageStatus,
			CreatedAt: checkMessageDetail.CreatedAt,
			UpdatedAt: chatconversation.UpdatedAt,
		}
	}
	return chatconversation, nil
}

func (r *subscriptionResolver)MessageUpdate(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error){
	panic("not implemented")
}

func (r *mutationResolver)UpdateMessageStatus(ctx context.Context, input model.UpdateMessageStatus) (model.ChatConversation, error){
	panic("not implemented")
}

func (r *subscriptionResolver)MessageStatusUpdate(ctx context.Context, messageID int, chatRoomID int) (<-chan model.ChatConversation, error){
	panic("not implemented")
}

func (r *mutationResolver)DeleteMessage(ctx context.Context, senderID int, messageID int) (model.ChatConversation, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	checkMessageDetail, err := checkMessageDetail(ctx, messageID)
	if err != nil{
		log.Println("Error to get message detail by messageId", err)
	}
	if checkMessageDetail.SenderID == senderID{
		_, err := crConn.Db.Exec("DELETE FROM chatconversation WHERE id = $1", messageID)
		if err != nil{
			log.Println("Error to delete message from chatconversation", err)
		}
		chatconversation = model.ChatConversation{
			MessageID: messageID,
			ChatRoomID: checkMessageDetail.ChatRoomID,
			SenderID: senderID,
			Message: checkMessageDetail.Message,
			MessageType: checkMessageDetail.MessageType,
			MessageParentID: checkMessageDetail.MessageParentID,
			MessageStatus: checkMessageDetail.MessageStatus,
			CreatedAt: checkMessageDetail.CreatedAt,
			UpdatedAt: chatconversation.UpdatedAt,
		}
	}
	return  chatconversation, nil
}

func(r *subscriptionResolver)MessageDelete(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error){
	panic("not implemented")
}

func (r *chatConversationResolver) Sender(ctx context.Context, obj *model.ChatConversation) (model.User, error) {
	crconn := ctx.Value("crConn").(*dal.DbConnection)
	var sender model.User
	rows, err := crconn.Db.Query("SELECT id, name, email, contact, profile_picture, bio, createdat FROM user_test WHERE id = $1", obj.SenderID)
	if err != nil{
		log.Println("Error to fetch user", err)
	}
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&sender.ID, &sender.Name, &sender.Email, &sender.Contact, &sender.ProfilePicture, &sender.Bio, &sender.CreatedAt)
		if err != nil{
			log.Println("Error to scan sender data in chatconversation function", err)
		}
	}
	return sender, nil
}

func CheckSenderExistence(ctx context.Context, chatRoomId int, senderId int)(int, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var sender model.Member
	row := crConn.Db.QueryRow("SELECT id FROM members_test WHERE chatroom_id  = $1 AND member_id = $2", chatRoomId, senderId)
	err :=  row.Scan(&sender.ID)
	if err != nil{
		log.Println("Error to scan id by using chatroom_id and member_id", err)
	}
	return sender.ID, nil
}
func checkMessageDetail(ctx context.Context, messageID int)(model.ChatConversation,error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var message model.ChatConversation
	row := crConn.Db.QueryRow("SELECT chatroom_id, sender_id, message_type, message_parent_id, message_status, createat FROM chatconversation WHERE id = $1", messageID)
	err := row.Scan(&message.ChatRoomID, &message.SenderID, &message.MessageType, &message.MessageParentID, &message.MessageStatus, &message.CreatedAt)
	if err != nil{
		log.Println("Error to fetch sender_id by message_id", err)
	}
	return message, nil
}