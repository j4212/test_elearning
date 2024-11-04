package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	c "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

type R2Cloudlfare struct {
	Bucket       string
	AccountID    string
	Key          string
	Secret       string
	PubBucketUrl string
}

type R2Stub struct {
	Client       *s3.Client
	Bucket       string
	PubBucketUrl string
}

func NewR2Stub(config *R2Cloudlfare) *R2Stub {
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", config.AccountID),
			}, nil
		},
	)

	cfg, err := c.LoadDefaultConfig(
		context.Background(),
		c.WithEndpointResolverWithOptions(resolver),
		c.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(config.Key, config.Secret, ""),
		),
		c.WithRegion("auto"),
	)

	if err != nil {
		logrus.Warnln("[R2 Cloudflare] Couldn't initialize configuration because error", err)
		return nil
	}

	client := s3.NewFromConfig(cfg)

	return &R2Stub{
		Client:       client,
		Bucket:       config.Bucket,
		PubBucketUrl: config.PubBucketUrl,
	}

}

func (r2 *R2Stub) Upload(file string) (string, error) {
	f, err := os.Open(file)

	if err != nil {
		return "", err
	}

	defer f.Close()

	// Getting extension file
	ext := path.Ext(file)

	// Random Oject name
	randID, _ := gonanoid.New(15)

	object := randID

	contentType := "application/octet-stream"

	logrus.Infoln("[R2 Cloudflare] Extension:", ext)
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		extReplaced := strings.ReplaceAll(ext, ".", "")
		contentType = fmt.Sprintf("image/%s", extReplaced)
	}

	_, err = r2.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(r2.Bucket),
		Key:         aws.String(object),
		Body:        f,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	errDir := os.RemoveAll("./temp")

	if errDir != nil {
		os.Mkdir("./temp", os.ModePerm)
	}

	logrus.Infoln("[R2 Cloudflare] File has been uploaded", object)
	return fmt.Sprintf("%s/%s", r2.PubBucketUrl, object), nil

}
