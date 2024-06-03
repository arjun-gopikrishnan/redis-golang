package tests

import (
	"testing"

	"github.com/arjun/redis-go/internal/keystore"
)

func TestNewStore(t *testing.T) {
	store,err := keystore.NewStore()
	if err!=nil{
		t.Fatalf("unable to set key")
	}
	if store == nil {
		t.Fatalf("expected new store to be initialized, got nil")
	}
	// if len(store.store) != 0 {
	// 	t.Fatalf("expected new store to be empty, got %d elements", len(store.store))
	// }
}

func TestSetKey(t *testing.T) {
	store,_ := keystore.NewStore()
	err := store.SetKey("testKey", "testValue", 10, "test metadata")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// if len(store.store) != 1 {
	// 	t.Fatalf("expected store to have 1 element, got %d", len(store.store))
	// }
}

// func TestGetKey(t *testing.T) {
// 	store,err := keystore.NewStore()
// 	store.SetKey("testKey", "testValue", 10, "test metadata")

// 	value, err := store.GetKey("testKey")
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// 	if value != "testValue" {
// 		t.Fatalf("expected value 'testValue', got %v", value)
// 	}

// 	_, err = store.GetKey("nonExistentKey")
// 	if err == nil {
// 		t.Fatalf("expected error for non-existent key, got none")
// 	}

// 	time.Sleep(11 * time.Second) // Wait for the key to expire
// 	_, err = store.GetKey("testKey")
// 	if err == nil {
// 		t.Fatalf("expected error for expired key, got none")
// 	}
// 	if err.Error() != "key has expired" {
// 		t.Fatalf("expected 'key has expired' error, got %v", err)
// 	}
// }

// func TestKeyExpiry(t *testing.T) {
// 	store,err := keystore.NewStore()
// 	store.SetKey("testKey", "testValue", 1*time.Second, "test metadata")

// 	time.Sleep(2 * time.Second) // Wait for the key to expire
// 	_, err := store.GetKey("testKey")
// 	if err == nil {
// 		t.Fatalf("expected error for expired key, got none")
// 	}
// 	if err.Error() != "key has expired" {
// 		t.Fatalf("expected 'key has expired' error, got %v", err)
// 	}
// }
