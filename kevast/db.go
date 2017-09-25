package kevast

import (
	"errors"
	"fmt"
	"sync"
)

var (
	errorEmptyKey    = errors.New("Value cannot be empty")
	errorEmptyValue  = errors.New("Key cannot be empty")
	errorKeyNotExist = errors.New("Key does not exist")
	errorNoTx        = errors.New("error no transaction started")
	errorEmptyTx     = errors.New("Nothing to commit. Aborting transaction")
)

type store map[string]string

// Kevast is the base struct for the key/value store
type Kevast struct {
	idx    int64
	stores []store
	mu     sync.Mutex
}

// NewDB creates an initilized Kevast instance
func NewDB() *Kevast {
	return &Kevast{
		idx:    0,
		stores: []store{store{}},
		mu:     sync.Mutex{},
	}
}

// Write inserts the value for the given key in to the store
func (s *Kevast) Write(key string, val string) error {

	if key == "" {
		return errorEmptyKey
	}

	if val == "" {
		return errorEmptyValue
	}

	s.mu.Lock()
	s.stores[s.idx][key] = val
	s.mu.Unlock()
	return nil
}

// Read retrieves the value for the given key.
// This is done by looping through the stores checking
// for the value or if the value has been deleted in a previous
// transaction
func (s *Kevast) Read(key string) (string, error) {

	if key == "" {
		return "", errorEmptyKey
	}

	for i := s.idx; i >= 0; i-- {
		s.mu.Lock()
		val, ok := s.stores[i][key]
		s.mu.Unlock()
		if !ok {
			continue
		}
		// keys with a value of "" denotes the entry has been deleted
		// in a transaction, but hasn't been commited yet
		if val == "" && ok {
			return "", fmt.Errorf(errNotFound, key)
		}
		return val, nil
	}

	return "", errorKeyNotExist
}

// Del deletes the key and value for the given key
func (s *Kevast) Del(key string) error {
	if key == "" {
		return errorEmptyKey
	}
	_, err := s.Read(key)
	if err != nil {
		return errorKeyNotExist
	}

	if s.idx == 0 {
		s.mu.Lock()
		delete(s.stores[s.idx], key)
		s.mu.Unlock()
		return nil
	}

	s.mu.Lock()
	s.stores[s.idx][key] = ""
	s.mu.Unlock()
	return nil
}

// Start will begin a transaction and store all
// changes in a temp store until commited or aborted
func (s *Kevast) Start() error {
	s.stores = append(s.stores, store{})
	s.idx++
	return nil
}

// Abort will delete the temp store for the current transaction
// and exit the transaction
func (s *Kevast) Abort() error {
	if s.idx == 0 {
		return errorNoTx
	}
	s.clearTx()
	return nil
}

// Commit will apply changes from the current transaction temp store
// to the parent store. This may be another transaction temp store if the
// current transaction is nested
func (s *Kevast) Commit() error {
	if s.idx == 0 {
		return errorNoTx
	}

	if len(s.stores[s.idx]) == 0 {
		s.clearTx()
		return errorEmptyTx
	}

	var wg sync.WaitGroup
	for k, v := range s.stores[s.idx] {
		wg.Add(1)
		if s.idx == 1 && v == "" {
			go func(k, v string) {
				s.mu.Lock()
				delete(s.stores[0], k)
				s.mu.Unlock()
				wg.Done()
			}(k, v)
			continue
		}
		s.mu.Lock()
		s.stores[s.idx-1][k] = v
		s.mu.Unlock()
		wg.Done()
	}
	wg.Wait()
	s.clearTx()
	return nil
}

// clearTx is a helper function to clear temp stores and free up the memory
func (s *Kevast) clearTx() {
	s.mu.Lock()
	s.stores = s.stores[:len(s.stores)-1]
	s.mu.Unlock()
	s.idx--
}
