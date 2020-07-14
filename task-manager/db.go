package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type DB struct {
	Conn       *bolt.DB
	BucketName []byte
}

func (db *DB) Open(fileName string, bucketName string) error {
	// Open connection to database
	conn, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return err
	}

	// Create bucket
	db.BucketName = []byte(bucketName)
	conn.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(db.BucketName)
		if err != nil {
			return err
		}
		fmt.Printf("Opening bucket %v...\n", bucket)
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

func (db *DB) List() []string {
	var list []string
	db.Conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.BucketName)

		bucket.ForEach(func(k, v []byte) error {
			list = append(list, string(v))
			return nil
		})

		return nil
	})

	return list
}
