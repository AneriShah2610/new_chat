package handler

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"

	model "github.com/aneri/new_chat/model"
)

var addUserChannel map[int]chan model.User

func init() {
	addUserChannel = map[int]chan model.User{}
}

// Retrieve user details
func (r *queryResolver) Users(ctx context.Context, name string) ([]model.User, error) {
	var users []model.User
	var user model.User
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	rows, err := crConn.Db.Query("SELECT id, name, email, contact, profile_picture, bio, createdat FROM user_test WHERE name != $1 ORDER BY name", name)
	if err != nil {
		log.Println("Error at 23 line of usr_handler", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Contact, &user.ProfilePicture, &user.Bio, &user.CreatedAt)
		if err != nil {
			log.Println("User data scan error")
		}
		users = append(users, user)
	}
	return users, nil
}

// Create New User
func (r *mutationResolver) NewUser(ctx context.Context, input model.NewUser) (model.User, error) {
	var user model.User

	var err error
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	user, err = UserData(ctx, input.Name)
	if err != nil {
		log.Println("Error to fetch user data")
	}
	if user.Name == "" {
		_, err := crConn.Db.Exec("INSERT INTO user_test (name, email, contact, profile_picture, bio, createdat) VALUES ($1, $2, $3, $4, $5, NOW())", input.Name, input.Email, input.Contact, input.ProfilePicture, input.Bio)
		if err != nil {
			log.Print("Error while inserting data", err)
		}
	}
	user = model.User{
		Name: input.Name,
		Email: input.Email,
		Contact: input.Contact,
		ProfilePicture: input.ProfilePicture,
		Bio: input.Bio,
	}
	for _, observer := range addUserChannel {
		observer <- user
	}
	return user, nil
}

// Listen New User Request & shown to live
func (r *subscriptionResolver) UserJoined(ctx context.Context) (<-chan model.User, error) {
	rand.Seed(time.Now().UnixNano())
	id := helper.Random(1, 10000000000000)
	userEvent := make(chan model.User, 1)
	go func() {

		<-ctx.Done()
		delete(addUserChannel, id)
	}()
	addUserChannel[id] = userEvent
	return userEvent, nil
}

func UserData(ctx context.Context, name string) (model.User, error) {
	var user model.User
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	row, err := crConn.Db.Query("SELECT id, name, email, contact, profile_picture, bio, createdat FROM user_test WHERE name = $1", name)
	if err != nil {
		log.Println("Error to read user data as per name", err)
	}
	defer row.Close()
	for row.Next() {
		err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Contact, &user.ProfilePicture, &user.Bio, &user.CreatedAt)
		if err != nil {
			log.Println("Error to read user details as per name", err)
		}
	}
	return user, nil
}
