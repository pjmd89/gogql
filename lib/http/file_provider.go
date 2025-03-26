package http

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

type FileProvider struct {
	workingDir string
	sessDir    string
	lock       sync.Mutex
}

type sessionAccess struct {
	SessData       any
	Id             string
	LastAccessTime int64
}

func newFileProvider(dirPath string) sessionProvider {
	return &FileProvider{
		workingDir: dirPath,
		sessDir:    filepath.Join(dirPath, "sessions"),
		lock:       sync.Mutex{},
	}
}

func (o *FileProvider) init(sessionID string, sessionData any) (err error) {

	err = os.MkdirAll(o.sessDir, 0700)
	if err != nil {
		return
	}
	if sessionData != nil {
		err = o.Set(sessionID, sessionData)
	}
	return
}

func (o *FileProvider) fullpath(filename, fileFormat string) string {
	return filepath.Join(o.sessDir, filename+"."+fileFormat)
}

func (o *FileProvider) Get(sessionID string) (r any, err error) {

	filepath := o.fullpath(sessionID, "gob")
	_, statErr := os.Stat(filepath)
	if os.IsNotExist(statErr) {
		err = errors.New("session is not seted")
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		return
	}

	defer file.Close()
	gobenconder := gob.NewDecoder(file)
	v := sessionAccess{}
	err = gobenconder.Decode(&v)
	if err == nil {
		r = v.SessData
	}
	return
}

func (o *FileProvider) Set(sessionID string, sessionData any) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	filepath := o.fullpath(sessionID, "gob")

	if sessionData == nil || reflect.ValueOf(sessionData).IsZero() {
		return
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}

	gobenconder := gob.NewEncoder(file)
	v := sessionAccess{
		Id:             sessionID,
		LastAccessTime: time.Now().Unix(),
		SessData:       sessionData,
	}
	err = gobenconder.Encode(v)
	if err != nil {
		return
	}

	file.Close()

	return
}

func (o *FileProvider) Count() (r int, err error) {
	files, err := os.ReadDir(o.sessDir)
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

func (o *FileProvider) garbageCollector(sessMaxLifeTime int64) {
	err := os.MkdirAll(o.sessDir, 0700)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		files, _ := os.ReadDir(o.sessDir)
		for _, v := range files {
			o.lock.Lock()
			content, _ := os.ReadFile(o.sessDir + "/" + v.Name())
			o.lock.Unlock()
			data := sessionAccess{}
			err := gob.NewDecoder(bytes.NewBuffer(content)).Decode(&data)
			if err == nil && (time.Now().Unix()-data.LastAccessTime) >= sessMaxLifeTime {
				o.Destroy(data.Id)
			}
		}

	}
}

func (o *FileProvider) updateSessionAccess(sessionID string) {
	o.lock.Lock()
	defer o.lock.Unlock()

	filename := o.fullpath(sessionID, "gob")
	f, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return
	}

	defer f.Close()
	accessData := sessionAccess{}

	decoder := gob.NewDecoder(f)
	encoder := gob.NewEncoder(f)
	err = decoder.Decode(&accessData)
	if err != nil {
		return
	}

	f.Seek(0, io.SeekStart)

	accessData.LastAccessTime = time.Now().Unix()
	err = encoder.Encode(accessData)
	if err != nil {
		log.Println(err.Error())
	}
}
