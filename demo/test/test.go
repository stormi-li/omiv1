package main

import (
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
)

func main() {
	opts := badger.DefaultOptions("./badger")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 写入带有过期时间的数据
	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("session"), []byte("session_data")).WithTTL(10 * time.Second)
		return txn.SetEntry(e)
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Data with TTL written successfully!")

	// 验证键过期
	time.Sleep(12 * time.Second) // 等待数据过期

	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("session"))
		if err == badger.ErrKeyNotFound {
			log.Println("Key has expired")
			return nil
		}
		return err
	})
	if err != nil && err != badger.ErrKeyNotFound {
		log.Fatal(err)
	}
}