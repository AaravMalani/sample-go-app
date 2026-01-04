package users

import (
	"context"
	"fmt"

	"github.com/AaravMalani/cvwo-forum-backend/model"
	"github.com/AaravMalani/cvwo-forum-backend/query"
	"github.com/pkg/errors"
)

func Search(n int) ([]model.User, error) {
	return nil, nil
}

const (
	UsersFindByUsername = "users.FindByUsername"
	UsersFindByEmail    = "users.FindByEmail"

	ErrNoRecord = "Unable to find record in %s"
)

func FindByUsername(username string) (*model.AuthDetails, error) {
	ctx := context.Background()

	var user model.AuthDetails
	u := query.User
	err := u.WithContext(ctx).Select(u.ID, u.Salt, u.Password).Where(u.Username.Eq(username)).Scan(&user)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrNoRecord, UsersFindByUsername))
	}

	return &user, nil
}

func FindByEmail(email string) (*model.AuthDetails, error) {
	ctx := context.Background()

	var user model.AuthDetails
	u := query.User
	err := u.WithContext(ctx).Select(u.ID, u.Salt, u.Password).Where(u.Email.Eq(email)).Scan(&user)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrNoRecord, UsersFindByEmail))
	}

	return &user, nil
}
func CreateUser(user *model.User) error {
	ctx := context.Background()
	u := query.User
	return u.WithContext(ctx).Create(user)
}

func GetModeratedTopics(userId string) ([]string, error) {
	ctx := context.Background()
	t := query.TopicModerator
	var topics []string
	err := t.WithContext(ctx).Where(t.UserID.Eq(userId)).Pluck(t.TopicID, &topics)
	return topics, err
}

func GetBannedTopics(userId string) ([]string, error) {
	ctx := context.Background()
	b := query.BannedUser
	var topics []string
	err := b.WithContext(ctx).Where(b.UserID.Eq(userId)).Pluck(b.TopicID, &topics)
	return topics, err
}
