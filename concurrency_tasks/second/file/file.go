package file

import (
	"sync"
	"time"
)

type File struct {
	path string
	data []byte // содержимое файла
	mu   sync.Mutex
}

func NewFile(p string) *File {
	return &File{path: p}
}

func (f *File) Write(b []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data = append(f.data, b...)
}

func (f *File) Read() []byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.data
}

func Copy(source *File, target *File) {
	source.mu.Lock()
	time.Sleep(time.Millisecond) //имитация задержки при доступе к файлу
	target.mu.Lock()
	copy(target.data, source.data)
	target.mu.Unlock()
	source.mu.Unlock()
}
