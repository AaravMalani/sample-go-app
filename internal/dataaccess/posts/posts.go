package posts

import (
	"context"
	"fmt"
	"time"

	"github.com/AaravMalani/cvwo-forum-backend/model"
	"github.com/AaravMalani/cvwo-forum-backend/query"
)

const (
	PostUpdate        = "posts.UpdatePost"
	ErrNoRowsAffected = "The comment is not found in %s"
)

func CreatePost(userId string, topicId string, title string, description string) (model.Post, error) {
	ctx := context.Background()
	p := query.Post

	post := model.Post{
		TopicID:     topicId,
		Title:       title,
		Description: description,
		AuthorID:    userId,
		CreatedAt:   time.Now(),
		EditedAt:    time.Now(),
	}
	err := p.WithContext(ctx).Create(&post)
	return post, err
}

func UpdatePost(userId string, postId string, description string) error {
	ctx := context.Background()
	p := query.Post
	result, err := p.WithContext(ctx).Where(p.AuthorID.Eq(userId)).Where(p.ID.Eq(postId)).Where(p.Deleted.Is(false)).Updates(map[string]interface{}{"description": description, "edited_at": time.Now()})
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected, PostUpdate)
	}
	return nil
}

func ListByTopic(topicId string, n int, lastPostTime time.Time, userId any) ([]model.VotedPost, error) {
	p := query.Post
	v := query.PostVote
	ctx := context.Background()

	if userId == nil {
		posts, err := p.WithContext(ctx).Where(p.TopicID.Eq(topicId)).Where(p.Deleted.Is(false)).Where(p.CreatedAt.Lt(lastPostTime)).Order(p.CreatedAt.Desc()).Limit(n).Find()
		if err != nil {
			return nil, err
		}

		votedPosts := []model.VotedPost{}
		for _, post := range posts {
			votedPosts = append(votedPosts, model.VotedPost{Post: *post, Vote: 0})
		}
		return votedPosts, nil
	}
	var votedPosts []model.VotedPost
	err := p.WithContext(ctx).Select(p.ALL, v.Vote).Where(p.TopicID.Eq(topicId)).Where(p.Deleted.Is(false)).Where(p.CreatedAt.Lt(lastPostTime)).Order(p.CreatedAt.Desc()).Limit(n).LeftJoin(v, v.PostID.EqCol(p.ID), v.UserID.Eq(userId.(string))).Scan(&votedPosts)

	return votedPosts, err

}
func VotePost(userId string, postId string, vote int16) error {
	ctx := context.Background()
	return query.Q.Transaction(func(tx *query.Query) error {
		v := tx.PostVote
		p := tx.Post
		downvoteDiff := 0
		upvoteDiff := 0

		postVotes, err := v.WithContext(ctx).Where(v.PostID.Eq(postId)).Where(v.UserID.Eq(userId)).Limit(1).Find()
		if err != nil {
			return err
		}
		if len(postVotes) == 1 {
			if vote == postVotes[0].Vote {
				return nil
			}
			switch postVotes[0].Vote {
			case 1:
				upvoteDiff += 1
			case -1:
				downvoteDiff += 1
			}
			downvoteDiff = int(postVotes[0].Vote)
			upvoteDiff = int(postVotes[0].Vote)

		}
		if vote == 0 {
			_, err = v.WithContext(ctx).Where(v.PostID.Eq(postId)).Where(v.UserID.Eq(userId)).Delete()
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
			_, err = v.WithContext(ctx).Where(v.PostID.Eq(postId)).Where(v.UserID.Eq(userId)).Update(v.Vote, vote)

		}
		_, err = p.WithContext(ctx).Where(p.ID.Eq(postId)).UpdateSimple(p.Upvotes.Sub(int32(upvoteDiff)), p.Downvotes.Sub(int32(downvoteDiff)))
		return err
	})
}
