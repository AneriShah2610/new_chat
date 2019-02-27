//go:generate gorunpkg github.com/99designs/gqlgen

package handler

import (
	"github.com/aneri/new_chat/model"
	"sync"

	"github.com/aneri/new_chat/graph"
)

type Resolver struct {
	AddMessages   map[int]*model.ChatRoom
	UpdateMessage map[int]*model.ChatRoom
	DeleteMessage map[int]*model.ChatRoom
	ChatRoomList  map[int]*model.Member
	mu            sync.Mutex // nolint: structcheck
}

func NewResolver() *Resolver {
	return &Resolver{
		AddMessages:   make(map[int]*model.ChatRoom),
		UpdateMessage: make(map[int]*model.ChatRoom),
		DeleteMessage: make(map[int]*model.ChatRoom),
		ChatRoomList:  make(map[int]*model.Member),
	}
}

var g = NewResolver()

func init() {
	g = NewResolver()
}

func New() graph.Config {
	return graph.Config{
		Resolvers: &Resolver{
			AddMessages:   map[int]*model.ChatRoom{},
			UpdateMessage: map[int]*model.ChatRoom{},
			DeleteMessage: map[int]*model.ChatRoom{},
			ChatRoomList:  map[int]*model.Member{},
		},
	}
}
func (r *Resolver) ChatConversation() graph.ChatConversationResolver {
	return &chatConversationResolver{r}
}
func (r *Resolver) ChatRoom() graph.ChatRoomResolver {
	return &chatRoomResolver{r}
}
func (r *Resolver) Member() graph.MemberResolver {
	return &memberResolver{r}
}
func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type chatConversationResolver struct{ *Resolver }

type chatRoomResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

type subscriptionResolver struct{ *Resolver }

type memberResolver struct{ *Resolver }
