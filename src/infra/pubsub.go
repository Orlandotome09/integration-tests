package infra

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func OpenPubSubConnection(ctx context.Context, projectID string, options ...option.ClientOption) *pubsub.Client {

	client, err := pubsub.NewClient(ctx, projectID, options...)
	if err != nil {
		logrus.Errorf("Failed do connect to google cloud project: %+v", errors.WithStack(err))
		return nil
	}

	return client
}
