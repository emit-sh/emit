package storage

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"os"
)

type S3Storage struct {
	Storage
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	client     *s3.S3
	bucket     string
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func NewDigitalOceanStorage() (storage *S3Storage, err error) {

	sKey := os.Getenv("SECRET_KEY")
	aKey := os.Getenv("ACCESS_KEY")

	// Initialize a client using Spaces
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(aKey, sKey, ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"), // This is counter intuitive, but it will fail with a non-AWS region name.
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	// Create a new Space
	/*
		params := &s3.CreateBucketInput{
			Bucket: aws.String("my-new-space-with-a-unique-name"),
		}

		_, err = s3Client.CreateBucket(params)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// List all Spaces in the region
		spaces, err := s3Client.ListBuckets(nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}



		for _, b := range spaces.Buckets {
			fmt.Printf("%s\n", aws.StringValue(b.Name))
		}

	*/
	// Upload a file to the Space
	/*
		object := s3.PutObjectInput{
			Body:   strings.NewReader("The contents of the file"),
			Bucket: aws.String("my-new-space-with-a-unique-name"),
			Key:    aws.String("file.ext"),
		}
		_, err = s3Client.PutObject(&object)
	*/
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	uploader := s3manager.NewUploaderWithClient(s3Client)
	downloader := s3manager.NewDownloaderWithClient(s3Client)

	bucket := os.Getenv("SPACES_BUCKET")

	return &S3Storage{
		uploader:   uploader,
		client:     s3Client,
		bucket:     bucket,
		downloader: downloader,
	}, nil

	/*

		bucket := os.Getenv("SPACES_BUCKET")
		//region := os.Getenv("SPACES_REGION")

		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			return nil,err
		}


		defaultResolver := endpoints.NewDefaultResolver()

		myCustomResolver := func(service, region string) (aws.Endpoint, error) {
			if service == endpoints.S3ServiceID {
				return aws.Endpoint{
					URL:           "https://nyc3.digitaloceanspaces.com",
					SigningRegion: "us-east-1",
				}, nil
			}

			return defaultResolver.ResolveEndpoint(service, region)
		}

		cfg.EndpointResolver = aws.EndpointResolverFunc(myCustomResolver)
		cfg.Region = "us-east-1"
		s3svc := s3.New(cfg)


		client := s3svc
		uploader := s3manager.NewUploaderWithClient(client)
		downloader := s3manager.NewDownloaderWithClient(client)

		return &S3Storage{
			uploader:   uploader,
			client:     client,
			bucket:     bucket,
			downloader: downloader,
		}, nil
	*/
}

func (d S3Storage) Put(path string, reader io.Reader, ttl int) error {

	_, err := d.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		Body:   reader,
	})

	if err != nil {

	}
	return err
}

func (d S3Storage) Get(path string) (reader io.ReadCloser, err error) {
	buf := &aws.WriteAtBuffer{}
	_, err = d.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})

	if err != nil {

	}

	reader = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

	//reader.Read(buf.Bytes())
	return
}

// not supported yet
func (d S3Storage) GetDownloadUrl(path string) (url string, err error) {
	return path, nil
}

// not supported yet
func (d S3Storage) GetUploadUrl(path string) (url string, err error) {
	return path, nil
}

func (d S3Storage) Exists(path string) (exists bool, err error) {
	return true, nil
}

func (d S3Storage) Type() StorageType {
	return DigitalOcean
}
