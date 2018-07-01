package storage

import (
	"io"
	"os"
)


type LocalFileStorage struct {
	Storage
	BasePath string
}

func New(path string) LocalFileStorage {
	return LocalFileStorage{BasePath:path}
}

func (store LocalFileStorage) Get(path string) (reader io.ReadCloser, err error) {
	reader, err = os.Open(store.BasePath + path)
	return
}

func (store LocalFileStorage) GetDownloadUrl(path string) (url string, err error) {
	_, err = os.Stat(store.BasePath + path)
	return path, nil
}

func (store LocalFileStorage) Put(path string, reader io.Reader, ttl int) (err error) {
	file, err := os.Create(store.BasePath + path);
	if err != nil {
		return
	}
	defer file.Close()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {

		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := file.Write(buf[:n]); err != nil {
			return err
		}
	}

	return
}

func (store LocalFileStorage) Exists(path string) (exists bool, err error) {
	exists = true
	return
}

func (store LocalFileStorage) Type() StorageType {
	return LocalFileSystem
}