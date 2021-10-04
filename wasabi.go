package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type Wasabi struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Bucket          string
}

func (w *Wasabi) upload(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", err
	}

	sess, err := w.connectAws()
	if err != nil {
		return "", err
	}

	up, err := s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Bucket: aws.String(w.Bucket),
		//ACL:    aws.String("public-read"),
		Key:  aws.String(filename),
		Body: file,
	})

	if err != nil {
		return "", err
	}

	return up.Location, nil
}

func (w *Wasabi) connectAws() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Endpoint: &w.Endpoint,
			Region:   aws.String(w.Region),
			Credentials: credentials.NewStaticCredentials(
				w.AccessKeyID,
				w.SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	return sess, err
}
