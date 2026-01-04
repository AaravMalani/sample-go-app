package topics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AaravMalani/cvwo-forum-backend/internal/api"
	"github.com/AaravMalani/cvwo-forum-backend/internal/dataaccess/topics"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	DefaultNumTopics = 20
	ListTopics       = "topics.HandleList"
	CreateTopic      = "topics.HandleCreate"

	SuccessfulListTopicsMessage  = "Successfully listed topics"
	SuccessfulCreateTopicMessage = "Successfully created topics"
	ErrRetrieveDatabase          = "Failed to retrieve database in %s"
	ErrRetrieveTopics            = "Failed to retrieve topics in %s"
	ErrEncodeView                = "Failed to encode topics in %s"
	ErrCreateTopic               = "Failed to create topic in %s"
	ErrTopicExists               = "A topic already exists with this name"
	ErrInvalidForm               = "Invalid form provided"
)

func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	query := r.URL.Query()
	numTopics, err := strconv.Atoi(query.Get("n"))
	if err != nil {
		numTopics = DefaultNumTopics
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = 0
	}

	topics, err := topics.ListBySize(numTopics, offset)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveTopics, ListTopics))
	}

	data, err := json.Marshal(topics)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListTopics))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListTopicsMessage},
	}, nil
}

func HandleCreate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()
	if !r.PostForm.Has("name") || !r.PostForm.Has("description") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	userId := r.Context().Value("user").(string)

	topic, err := topics.CreateTopic(name, description, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.Response{
				ErrorCode: 400,
				Messages:  []string{ErrTopicExists},
			}, nil
		} else {
			return nil, errors.Wrap(err, fmt.Sprintf(ErrCreateTopic, CreateTopic))
		}
	}

	data, err := json.Marshal(topic)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, CreateTopic))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulCreateTopicMessage},
	}, nil
}
