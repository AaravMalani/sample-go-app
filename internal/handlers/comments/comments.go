package comments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AaravMalani/cvwo-forum-backend/internal/api"
	"github.com/AaravMalani/cvwo-forum-backend/internal/dataaccess/comments"
	"github.com/AaravMalani/cvwo-forum-backend/internal/dataaccess/topics"
	"github.com/pkg/errors"
)

const (
	DefaultNumComments = 10
	ListComments       = "comments.HandleList"
	CommentCreate      = "comments.HandleCreate"
	CommentUpdate      = "comments.HandleUpdate"
	CommentVote        = "comments.HandleVote"

	SuccessfulListCommentsMessage      = "Successfully listed comments"
	SuccessfulUpdateCommentMessage     = "Successfully updated comment"
	SuccessfulCreateCommentMessage     = "Successfully created comment"
	SuccessfulUpdateCommentVoteMessage = "Successfully voted for comment"
	ErrRetrieveDatabase                = "Failed to retrieve database in %s"
	ErrRetrieveComments                = "Failed to retrieve comments in %s"
	ErrEncodeView                      = "Failed to encode comments in %s"
	ErrVoteComment                     = "Failed to vote for comment in %s"
	ErrInvalidForm                     = "Invalid form provided"
	ErrBannedUser                      = "You are banned from this topic"
)

func HandleList(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	query := r.URL.Query()
	numComments, err := strconv.Atoi(query.Get("n"))
	userId := r.Context().Value("user")
	if err != nil {
		numComments = DefaultNumComments
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = 0
	}
	postId := query.Get("post_id")

	comments, err := comments.ListByPost(postId, numComments, offset, userId)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveComments, ListComments))
	}

	data, err := json.Marshal(comments)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListComments))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulListCommentsMessage},
	}, nil
}

func HandleCreate(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()

	if !r.PostForm.Has("description") || !r.PostForm.Has("post_id") || !r.PostForm.Has("topic_id") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	// To-do, scrape description for media links and add media IDs
	description := r.PostForm.Get("description")
	postId := r.PostForm.Get("post_id")
	topicId := r.PostForm.Get("topic_id")
	repliedCommentId := r.PostForm.Get("replied_comment_id")
	userId := r.Context().Value("user").(string)

	if topics.IsUserBanned(topicId, userId) {
		return &api.Response{
			ErrorCode: 403,
			Messages:  []string{ErrBannedUser},
		}, nil
	}

	comment, err := comments.CreateComment(userId, postId, topicId, repliedCommentId, description)
	if err != nil {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	data, err := json.Marshal(comment)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, CommentCreate))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{SuccessfulCreateCommentMessage},
	}, nil
}

func HandleEdit(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()

	if !r.PostForm.Has("description") || !r.PostForm.Has("comment_id") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	// To-do, scrape description for media links and add media IDs
	description := r.PostForm.Get("description")
	commentId := r.PostForm.Get("comment_id")
	userId := r.Context().Value("user").(string)

	err := comments.UpdateComment(userId, commentId, description)
	if err != nil {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}

	return &api.Response{
		Payload:  api.Payload{},
		Messages: []string{SuccessfulUpdateCommentMessage},
	}, nil
}

func HandleVote(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	r.ParseForm()

	if !r.PostForm.Has("comment_id") || !r.PostForm.Has("vote") {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}
	commentId := r.PostForm.Get("comment_id")
	vote, err := strconv.ParseInt(r.PostForm.Get("vote"), 10, 16)
	if err != nil || !(vote == 1 || vote == -1 || vote == 0) {
		return &api.Response{
			ErrorCode: 400,
			Messages:  []string{ErrInvalidForm},
		}, nil
	}
	userId := r.Context().Value("user").(string)
	err = comments.VoteComment(userId, commentId, int16(vote))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrVoteComment, CommentVote))
	}
	return &api.Response{
		Payload:  api.Payload{},
		Messages: []string{SuccessfulUpdateCommentVoteMessage},
	}, nil
}
