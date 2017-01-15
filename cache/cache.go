package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
)

//go:generate counterfeiter . Cache
type Cache interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type DiskCache struct {
	baseDir string
}

func NewDiskCache(baseDir string) *DiskCache {
	return &DiskCache{
		baseDir: baseDir,
	}
}

func (self *DiskCache) Set(key string, data []byte) error {
	dirPath, filePath := self.getPaths(key)

	err := os.MkdirAll(dirPath, 0770)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0660)
	if err != nil {
		return err
	}

	return nil
}

func (self *DiskCache) Get(key string) ([]byte, error) {
	_, filePath := self.getPaths(key)

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (self *DiskCache) Delete(key string) error {
	_, filePath := self.getPaths(key)
	return os.Remove(filePath)
}

func (self *DiskCache) getPaths(key string) (dirPath, filePath string) {
	hexdigest := hexDigest([]byte(key))

	dirPath = path.Join(self.baseDir, hexdigest[0:4])
	filePath = path.Join(dirPath, hexdigest[4:])
	return
}

func hexDigest(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)
	return hex.EncodeToString(digest)
}
