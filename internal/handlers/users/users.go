package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/AaravMalani/cvwo-forum-backend/internal/api"
	"github.com/AaravMalani/cvwo-forum-backend/internal/auth"
	"github.com/AaravMalani/cvwo-forum-backend/internal/dataaccess/users"
	"github.com/AaravMalani/cvwo-forum-backend/internal/utils"
	"github.com/AaravMalani/cvwo-forum-backend/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	UserHandleLogin      = "users.HandleLogin"
	UserHandleRegister   = "users.HandleRegister"
	UserHandleModerators = "users.HandleModerators"
	UserHandleBanned     = "users.HandleBanned"

	SuccessfulHandleLoginMessage      = "Successfully logged in"
	SuccessfulHandleModeratorsMessage = "Succesfully listed moderator topics"
	SuccessfulHandleBannedMessage     = "Succesfully listed banned topics"
	ErrInvalidForm                    = "Invalid form"
	ErrInvalidCredentials             = "Invalid credentials"
	ErrUsernameEmailTaken             = "Username or email taken"

	ErrRetrieveDatabase = "Failed to retrieve database in %s"
	ErrRetrieveUser     = "Failed to retrieve user in %s"
	ErrRetrieveTopics   = "Failed to retrieve topics in %s"
	ErrCreateUser       = "Failed to create user in %s"
	ErrEncodeView       = "Failed to retrieve user in %s"

	ErrSaltCreation = "Failed to process request in %s"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()
	if (!r.PostForm.Has("username") && !r.PostForm.Has("email")) || !r.PostForm.Has("password") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	var user *model.AuthDetails
	var err error
	if r.PostForm.Has("username") {
		username := r.PostForm.Get("username")
		user, err = users.FindByUsername(username)
	} else {
		email := r.PostForm.Get("email")
		user, err = users.FindByEmail(email)
	}
	if err != nil {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidCredentials},
		}, nil
	}
	user_password := r.PostForm.Get("password")
	if !auth.ValidateLogin(user.Salt, user_password, user.Password) {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidCredentials},
		}, nil
	}

	auth.AddJWT(w, user.ID)
	return &api.Response{
		Messages: []string{SuccessfulHandleLoginMessage},
	}, nil
}

func HandleRegister(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()
	if (!r.PostForm.Has("username") || !r.PostForm.Has("email")) || !r.PostForm.Has("password") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	// TODO: Validate emails

	isValidUsername := regexp.MustCompile("^[a-zA-Z0-9_]+$").MatchString
	if !isValidUsername(r.PostForm.Get("username")) {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	salt, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrSaltCreation, UserHandleRegister))
	}

	hash := auth.GeneratePasswordHash(salt, r.PostForm.Get("password"))
	user := model.User{Username: r.PostForm.Get("username"), Email: r.PostForm.Get("email"), Password: hash, Salt: salt, CreatedAt: time.Now()}

	err = users.CreateUser(&user)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.Response{
				ErrorCode: 400,
				Messages:  []string{ErrUsernameEmailTaken},
			}, nil
		} else {
			return nil, errors.Wrap(err, fmt.Sprintf(ErrCreateUser, UserHandleRegister))
		}
	}

	auth.AddJWT(w, *user.ID)
	return &api.Response{
		Messages: []string{SuccessfulHandleLoginMessage},
	}, nil
}

func HandleModerators(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userId := r.Context().Value("user")
	topics, err := users.GetModeratedTopics(userId.(string))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveTopics, UserHandleModerators))
	}

	data, err := json.Marshal(topics)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, UserHandleModerators))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulHandleModeratorsMessage},
	}, nil
}

func HandleBanned(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	userId := r.Context().Value("user")
	topics, err := users.GetBannedTopics(userId.(string))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveTopics, UserHandleBanned))
	}

	data, err := json.Marshal(topics)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, UserHandleBanned))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulHandleBannedMessage},
	}, nil
}
