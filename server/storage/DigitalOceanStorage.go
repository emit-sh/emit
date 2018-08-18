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

}

func (d S3Storage) Put(path string, reader io.Reader, ttl int) error {

	_, err := d.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		Body:   reader,
	})

	return err
}

func (d S3Storage) Get(path string) (reader io.ReadCloser, err error) {
	buf := &aws.WriteAtBuffer{}
	_, err = d.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return
	}

	reader = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
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
