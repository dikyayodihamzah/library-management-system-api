package lib

import "sync"

type fileCache struct {
	sync.RWMutex
	files map[string]bool
}

var FileCache = fileCache{
	files: make(map[string]bool),
}

func (r *fileCache) IsExists(filename string) (status bool, exists bool) {
	r.RLock()
	status, exists = r.files[filename]
	r.RUnlock()
	return status, exists
}

func (r *fileCache) SetProcessing(filename string) {
	r.Lock()
	r.files[filename] = false
	r.Unlock()
}

func (r *fileCache) SetFinished(filename string) {
	r.Lock()
	r.files[filename] = true
	r.Unlock()
}

func (r *fileCache) DelFinished(filename string) {
	r.Lock()
	delete(r.files, filename)
	r.Unlock()
}
