package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

func downloadFromS3(src string, file *os.File) error {
	awsCfg := aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			c.AccessKey,
			c.SecretKey,
			"",
		),
		Region: aws.String(c.Region),
	}
	sess, err := session.NewSession(&awsCfg)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(sess)

	param := s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(src),
	}
	if _, err := downloader.Download(file, &param); err != nil {
		return err
	}

	return nil
}
