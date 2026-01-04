package topics

import (
	"context"
	"fmt"
	"time"

	"github.com/AaravMalani/cvwo-forum-backend/model"
	"github.com/AaravMalani/cvwo-forum-backend/query"
	"github.com/pkg/errors"
)

const (
	TopicsListBySize = "topics.ListBySize"

	ErrDbError = "Unable to find records in %s"
)

func ListBySize(n int, offset int) ([]model.EnumeratedTopic, error) {
	ctx := context.Background()
	t := query.Topic
	u := query.UserFeed
	var users []model.EnumeratedTopic
	err := t.WithContext(ctx).Select(t.ALL, u.ALL.Count()).LeftJoin(u, u.TopicID.EqCol(t.ID)).Group(t.ID).Scan(&users)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrDbError, TopicsListBySize))
	}
	return users, nil
}

func CreateTopic(name string, description string, userId string) (model.Topic, error) {
	ctx := context.Background()

	topic := model.Topic{
		Name:            name,
		Description:     description,
		LeadModeratorID: userId,
		CreatedAt:       time.Now(),
	}

	err := query.Q.Transaction(func(tx *query.Query) error {
		t := tx.Topic
		err := t.WithContext(ctx).Create(&topic)
		if err != nil {
			return err
		}

		m := tx.TopicModerator
		topicModerator := model.TopicModerator{
			UserID:  userId,
			TopicID: *topic.ID,
		}
		err = m.WithContext(ctx).Create(&topicModerator)
		return err
	})
	return topic, err
}
func IsUserBanned(topicId string, userId string) bool {
	ctx := context.Background()
	b := query.BannedUser

	users, err := b.WithContext(ctx).Select(b.ALL).Where(b.UserID.Eq(userId)).Where(b.TopicID.Eq(topicId)).Count()
	return err == nil && users > 0
}
