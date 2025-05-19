package post

import (
	"strings"
	"sync"
)

// Post представляет отдельное сообщение
type Post struct {
	Username string // Имя пользователя, который опубликовал сообщение
	Text     string // Текст сообщения
}

// Хранилище постов
type PostStorage struct {
	mu        sync.RWMutex // Чтение/запись мьютекс
	PostArr   []Post       // Массив для хранения постов
	PostLimit int          // Максимальное количество постов
	Flag      bool         // Флаг для завершения
}
// NewPostStorage создает новое хранилище постов
func NewPostStorage(limit int) *PostStorage {
	return &PostStorage{
		PostLimit: limit,
		Flag:      true, 
	}
}

// AddPost добавляет пост в массив
func (ps *PostStorage) AddPost(p Post) {
	ps.mu.Lock()         // Захватываем мьютекс для записи, чтобы предотвратить конкурентный доступ к массиву postArr.
	defer ps.mu.Unlock() // Освобождаем мьютекс при выходе из функции
	// Проверяем, не превышает ли количество постов заданный лимит
	if len(ps.PostArr) < ps.PostLimit {
		ps.PostArr = append(ps.PostArr, p) // Если лимит не превышен, добавляем новый пост в массив
	}
}


// IsPostAboutGo проверяет, говорит ли сообщение о golang или gopher
func (p *Post) IsPostAboutGo() bool {
	hasGolang := strings.Contains(strings.ToLower(p.Text), "golang") // Проверка на наличие "golang"
	hasGopher := strings.Contains(strings.ToLower(p.Text), "gopher") // Проверка на наличие "gopher"
	return hasGolang || hasGopher                                    // Возвращаем true, если одно из условий выполнено
}

// messagedata содержит тестовые сообщения
var Messagedata = []Post{
	{"alice", "Just learned about golang!"},   // Сообщение о golang
	{"bob", "Working on a new project."},      // Сообщение не о golang
	{"charlie", "golang is awesome!"},         // Сообщение о golang
	{"dave", "Not related to programming."},   // Сообщение не о golang
	{"eve", "Excited for the golang meetup!"}, // Сообщение о golang
}


