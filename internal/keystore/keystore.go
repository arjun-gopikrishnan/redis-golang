package keystore

import "errors"


type Store struct{
	store map[string] storageKey
}

type storageKey struct{
	keyname string
	value string
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

func(s *Store) SetKey(keyName string,keyValue string,expiryDuration int,clientMetaData string) error{

	newStorageKeyInstance := storageKey{
		keyname: keyName,
		value: keyValue,
		expiry: expiryDuration,
		metaData: clientMetaData,
		active: true,
	}

	s.store[keyName] = newStorageKeyInstance

	return nil
}

func (s *Store) GetKey(keyName string) (storageKey,error){
	value, exists := s.store[keyName]

	if !exists {
		return value,  errors.New("key does not exist within this store")
	}


	return value, nil
}
// func (k *Key) SetKey(key string,value string){
// 	k.Key
// }