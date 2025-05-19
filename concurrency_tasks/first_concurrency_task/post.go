package main

import (
	"strings"
	"sync"
)

// Post представляет отдельное сообщение
type Post struct {
	Username string // Имя пользователя, который опубликовал сообщение
	Text     string // Текст сообщения
}

//Хранилище постов нужно добавить (массив, mutax, flag)

// Глобальные переменные (взять и перенсти)
var (
	postArr   []Post         // Массив для хранения постов
	wg        sync.WaitGroup // Групповой мьютекс для управления состоянием горутин (переименовать/исправить)
	postLimit = 5000         // Максимальное количество постов (перенести в main.go)
	flag      = true         // Флаг для завершения
	mu        sync.RWMutex   // Чтение/запись мьютекс
)

// addPost добавляет пост в массив
func addPost(p Post) {
	mu.Lock()         // Захватываем мьютекс для записи, чтобы предотвратить конкурентный доступ к массиву postArr.
	defer mu.Unlock() // Освобождаем мьютекс при выходе из функции
	// Проверяем, не превышает ли количество постов заданный лимит
	if len(postArr) < postLimit {
		postArr = append(postArr, p) // Если лимит не превышен, добавляем новый пост в массив
	}
}

// IsTalkingAboutGo проверяет, говорит ли сообщение о golang или gopher (переименуй IsPostAboutGo)
func (p *Post) IsTalkingAboutGo() bool {
	hasGolang := strings.Contains(strings.ToLower(p.Text), "golang") // Проверка на наличие "golang"
	hasGopher := strings.Contains(strings.ToLower(p.Text), "gopher") // Проверка на наличие "gopher"
	return hasGolang || hasGopher                                    // Возвращаем true, если одно из условий выполнено
}

// messagedata содержит тестовые сообщения
var messagedata = []Post{
	{"alice", "Just learned about golang!"},   // Сообщение о golang
	{"bob", "Working on a new project."},      // Сообщение не о golang
	{"charlie", "golang is awesome!"},         // Сообщение о golang
	{"dave", "Not related to programming."},   // Сообщение не о golang
	{"eve", "Excited for the golang meetup!"}, // Сообщение о golang
}