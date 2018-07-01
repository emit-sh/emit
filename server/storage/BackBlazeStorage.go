package storage

import (
	"github.com/kurin/blazer/b2"
	"io"
	"context"
)

type BackBlazeStorage struct {
	Storage
	Client *b2.Client
}

func (b BackBlazeStorage) Get(path string) (reader io.ReadCloser, err error) {
	ctx := context.Background()
	bucket, err := b.getBucket("sharing-test-bucket")
	reader = bucket.Object(path).NewReader(ctx)
	return
}

func (b BackBlazeStorage) Put(path string, reader io.Reader, ttl int) error {

	ctx := context.Background()
	bucket, err := b.getBucket("sharing-test-bucket")
	if err != nil {
		return err
	}
	obj := bucket.Object(path)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, reader); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

//TODO implement this
func (b BackBlazeStorage) GetDownloadUrl(path string) (url string, err error) {
	return path,nil
}

func (b BackBlazeStorage) Exists(path string) (exists bool, err error) {
	return true,nil
}

func (b BackBlazeStorage) Type() StorageType {
	return BackBlaze
}

func (b BackBlazeStorage) getBucket(basePath string) (bucket *b2.Bucket, err error) {

	//uri, err := url.Parse(basePath)
	if err != nil {
		return nil, err
	}

	att := b2.BucketAttrs{}

	att.Type = b2.BucketType("allPublic")

	life := b2.LifecycleRule{Prefix:"*",DaysNewUntilHidden:5,DaysHiddenUntilDeleted:1}

	att.LifecycleRules = append(att.LifecycleRules, life)

	ctx := context.Background()
	bucket, err = b.Client.NewBucket(ctx,basePath,&att)
	if err == nil {
		return
	}


	ctx = context.Background()
	bucket, err = b.Client.Bucket(ctx, basePath)
	if err != nil {
		return nil, err
	}
	return
}