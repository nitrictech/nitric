package s3_plugin

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

// S3Plugin - Is the concrete implementation of AWS S3 for the Nitric Storage Plugin
type S3Plugin struct {
	sdk.UnimplementedStoragePlugin
	client s3iface.S3API
}

func (s *S3Plugin) getBucketByName(bucket string) (*s3.Bucket, error) {
	out, err := s.client.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		return nil, fmt.Errorf("Encountered an error retrieving the bucket list: %v", err)
	}

	for _, b := range out.Buckets {
		// TODO: This could be rather slow, it's interesting that they don't return this in the list buckets output
		tagout, err := s.client.GetBucketTagging(&s3.GetBucketTaggingInput{
			Bucket: b.Name,
		})

		if err != nil {
			return nil, err
		}

		for _, tag := range tagout.TagSet {
			fmt.Println(fmt.Sprintf("Checking Tags %s -> %s", *tag.Value, bucket))
			if *tag.Key == "x-nitric-name" && *tag.Value == bucket {
				return b, nil
			}
		}
	}

	return nil, fmt.Errorf("Unable to find bucket with name: %s", bucket)
}

// Put - Writes a new item to a bucket
func (s *S3Plugin) Put(bucket string, key string, object []byte) error {
	if b, err := s.getBucketByName(bucket); err == nil {
		contentType := http.DetectContentType(object)

		_, err := s.client.PutObject(&s3.PutObjectInput{
			Bucket:      b.Name,
			Body:        bytes.NewReader(object),
			ContentType: &contentType,
			Key:         aws.String(key),
		})

		if err != nil {
			return err
		}

	} else {
		return err
	}

	return nil
}

// Get - Retrieves an item from a bucket
func (s *S3Plugin) Get(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// New creates a new default S3 storage plugin
func New() (sdk.StoragePlugin, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		// FIXME: Use ENV configuration
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("Error creating new AWS session %v", sessionError)
	}

	s3Client := s3.New(sess)

	return &S3Plugin{
		client: s3Client,
	}, nil
}

// NewWithClient creates a new S3 Storage plugin and injects the given client
func NewWithClient(client s3iface.S3API) (sdk.StoragePlugin, error) {
	return &S3Plugin{
		client: client,
	}, nil
}
