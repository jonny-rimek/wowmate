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
		return nil
		// publishing encrypted messages to SNS doesn't work from SAM+CDK, I suspect that SAM uses an incomplete name
		// CDK, auto generates a name, but those names aren't used locally
		// e.g. wm-dev-DynamoDBtableF8E87752-HSV525WR7KN3 is the name of the ddb in the cloud
		// locally it the name it knows is wm-dev-DynamoDBtableF8E87752 the last bit is missing
		// the same is gonna be the problem for the KMS key, and I don't know how or if I can pass in the complete key

		// _, err = snsSvc.Publish(&sns.PublishInput{
		// 	Message:  aws.String(input),
		// 	TopicArn: topicArn,
		// })
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
