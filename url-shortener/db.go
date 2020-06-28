package urlshort

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/boltdb/bolt"
)

var bucketName []byte = []byte("PathURLS")
var db *bolt.DB

func AddPathURL(path string, URL string) error {
	db = connect()
	defer db.Close()

	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return fmt.Errorf("Invalid URL...\nExample URL: https://google.com\n")
	}

	insertKeyValue(path, URL)
	return nil
}

func GetPathURL(path string) (string, error) {
	db = connect()
	defer db.Close()

	value := getValue(path)
	if len(value) == 0 {
		return "", fmt.Errorf("No URL found from given path %v", path)
	}

	return value, nil
}

func connect() *bolt.DB {
	db, err := bolt.Open("paths.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	return db
}

func insertKeyValue(key string, value string) {
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
}

func getValue(key string) string {
	var value []byte

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		value = bucket.Get([]byte(key))

		return nil
	})

	return string(value)
}
