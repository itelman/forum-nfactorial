package flash

import (
	"errors"
	"sync"
)

const (
	flash = "flash"
)

var (
	ErrDataNotFound = errors.New("flash: data not found")
)

type FlashManager interface {
	UpdateFlash(val string)
	PopFlash() string
}

type flashManager struct {
	store map[string]string
	mutex sync.RWMutex
}

func NewFlashManager() *flashManager {
	return &flashManager{
		store: make(map[string]string),
	}
}

func (s *flashManager) getFlash() (string, error) {
	s.mutex.RLock() // Read lock allows multiple readers
	defer s.mutex.RUnlock()

	val, exists := s.store[flash]
	if !exists {
		return "", ErrDataNotFound
	}

	return val, nil
}

func (s *flashManager) UpdateFlash(val string) {
	s.mutex.Lock() // Locking for exclusive write access
	defer s.mutex.Unlock()

	s.store[flash] = val
}

func (s *flashManager) deleteFlash() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, flash)
}

func (s *flashManager) PopFlash() string {
	flash, err := s.getFlash()
	if err != nil {
		return ""
	}

	s.deleteFlash()

	return flash
}
