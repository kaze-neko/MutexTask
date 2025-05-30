package sensor

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// датчик, способный замерять некоторые значения и по запросу отдавать их сторонним системам
type Sensor struct {
	mu       sync.Mutex
	value    int
	isActive bool           // флаг активности: true - датчик включён, false - выключен
	wg       sync.WaitGroup // механизм синхронизации включения и выключения
	log      []string       // лог событий датчика
}

// получение значения датчика
func (s *Sensor) GetValue() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return returnWithDelay(s.value)
}

// получение статуса активности датчика
func (s *Sensor) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isActive
}

// получени логов датчика
func (s *Sensor) GetLog() []string {
	return s.log
}

// включение датчика (false, если датчик уже активен)
func (s *Sensor) On() bool {
	if s.IsActive() {
		return false
	}
	s.mu.Lock()
	s.isActive = true
	s.log = append(s.log, "Датчик включен!")
	s.wg.Add(1)
	s.mu.Unlock()
	// раз в 1 мс датчик производит замер
	go func() {
		defer s.wg.Done()
		for {
			if !s.isActive {
				s.log = append(s.log, "Датчик выключен!")
				break
			}
			time.Sleep(time.Millisecond) // задержка между замерами
			s.mu.Lock()
			s.value = rand.Intn(42) // имитация замера
			s.log = append(s.log, "Замер произведён, новое значение "+strconv.Itoa(s.value))
			s.mu.Unlock()
		}
	}()
	return true
}

// выключение датчика (false, если датчик уже не активен)
func (s *Sensor) Off() bool {
	if !s.IsActive() {
		return false
	}
	s.mu.Lock()
	s.isActive = false
	s.mu.Unlock()
	s.wg.Wait()
	return true
}

// функция имитация задержки при чтении значения (НЕ МЕНЯТЬ!)
func returnWithDelay[T any](value T) T {
	time.Sleep(time.Nanosecond * 100000)
	return value
}
