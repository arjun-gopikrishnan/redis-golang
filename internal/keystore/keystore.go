package keystore

import (
	"errors"
	"sync"
)


type Store struct{
	mu sync.RWMutex
	store map[string] storageKey
}

type storageKey struct{
	keyname string
	value interface{}
	expiry int
	metaData string
	active bool
}

func NewStore() (*Store,error){
	newStoreInstance := &Store{
		store: make(map[string]storageKey),
	}
	return newStoreInstance,nil
}

func (s *Store) GetKey(keyName string) (storageKey,error){
	value, exists := s.store[keyName]

	if !exists {
		return value,  errors.New("key does not exit within this store")
	}
	return value, nil
}

func(s *Store) SetKey(keyName string,keyValue interface{},expiryDuration int,clientMetaData string) error{

	newStorageKeyInstance := storageKey{
		keyname: keyName,
		value: keyValue,
		expiry: expiryDuration,
		metaData: clientMetaData,
		active: true,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[keyName] = newStorageKeyInstance

	return nil
}


// func (k *Key) SetKey(key string,value string){
// 	k.Key
// }