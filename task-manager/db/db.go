package db

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type DB struct {
	Conn       *bolt.DB
	BucketName []byte
}

type Task struct {
	Key  []byte
	Task []byte
}

func (db *DB) Open(fileName string, bucketName string) error {
	// Return if connection is already established
	if db.Conn != nil {
		return nil
	}

	// Open connection to database
	conn, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return err
	}

	// Create bucket
	db.BucketName = []byte(bucketName)
	conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(db.BucketName)
		if err != nil {
			return err
		}
		return nil
	})

	db.Conn = conn
	return nil
}

func (db *DB) Insert(value string) {
	db.Conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.BucketName)

		// Generate ID based on current time
		t := time.Now()
		id, err := t.MarshalText()
		if err != nil {
			log.Fatal(err)
		}

		// Store new record
		err = bucket.Put(id, []byte(value))
		return err
	})
}

func (db *DB) List() []Task {
	var list []Task

	db.Conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.BucketName)
		bucket.ForEach(func(k, v []byte) error {
			list = append(list, Task{
				Key:  k,
				Task: v,
			})
			return nil
		})
		return nil
	})

	return list
}

func (db *DB) Delete(key []byte) {
	db.Conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.BucketName)
		bucket.Delete(key)

		return nil
	})
}
