package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type User struct {
	ID             int        `json:"id"`
	Username       string     `json:"username"`
	FirstName      *string    `json:"firstName"`
	LastName       *string    `json:"lastName"`
	Email          string     `json:"email"`
	Contact        *string    `json:"contact"`
	Bio            *string    `json:"bio"`
	ProfilePicture *string    `json:"profilePicture"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

type NewUser struct {
	UserName       string  `json:"userName"`
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	Email          string  `json:"email"`
	Contact        *string `json:"contact"`
	ProfilePicture *string `json:"profilePicture"`
	Bio            *string `json:"bio"`
}

type UpdateUser struct {
	UserName  *string `json:"userName"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type ChatRoom struct {
	ChatRoomID             int                           `json:"chatRoomID"`
	CreatorID              int                           `json:"creatorID"`
	ChatRoomName           *string                       `json:"chatRoomName"`
	ChatRoomType           ChatRoomType                  `json:"chatRoomType"`
	CreatedAt              time.Time                     `json:"createdAt"`
	UpdateByID             *int                          `json:"updateByID"`
	UpdatedAt              *time.Time                    `json:"updatedAt"`
	DeleteAt               *time.Time                    `json:"deleteAt"`
	ChatConversations      []ChatConversation            `json:"chatConversation"`
	AddMessageObservers    map[int]chan ChatConversation `json:"addMessageObserver"`
	UpdateMessageObservers map[int]chan ChatConversation `json:"updateMessageObserver"`
	DeleteMessageObservers map[int]chan ChatConversation `json:"deleteMessageObserver"`
}

type NewGroupChatRoom struct {
	CreatorID    int          `json:"creatorID"`
	ChatRoomName string       `json:"chatRoomName"`
	ChatRoomType ChatRoomType `json:"chatRoomType"`
	ReceiverID   []int        `json:"receiverID"`
}

type NewPrivateChatRoom struct {
	CreatorID    int          `json:"creatorID"`
	ChatRoomType ChatRoomType `json:"chatRoomType"`
	ReceiverID   int          `json:"receiverID"`
}

type UpdateChatRoomDetail struct {
	ChatRoomID   int     `json:"chatRoomID"`
	ChatRoomName *string `json:"chatRoomName"`
	UpdateByID   *int    `json:"updateByID"`
}

type Member struct {
	ID                    int                         `json:"id"`
	ChatRoomID            int                         `json:"chatRoomID"`
	MemberID              int                         `json:"memberID"`
	JoinAt                time.Time                   `json:"joinAt"`
	DeleteAt              *time.Time                  `json:"deleteAt"`
	ChatRoomLists         [][]ChatRoomList            `json:"chatRoomLists"`
	ChatRoomListObservers map[int]chan []ChatRoomList `json:"chatRoomList"`
}
type NewChatRoomMembers struct {
	ChatRoomID int   `json:"chatRoomID"`
	CreatorID int    `json:creatorID`
	MemberIDs  []int `json:"memberIDs"`
}

type LeaveChatRoom struct {
	ChatRoomID int `json:"chatRoomID"`
	MemberID   int `json:"memberID"`
}

type DeleteChat struct {
	ChatRoomID int `json:"chatRoomID"`
	MemberID   int `json:"memberID"`
}

type ChatConversation struct {
	MessageID       int         `json:"messageId"`
	ChatRoomID      int         `json:"chatRoomID"`
	SenderID        int         `json:"senderID"`
	Message         string      `json:"message"`
	MessageType     MessageType `json:"messageType"`
	MessageStatus   State       `json:"messageStatus"`
	MessageParentID *int        `json:"messageParentId"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       *time.Time  `json:"updatedAt"`
}

type NewMessage struct {
	ChatRoomID      int         `json:"chatRoomID"`
	SenderID        int         `json:"senderID"`
	Message         string      `json:"message"`
	MessageType     MessageType `json:"messageType"`
	MessageParentID *int        `json:"messageParentId"`
	MessageStatus   State       `json:"messageStatus"`
}

type UpdateMessage struct {
	Message    *string `json:"message"`
	SenderID   int     `json:"senderID"`
	MessageID  int     `json:"messageID"`
	ChatRoomID int     `json:"chatRoomID`
}
type DeleteMessage struct {
	ChatRoomID int `json:"chatRoomID"`
	MessageID  int `json:"messageID"`
	DeleteByID int `json:"DeleteByID"`
}
type ChatRoomList struct {
	ChatRoomID   int       `json:"chatRoomID"`
	Name         *string   `json:"name"`
	ChatRoomType string    `json:"chatRoomType"`
	CreatedAt    time.Time `json:"createdAt"`
	TotalMember  *int		`json:"totalMember"`
}

type UpdateMessageStatus struct {
	MessageStatus State `json:"messageStatus"`
}
type MemberCountsWithMemberDetailsByChatRoom struct {
	MemberCount int      `json:"memberCount"`
	Members     []Member `json:"members"`
}

type RemoveMembersFromChatRoom struct {
	ChatRoomID int `json:"chatRoomID"`
	CreatorID  int `json:"creatorID"`
	MemberID   int `json:"memberID"`
}
type RemoveChatRoomPermenantly struct {
	ChatRoomID int `json:"chatRoomID"`
	Creator    int `json:"creator"`
}

type ChatRoomType string

const (
	ChatRoomTypePrivate ChatRoomType = "PRIVATE"
	ChatRoomTypeGroup   ChatRoomType = "GROUP"
)

func (e ChatRoomType) IsValid() bool {
	switch e {
	case ChatRoomTypePrivate, ChatRoomTypeGroup:
		return true
	}
	return false
}

func (e ChatRoomType) String() string {
	return string(e)
}

func (e *ChatRoomType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChatRoomType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChatRoomType", str)
	}
	return nil
}

func (e ChatRoomType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MessageType string

const (
	MessageTypeText  MessageType = "TEXT"
	MessageTypeImage MessageType = "IMAGE"
	MessageTypeVideo MessageType = "VIDEO"
	MessageTypeGif   MessageType = "GIF"
	MessageTypeAudio MessageType = "AUDIO"
)

func (e MessageType) IsValid() bool {
	switch e {
	case MessageTypeText, MessageTypeImage, MessageTypeVideo, MessageTypeGif, MessageTypeAudio:
		return true
	}
	return false
}

func (e MessageType) String() string {
	return string(e)
}

func (e *MessageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageType", str)
	}
	return nil
}

func (e MessageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type State string

const (
	StateSend   State = "SEND"
	StateUnread State = "UNREAD"
	StateRead   State = "READ"
	StateDelete State = "DELETE"
	StateAdd    State = "ADD"
)

func (e State) IsValid() bool {
	switch e {
	case StateSend, StateUnread, StateRead, StateDelete, StateAdd:
		return true
	}
	return false
}

func (e State) String() string {
	return string(e)
}

func (e *State) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = State(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid State", str)
	}
	return nil
}

func (e State) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Lets redefine the base ID type to use an id from an external library
func MarshalID(id int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (int, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	return strconv.Atoi(id)
}
