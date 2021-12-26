package blobstorage

import (
	"io"
	"io/ioutil"
	"os"
)

type FSStorage struct {
	path string
}

func NewFSStorage(path string) (*FSStorage, error) {
	return &FSStorage{
		path: path,
	}, nil
}

func (st *FSStorage) Put(data io.ReadSeeker, objectName string) error {
	newFile, err := os.Create(st.path + objectName)
	if err != nil {
		return err
	}
	_, err = io.Copy(newFile, data)
	if err != nil {
		return err
	}
	newFile.Sync()
	newFile.Close()
	return nil
}

func (st *FSStorage) Get(objectName string) ([]byte, error) {
	fileData, err := ioutil.ReadFile(st.path + objectName)
	if err != nil {
		return nil, err
	}
	// _, err = io.Copy(newFile, data)
	// if err != nil {
	// 	return err
	// }
	// newFile.Sync()
	// newFile.Close()
	return fileData, nil
}
