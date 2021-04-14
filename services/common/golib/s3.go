package golib

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Event base structure for an s3 event, to unmarshal a json into it
type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

// DownloadS3 downloads a file from s3
func DownloadS3(ctx aws.Context, downloader *s3manager.Downloader, bucketName string, objectKey string, fileContent io.WriterAt) error {
	var err error

	if os.Getenv("LOCAL") == "true" {
		_, err = downloader.Download(
			fileContent,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
			},
		)
	} else {
		_, err = downloader.DownloadWithContext(
			ctx,
			fileContent,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
			},
		)
	}
	if err != nil {
		return err
	}

	return nil
}

/*
func downloadS32(cfg aws2.Config, bucketName string, objectKey string, fileContent io.WriterAt) error {
	client := s32.NewFromConfig(cfg)
	downloader := manager2.NewDownloader(client)
	downloader.Download(context.TODO(), fileContent, &s32.GetObjectInput{
		Bucket: aws2.String(bucketName),
		Key:    aws2.String(objectKey),
	})
	return nil
}
*/

// RenameFileS3 creates a copy of a file and deletes the old one, because you can't rename files =(
func RenameFileS3(ctx aws.Context, s3Svc *s3.S3, oldName, newName, bucketName string) error {
	source := fmt.Sprintf("%s/%s", bucketName, oldName)
	escapedSource := url.PathEscape(source)

	_, err := s3Svc.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Key:        &newName,
		Bucket:     &bucketName,
		CopySource: &escapedSource,
	})
	if err != nil {
		return err
	}

	_, err = s3Svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Key:    &oldName,
		Bucket: &bucketName,
	})
	if err != nil {
		return err
	}

	return nil
}

// SizeOfS3Object checks the file size without actually downloading it in KB
func SizeOfS3Object(ctx aws.Context, s3Svc *s3.S3, bucketName string, objectKey string) (int, error) {
	var output *s3.HeadObjectOutput
	var err error
	// if it's local do without context
	if os.Getenv("LOCAL") == "true" {
		output, err = s3Svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	} else {
		output, err = s3Svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	}
	if err != nil {
		return 0, fmt.Errorf("unable to to send head request to item %q, %v", objectKey, err)
	}

	return int(*output.ContentLength / 1024), nil
}

// checks the file size without actually downloading it in MB
/*
func sizeOfS3Object2(cfg aws2.Config, bucketName string, objectKey string) (int, error) {
	client := s32.NewFromConfig(cfg)

	input := s32.HeadObjectInput{
		Bucket: aws2.String(bucketName),
		Key:    aws2.String(objectKey),
	}
	output, err := client.HeadObject(context.TODO(), &input)
	if err != nil {
		return 0, fmt.Errorf("Unable to to send head request to item %q, %v", objectKey, err)
	}

	return int(output.ContentLength / 1024 / 1024), nil
}
*/
