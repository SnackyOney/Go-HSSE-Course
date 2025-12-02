package stg_ex

import (
	"bank-service/exceptions"
	"bank-service/stg"
	"sync"
)

type UserInfo struct {
	id     int
	budget int
}

type storage struct {
	database map[int]*UserInfo
	mux      sync.RWMutex
}

func Newstorage() stg.StorageInterface {
	return &storage{
		database: make(map[int]*UserInfo),
	}
}

func (storage *storage) Get(userid int) (int, error) {
	storage.mux.RLock()
	defer storage.mux.RUnlock()
	if storage.database[userid] == nil {
		storage.database[userid] = &UserInfo{userid, 0}
	}
	return storage.database[userid].budget, nil
}

func (storage *storage) Deposit(userid int, amount int) (int, error) {
	storage.mux.Lock()
	defer storage.mux.Unlock()
	if storage.database[userid] == nil {
		storage.database[userid] = &UserInfo{userid, 0}
	}
	storage.database[userid].budget += amount
	return storage.database[userid].budget, nil
}

func (storage *storage) Transfer(from, to int, amount int) (int, int, error) {
	storage.mux.Lock()
	defer storage.mux.Unlock()

	if storage.database[from] == nil {
		storage.database[from] = &UserInfo{from, 0}
	}

	if storage.database[to] == nil {
		storage.database[to] = &UserInfo{to, 0}
	}

	if storage.database[from].budget < amount {
		return 0, 0, exceptions.ErrInvalidAmount
	}
	if storage.database[to].id == from {
		return 0, 0, exceptions.ErrSelfTransfer
	}

	storage.database[from].budget -= amount
	storage.database[to].budget += amount

	return storage.database[from].budget, storage.database[to].budget, nil
}
