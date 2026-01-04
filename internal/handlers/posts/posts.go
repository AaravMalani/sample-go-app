package posts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AaravMalani/cvwo-forum-backend/internal/api"
	"github.com/AaravMalani/cvwo-forum-backend/internal/dataaccess/posts"
	"github.com/AaravMalani/cvwo-forum-backend/internal/dataaccess/topics"
	"github.com/pkg/errors"
)

const (
	DefaultNumPosts                 = 10
	ListPosts                       = "posts.HandleList"
	PostCreate                      = "posts.HandleCreate"
	PostVote                        = "posts.HandleVote"
	SuccessfulListPostsMessage      = "Successfully listed posts"
	SuccessfulCreatePostMessage     = "Successfully created post"
	SuccessfulUpdatePostMessage     = "Successfully updated post"
	SuccessfulUpdatePostVoteMessage = "Successfully voted for post"
	ErrRetrieveDatabase             = "Failed to retrieve database in %s"
	ErrRetrievePosts                = "Failed to retrieve posts in %s"
	ErrEncodeView                   = "Failed to encode posts in %s"
	ErrInvalidForm                  = "Invalid form provided"
	ErrVotePost                     = "Failed to vote for post in %s"
	ErrBannedUser                   = "You are banned from this topic"
)

func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	query := r.URL.Query()
	numPosts, err := strconv.Atoi(query.Get("n"))
	userId := r.Context().Value("user")
	if err != nil {
		numPosts = DefaultNumPosts
	}
	lastPostTime, err := strconv.ParseInt(query.Get("last_post_time"), 10, 64)
	if err != nil {
		lastPostTime = time.Now().Unix()
	}
	topicId := query.Get("topic_id")

	posts, err := posts.ListByTopic(topicId, numPosts, time.Unix(lastPostTime, 0), userId)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrievePosts, ListPosts))
	}

	data, err := json.Marshal(posts)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListPosts))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListPostsMessage},
	}, nil
}

func HandleCreate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()

	if !r.PostForm.Has("title") || !r.PostForm.Has("description") || !r.PostForm.Has("topic_id") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	// To-do, scrape description for media links and add media IDs
	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	topicId := r.PostForm.Get("topic_id")

	userId := r.Context().Value("user").(string)

	if topics.IsUserBanned(topicId, userId) {
		return &api.Response{
			ErrorCode: 403,
			Messages:  []string{ErrBannedUser},
		}, nil
	}

	post, err := posts.CreatePost(userId, topicId, title, description)
	if err != nil {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	data, err := json.Marshal(post)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, PostCreate))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulCreatePostMessage},
	}, nil
}

func HandleEdit(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()

	if !r.PostForm.Has("description") || !r.PostForm.Has("post_id") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	// To-do, scrape description for media links and add media IDs
	description := r.PostForm.Get("description")
	postId := r.PostForm.Get("post_id")
	userId := r.Context().Value("user").(string)

	err := posts.UpdatePost(userId, postId, description)
	if err != nil {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	return &api.Response{
		Payload:  api.Payload{},
		Messages: []string{SuccessfulUpdatePostMessage},
	}, nil
}
func HandleVote(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()

	if !r.PostForm.Has("post_id") || !r.PostForm.Has("vote") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}
	postId := r.PostForm.Get("post_id")
	vote, err := strconv.ParseInt(r.PostForm.Get("vote"), 10, 16)
	if err != nil || !(vote == 1 || vote == -1 || vote == 0) {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}
	userId := r.Context().Value("user").(string)
	err = posts.VotePost(userId, postId, int16(vote))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrVotePost, PostVote))
	}
	return &api.Response{
		Payload:  api.Payload{},
		Messages: []string{SuccessfulUpdatePostVoteMessage},
	}, nil
}
