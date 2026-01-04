package comments

import (
	"context"
	"fmt"
	"time"

	"github.com/AaravMalani/cvwo-forum-backend/model"
	"github.com/AaravMalani/cvwo-forum-backend/query"
	"github.com/pkg/errors"
)

const (
	CommentCreate = "comments.CreateComment"
	CommentUpdate = "comments.UpdateComment"

	ErrPostNotFound   = "The post is not found"
	ErrPostDeleted    = "The post is not found"
	ErrNoRowsAffected = "The comment is not found in %s"
)

func CreateComment(userId string, postId string, topicId string, repliedCommentId string, description string) (model.Comment, error) {
	ctx := context.Background()
	c := query.Comment
	p := query.Post

	post, err := p.WithContext(ctx).Select(p.Deleted, p.TopicID).Where(p.ID.Eq(postId)).First()
	if err != nil || topicId != post.TopicID {
		return model.Comment{}, errors.Wrap(err, ErrPostNotFound)
	}
	if post.Deleted {
		return model.Comment{}, errors.New(ErrPostDeleted)
	}

	comment := model.Comment{
		TopicID:     topicId,
		PostID:      postId,
		Description: description,
		AuthorID:    userId,
		CreatedAt:   time.Now(),
		EditedAt:    time.Now(),
	}
	if repliedCommentId != "" {
		comment.RepliedCommentID = repliedCommentId
	}
	err = c.WithContext(ctx).Create(&comment)
	return comment, err
}

func UpdateComment(userId string, commentId string, description string) error {

	ctx := context.Background()
	c := query.Comment
	result, err := c.WithContext(ctx).Where(c.AuthorID.Eq(userId)).Where(c.ID.Eq(commentId)).Where(c.Deleted.Is(false)).Updates(map[string]interface{}{"description": description, "edited_at": time.Now()})
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected, CommentUpdate)
	}
	return nil
}

func ListByPost(postId string, n int, offset int, userId any) ([]model.VotedComment, error) {
	c := query.Comment
	v := query.CommentVote
	ctx := context.Background()

	if userId == nil {
		comments, err := c.WithContext(ctx).Where(c.PostID.Eq(postId)).Where(c.Deleted.Is(false)).Order(c.CreatedAt).Limit(n).Offset(offset).Find()
		if err != nil {
			return nil, err
		}

		votedComments := []model.VotedComment{}
		for _, comment := range comments {
			votedComments = append(votedComments, model.VotedComment{Comment: *comment, Vote: 0})
		}
		return votedComments, nil
	}
	var votedComments []model.VotedComment
	err := c.WithContext(ctx).Select(c.ALL, v.Vote).Where(c.PostID.Eq(postId)).Where(c.Deleted.Is(false)).Order(c.CreatedAt).Limit(n).Offset(offset).LeftJoin(v, v.CommentID.EqCol(c.PostID), v.UserID.Eq(userId.(string))).Scan(&votedComments)
	return votedComments, err
}

func VoteComment(userId string, commentId string, vote int16) error {
	ctx := context.Background()
	return query.Q.Transaction(func(tx *query.Query) error {
		v := tx.CommentVote
		c := tx.Comment
		downvoteDiff := 0
		upvoteDiff := 0

		commentVotes, err := v.WithContext(ctx).Where(v.CommentID.Eq(commentId)).Where(v.UserID.Eq(userId)).Limit(1).Find()
		if err != nil {
			return err
		}
		if len(commentVotes) == 1 {
			if vote == commentVotes[0].Vote {
				return nil
			}
			switch commentVotes[0].Vote {
			case 1:
				upvoteDiff += 1
			case -1:
				downvoteDiff += 1
			}
			downvoteDiff = int(commentVotes[0].Vote)
			upvoteDiff = int(commentVotes[0].Vote)

		}
		if vote == 0 {
			_, err = v.WithContext(ctx).Where(v.CommentID.Eq(commentId)).Where(v.UserID.Eq(userId)).Delete()
			if err != nil {
				return err
			}
		} else {
			switch vote {
			case 1:
				upvoteDiff += -1
			case -1:
				downvoteDiff += -1

			}
			_, err = v.WithContext(ctx).Where(v.CommentID.Eq(commentId)).Where(v.UserID.Eq(userId)).Update(v.Vote, vote)

		}
		_, err = c.WithContext(ctx).Where(c.ID.Eq(commentId)).UpdateSimple(c.Upvotes.Sub(int32(upvoteDiff)), c.Downvotes.Sub(int32(downvoteDiff)))
		return err
	})

}
