package handler

import (
	"context"
	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/model"
	"log"
)

func (r *mutationResolver) NewChatRoomMember(ctx context.Context, input model.NewChatRoomMember, receiverId *int) (model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	member := model.Member{
		ChatRoomID: input.ChatRoomID,
		MemberID: input.MemberID,
	}
	// Fetch chatRoomType
	chatRoomType, err := CheckChatRoomTypeByChatID(ctx, input.ChatRoomID)
	if err != nil{
		log.Println("Error at CheckChatRoomTypeByChatID function", err)
	}
	if chatRoomType == "GROUP" {
		// Fetch Member existence in  chatRoom
		checkMemberExist, err := CheckMemberExistence(ctx, input.ChatRoomID, input.MemberID)
		if err != nil{
			log.Println("Error at checkMemberExist function", err)
		}
		if checkMemberExist.MemberID == 0 {
			_, err := crConn.Db.Exec("INSERT INTO members_test (chatroom_id, member_id, joinat) values ($1, $2, NOW())", input.ChatRoomID, input.MemberID)
			if err != nil {
				log.Println("Error to insert new member in chatroom", err)
			}
		}
	} else if( chatRoomType == "PRIVATE" ){
		// Fetch Total members in chatroom
		totalMember, err := ChatRoomTotalMemberByChatRoomId(ctx, input.ChatRoomID)
		if err != nil{
			log.Println("Error at ChatRoomTotalMemberByChatRoomId function", err)
		}
			if totalMember < 2{
				_, err := crConn.Db.Exec("INSERT INTO members_test (chatroom_id, member_id, joinat) values ($1, $2, NOW()),($1,$3,NOW())", input.ChatRoomID, input.MemberID, receiverId)
				if err != nil {
					log.Println("Error to insert new member in chatroom", err)
				}
			}
	}
	return member, nil
}
func (r *memberResolver) Member(ctx context.Context, obj *model.Member) (model.User, error) {
	var memberInfo model.User
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	row, _ := crConn.Db.Query("SELECT user_test.id, name, email, contact, profile_picture, bio, createdat FROM user_test,members_test WHERE user_test.id = $1 and user_test.id = members_test.member_id", obj.MemberID)
	defer row.Close()
	for row.Next() {
		err := row.Scan(&memberInfo.ID, &memberInfo.Name, &memberInfo.Email, &memberInfo.Contact, &memberInfo.ProfilePicture, &memberInfo.Bio, &memberInfo.CreatedAt)
		if err != nil {
			log.Println("Error to scan user data as per memberid at line 24 of members_handler", err)
		}
	}
	return  memberInfo, nil
}

func CheckMemberExistence(ctx context.Context, chatRoomId int, memberId int) (model.Member, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var member model.Member
	rows, _ := crConn.Db.Query("SELECT id, chatroom_id, member_id, joinat FROM members_test WHERE chatroom_id = $1 and member_id = $2", chatRoomId, memberId)
	for rows.Next() {
		err := rows.Scan(&member.ID, &member.ChatRoomID, &member.MemberID, &member.JoinAt)
		if err != nil {
			log.Println("Error to read member data as per chatroom id", err)
		}
	}
	return member, nil
}
func CheckChatRoomTypeByChatID(ctx context.Context, chatRoomId int) (string, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var chatroom model.ChatRoom
	row := crConn.Db.QueryRow("SELECT chatroom_type FROM chatroom_test WHERE id = $1", chatRoomId)
	err := row.Scan(&chatroom.ChatRoomType)
	if err != nil{
		log.Println("Error to fetch chatroom type by chatroom id  of members_handler", err)
	}
	return  chatroom.ChatRoomType.String(), nil
}
func ChatRoomTotalMemberByChatRoomId(ctx context.Context, chatRoomId int)(int, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var totalChatRoomMember int
	row := crConn.Db.QueryRow("select count(member_id) from members_test join chatroom_test on members_test.chatroom_id = chatroom_test.id where chatroom_id = $1", chatRoomId)
	err := row.Scan(&totalChatRoomMember)
	if err != nil{
		log.Println("Error to count total member in chatroom at  of members_handler", err)
	}
	return totalChatRoomMember, nil
}