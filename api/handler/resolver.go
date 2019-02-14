//go:generate gorunpkg github.com/99designs/gqlgen

package handler

import (
	"sync"

	graph "github.com/aneri/new_chat/graph"
)

type Resolver struct {
	mu sync.Mutex
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
