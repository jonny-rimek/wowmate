#!/bin/bash
for value in {1..750} #run 20 times in parallel to get to 15k events
do
	AWS_PAGER="" aws sqs send-message --queue-url https://sqs.us-east-1.amazonaws.com/461497339039/wm-dev-ConvertProcessingQueueE8D6E023-17QQELE95GJC8 --message-body file://misc/s3event.json

	echo "$value"
done
echo "All done"
