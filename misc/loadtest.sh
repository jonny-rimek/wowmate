#!/bin/bash
for value in {1..2000}
do
	AWS_PAGER="" aws sqs send-message --queue-url https://sqs.us-east-1.amazonaws.com/940880032268/wm-ConvertProcessingQueueE8D6E023-15M2CISQMJF1E --message-body file://misc/s3event.json

	echo $value
done
echo All done
