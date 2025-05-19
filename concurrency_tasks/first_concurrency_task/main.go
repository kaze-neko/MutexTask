package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"first_concurrency_task/post"
)

var wg sync.WaitGroup // wg используется для ожидания завершения всех горутин
const PostLimit = 5000 // Максимальное количество постов в хранилище

// shuffle перемешивает срез (slice) постов с использованием алгоритма Фишера-Йетса
func shuffle(posts []post.Post) {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел
	for i := len(posts) - 1; i > 0; i-- {
		j := rand.Intn(i + 1) // Генерация случайного индекса
		posts[i], posts[j] = posts[j], posts[i] // Обмен элементов
	}
}

func TestAddPostsConcurrent() { // Проверка добавления постов из 1000 горутин
	// Инициализация хранилища постов
	postStorage := post.NewPostStorage(PostLimit)

	// Перемешиваем массив сообщений
	shuffle(post.Messagedata)

	// Запускаем 1000 горутин для добавления постов
	for i := 0; i < 1000; i++ {
		wg.Add(1) // Увеличиваем счётчик ожидания
		go func() {
			defer wg.Done()                                              // Уменьшаем счётчик после завершения горутины
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) // Случайная задержка перед добавлением постов
			for j := 0; j < 5; j++ {                                     // Добавляем 5 постов из списка
				index := j % len(post.Messagedata) // Используем индекс j для доступа к перемешанному массиву
				postStorage.AddPost(post.Messagedata[index]) // Добавление поста в массив
			}
		}()
	}
	// Ожидаем завершения всех горутин добавления постов
	wg.Wait() 
	// Проверка количества постов
	if len(postStorage.PostArr) == PostLimit {
		fmt.Printf("Успех! Количество постов: %d\n", len(postStorage.PostArr))
	} else {
		fmt.Printf("Ошибка: Неверное количество постов: %d\n", len(postStorage.PostArr))
	}
}

func TestReadPostsAboutGo() { //Проверка подсчета постов, обсуждающих Go.
	// Инициализация хранилища постов
	postStorage := post.NewPostStorage(PostLimit)

	var wg1 sync.WaitGroup // Для добавления постов
	var wg2 sync.WaitGroup // Для чтения постов

	// Запускаем читателей
	var mu sync.Mutex        // Мьютекс для защиты currentPost и cnt
	currentPost := 0         // Переменная для отслеживания текущего поста
	cnt := 0                 // Счётчик постов о Go

	// Запускаем горутину добавления постов
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for i := 0; i < 1000; i++ {
			// Симуляция рандомного sleep и добавления 5 постов
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			for j := 0; j < 5; j++ {
				index := j % len(post.Messagedata) // Используем индекс j для доступа к перемешанному массиву
				postStorage.AddPost(post.Messagedata[index]) // Добавление поста в массив
			}
		}
	}()

	// Запускаем горутину чтения постов
	wg2.Add(1)
	go func() {
		defer wg2.Done() // Уменьшаем счётчик после завершения горутины
		for {
			mu.Lock() // Захватываем мьютекс для безопасного доступа к общим данным
			flag := postStorage.Flag // Получаем значение флага завершения
			length := len(postStorage.PostArr) // Получаем длину массива постов
			pos := currentPost // Сохраняем текущую позицию поста
			mu.Unlock() // Освобождаем мьютекс

			// Проверяем условие завершения: если флаг == false и текущая позиция равна длине массива
			if !flag && pos == length {
				break // Завершаем цикл, если нет новых постов
			}

			// Проверяем, если текущая позиция меньше длины массива
			if pos < length {
				mu.Lock() // Захватываем мьютекс для безопасного доступа к посту
				post := postStorage.PostArr[pos] // Получаем текущий пост
				currentPost++ // Увеличиваем текущую позицию
				mu.Unlock() // Освобождаем мьютекс

				// Проверяем, обсуждает ли пост Go
				if post.IsPostAboutGo() {
					mu.Lock() // Захватываем мьютекс для безопасного увеличения счётчика
					cnt++ // Увеличиваем счётчик постов о Go
					mu.Unlock() // Освобождаем мьютекс
				}
			} else {
				// Если текущая позиция больше или равна длине массива, делаем небольшую задержку
				time.Sleep(10 * time.Millisecond) // Задержка перед следующей проверкой
			}
		}
	}()


	wg1.Wait() // Ждем завершения добавления постов
	postStorage.Flag = false // Устанавливаем флаг завершения для читателя
	wg2.Wait()   // Ждем завершения всех горутин чтения постов


	// Проверка количества постов о golang
	if cnt == 3000 {
		fmt.Printf("Количество постов о golang: %d\n", cnt)
	} else {
		fmt.Printf("Количество постов о golang не соответствуют: %d\n", cnt)
	}
}

func main() {
	TestAddPostsConcurrent()
	TestReadPostsAboutGo()
}

