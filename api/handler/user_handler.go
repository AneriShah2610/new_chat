package handler

import (
	"context"
	"database/sql"
	er "github.com/aneri/new_chat/error"
	"math/rand"
	"time"

	"github.com/aneri/new_chat/api/dal"
	"github.com/aneri/new_chat/api/helper"

	"github.com/aneri/new_chat/model"
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
	isUserExist, err := CheckUserExistence(ctx, name)
	if err != nil {
		er.DebugPrintf(err)
		return []model.User{}, er.InternalServerError
	}
	if isUserExist{
		rows, err := crConn.Db.Query("SELECT id, username, first_name, last_name, email, contact, bio, profile_picture, created_at, updated_at FROM users WHERE username != $1 ORDER BY username ASC", name)
		if err != nil {
			er.DebugPrintf(err)
			return []model.User{}, er.InternalServerError
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Contact, &user.Bio, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				er.DebugPrintf(err)
				return []model.User{}, er.InternalServerError
			}
			users = append(users, user)
		}
	}
	return users, nil
}

// Create New User
func (r *mutationResolver) NewUser(ctx context.Context, input model.NewUser) (model.User, error) {
	var user model.User
	user = model.User{
		Username:       input.UserName,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Contact:        input.Contact,
		Bio:            input.Bio,
		ProfilePicture: input.ProfilePicture,
	}
	var err error
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	isUserExist, err := CheckUserExistence(ctx, input.UserName)
	if err != nil {
		er.DebugPrintf(err)
		return model.User{}, er.InternalServerError
	}
	if !isUserExist {
		row := crConn.Db.QueryRow("INSERT INTO users (username, first_name, last_name, email, contact, bio, profile_picture, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", input.UserName, input.FirstName, input.LastName, input.Email, input.Contact, input.Bio, input.ProfilePicture, time.Now().UTC())
		err = row.Scan(&user.ID)
		if err != nil {
			er.DebugPrintf(err)
			return model.User{}, er.InternalServerError
		}
	}

	for _, observer := range addUserChannel {
		observer <- user
	}
	return user, nil
}
func (r *queryResolver) MemberLogIn(ctx context.Context, name string) (*model.User, error){
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var user model.User
	isUserExist, err := CheckUserExistence(ctx, name)
	if err != nil {
		er.DebugPrintf(err)
		return &model.User{}, er.InternalServerError
	}
	if isUserExist{
		row := crConn.Db.QueryRow("SELECT id, username, first_name, last_name, email, contact, bio, profile_picture, created_at, updated_at FROM users WHERE username = $1", name)
		err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Contact, &user.Bio, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)
		if err != nil{
			er.DebugPrintf(err)
			return &model.User{}, er.InternalServerError
		}
	} else {
		return nil, er.UserDoesNotExists
	}

	return &user, nil
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

func CheckUserExistence(ctx context.Context, userName string) (bool, error) {
	crConn := ctx.Value("crConn").(*dal.DbConnection)
	var isUserExist bool
	row := crConn.Db.QueryRow("SELECT true FROM users WHERE username = $1", userName)

	err := row.Scan(&isUserExist)
	if err != nil && err != sql.ErrNoRows {
		er.DebugPrintf(err)
		return false, er.InternalServerError
	}
	return isUserExist, nil
}