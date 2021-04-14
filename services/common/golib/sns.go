package golib

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sirupsen/logrus"
)

// SNSPublishMsg publishes a message to an SNS topic
func SNSPublishMsg(ctx aws.Context, snsSvc *sns.SNS, input string, topicArn *string) error {
	if input == "" {
		return fmt.Errorf("input can't be empty")
	}
	logrus.Debug("sns input to publish: ", input)

	var err error

	if os.Getenv("LOCAL") == "true" {
		_, err = snsSvc.Publish(&sns.PublishInput{
			Message:  aws.String(input),
			TopicArn: topicArn,
		})
	} else {
		_, err = snsSvc.PublishWithContext(ctx, &sns.PublishInput{
			Message:  aws.String(input),
			TopicArn: topicArn,
		})
	}
	if err != nil {
		return fmt.Errorf("failed publishing a message to sns: %v", err)
	}

	logrus.Debug("message successfully sent to topic")
	return nil
}

/*
func SNSPublishMsg2(client *sns2.Client, input string, topicArn *string) error {
	if input == "" {
		return fmt.Errorf("combatlog_uuid can't be empty")
	}
	logrus.Debug("sns2 input to publish: ", input)

	_, err := client.Publish(context.TODO(), &sns2.PublishInput{
		Message:  aws2.String(input),
		TopicArn: topicArn,
	})
	if err != nil {
		return err
	}

	logrus.Debug("message successfully sent to topic")
	return nil
}
*/
