package handler

import (
	"context"
	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/model"
	"log"
	"time"
)

// Retrieve chat conversation by chatRoom Id
func (r *queryResolver) ChatconversationByChatRoomID(ctx context.Context, chatRoomID int, memberID int)([]model.ChatConversation, error) {
	crConn :=ctx.Value("crConn").(*dal.DbConnection)
	var chatconversation model.ChatConversation
	var chattconversationarr []model.ChatConversation
	checkChatRoomMember, err := CheckChatRoomMember(ctx, chatRoomID, memberID)
	if err != nil{
		log.Println("Error", err)
	}
	if checkChatRoomMember != nil{
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
	}else if(checkChatRoomMember == nil){
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
func (r *subscriptionResolver) PostMessage(ctx context.Context, chatRoomID int) (<-chan model.ChatConversation, error) {
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
func CheckChatRoomMember(ctx context.Context, chatRoomId int, memberId int)( *time.Time, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	row := crConn.Db.QueryRow("SELECT deleteat from members_test WHERE chatroom_id = $1 and member_id = $2", chatRoomId, memberId)
	err := row.Scan(&member.DeleteAt)
	if err != nil{
		log.Println("Error to fetch chatroom delete time", err)
	}
	return member.DeleteAt, nil
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