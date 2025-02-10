package sessions

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

type FileProvider struct {
	dirPath string
	lock    sync.Mutex
}

func newFileProvider(dirPath string) SessionProvider {
	return &FileProvider{
		dirPath: dirPath,
	}
}

func (o *FileProvider) Init(sessionID string, sessionData interface{}) (err error) {

	err = os.MkdirAll(o.dirPath, 0700)
	if err != nil {
		return
	}

	err = o.Set(sessionID, sessionData)
	return
}

func (o *FileProvider) fullpath(filename, fileFormat string) string {
	return filepath.Join(o.dirPath, filename+"."+fileFormat)
}

func (o *FileProvider) Get(sessionID string, dataReceiver any) (r interface{}, err error) {

	filepath := o.fullpath(sessionID, "gob")
	_, statErr := os.Stat(filepath)
	if os.IsNotExist(statErr) {
		err = errors.New("session is not seted")
		return
	}

	if dataReceiver == nil || reflect.ValueOf(dataReceiver).Kind() != reflect.Pointer {
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		return
	}

	gobenconder := gob.NewDecoder(file)
	err = gobenconder.Decode(dataReceiver)
	if err == nil {
		r = dataReceiver
	}

	file.Close()

	return
}

func (o *FileProvider) Set(sessionID string, sessionData interface{}) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	filepath := o.fullpath(sessionID, "gob")

	if sessionData == nil {
		return
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0600)
	if err == nil {
		defer file.Close()
		gobenconder := gob.NewEncoder(file)
		err = gobenconder.Encode(sessionData)
	}

	return
}

func (o *FileProvider) Count() (r int, err error) {
	files, err := os.ReadDir(o.dirPath)
	if err == nil {
		r = len(files)
	}
	return
}

func (o *FileProvider) Destroy(sessionID string) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	filepath := o.fullpath(sessionID, "gob")
	err = os.Remove(filepath)
	return
}
