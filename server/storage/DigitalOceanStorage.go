package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"golang.org/x/oauth2"
	"io"
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

func NewDigitalOceanStorage() *S3Storage {

	bucket := os.Getenv("SPACES_BUCKET")
	region := os.Getenv("SPACES_REGION")

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		os.Exit(1)
	}

	cfg.Region = region
	s3svc := s3.New(cfg)

	client := s3svc
	uploader := s3manager.NewUploaderWithClient(client)
	downloader := s3manager.NewDownloaderWithClient(client)

	return &S3Storage{
		uploader:   uploader,
		client:     client,
		bucket:     bucket,
		downloader: downloader,
	}
}

func (d S3Storage) Put(path string, reader io.Reader, ttl int) error {

	_, err := d.uploader.Upload(&s3manager.UploadInput{
		ACL:    "public-read",
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

	reader.Read(buf.Bytes())
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
