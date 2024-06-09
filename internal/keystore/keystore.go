package keystore

import (
	"errors"
	"sync"
	"time"
)


type Store struct{
	mu sync.RWMutex
	store map[string] storageKey
}

type storageKey struct{
	keyname string
	value interface{}
	Expiry time.Time
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

	currTimeStamp := time.Now()

	value, exists := s.store[keyName]

	if !exists {
		return value,  errors.New("key does not exit within this store")
	}

	if currTimeStamp.After(value.Expiry){
		return value,errors.New("key has expired")
	}

	return value, nil
}

func(s *Store) SetKey(keyName string,keyValue interface{},expiryDuration int,clientMetaData string) error{

	newStorageKeyInstance := storageKey{
		keyname: keyName,
		value: keyValue,
		Expiry: time.Now().Add(time.Duration(expiryDuration)),
		metaData: clientMetaData,
		active: true,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[keyName] = newStorageKeyInstance

	return nil
}

func (s *storageKey) Value() string {

    if str, ok := s.value.(string); ok {
        return str
    }
    return "" // or handle the case where value is not a string
}


// func (k *Key) SetKey(key string,value string){
// 	k.Key
// }