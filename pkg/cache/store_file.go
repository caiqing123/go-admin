package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	files "api/pkg/file"
	"api/pkg/logger"
)

type (
	file struct {
		dir string
	}

	fileContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

const perm = 0o666

// New creates an instance of File cache
func New(dir string) Store {
	err := files.IsNotExistMkDir(dir)
	logger.LogIf(err)
	return &file{dir}
}

func (f *file) Increment(parameters ...interface{}) {
}

func (f *file) Decrement(parameters ...interface{}) {
}

func (f *file) IsAlive() error {
	return nil
}

func (f *file) createName(key string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	return filepath.Join(f.dir, fmt.Sprintf("%s.cachego", hash))
}

func (f *file) read(key string) (*fileContent, error) {
	value, err := ioutil.ReadFile(f.createName(key))
	if err != nil {
		return nil, err
	}

	content := &fileContent{}
	if err := json.Unmarshal(value, content); err != nil {
		return nil, err
	}

	if content.Duration == 0 {
		return content, nil
	}

	if content.Duration <= time.Now().Unix() {
		f.Forget(key)
		return nil, errors.New("cache expired")
	}

	return content, nil
}

// Has checks if the cached key exists into the File storage
func (f *file) Has(key string) bool {
	_, err := f.read(key)
	return err == nil
}

// Forget the cached key from File storage
func (f *file) Forget(key string) {
	_, err := os.Stat(f.createName(key))
	if err == nil && !os.IsNotExist(err) {
		_ = os.Remove(f.createName(key))
	}
}

// Get retrieves the cached value from key of the File storage
func (f *file) Get(key string) string {
	content, err := f.read(key)
	if err != nil {
		return ""
	}
	return content.Data
}

// FetchMulti retrieve multiple cached values from keys of the File storage
func (f *file) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		result[key] = f.Get(key)
	}

	return result
}

// Flush removes all cached keys of the File storage
func (f *file) Flush() {
	dir, _ := os.Open(f.dir)

	defer func() {
		_ = dir.Close()
	}()

	names, _ := dir.Readdirnames(-1)

	for _, name := range names {
		_ = os.Remove(filepath.Join(f.dir, name))
	}
}

// Set a value in File storage by key
func (f *file) Set(key string, value string, lifeTime time.Duration) {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &fileContent{duration, value}

	data, _ := json.Marshal(content)

	_ = ioutil.WriteFile(f.createName(key), data, perm)
}

// Forever 清空
func (f *file) Forever(key string, value string) {
	duration := int64(0)

	content := &fileContent{duration, value}

	data, _ := json.Marshal(content)

	_ = ioutil.WriteFile(f.createName(key), data, perm)
}
