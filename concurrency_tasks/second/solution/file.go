package file

import (
	"fmt"
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
	for {
		if source.mu.TryLock() {
			time.Sleep(time.Millisecond) //имитация задержки при доступе к файлу
			if target.mu.TryLock() {
				break // обе блокировки захвачены, можно копировать
			}
			source.mu.Unlock() // снимаем блокировку source, если блокировку taget получить не удалось
		}
		time.Sleep(time.Millisecond) // следующая попытка через 1мс
	}
	// в этой точке получены обе блокировки
	copy(target.data, source.data)
	target.mu.Unlock()
	source.mu.Unlock()
}
