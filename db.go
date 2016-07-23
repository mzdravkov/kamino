package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// Takes a Bolt database, tenant's name as a string and port for the tenant as an int
// and puts the key-value pair (tenant, port) in the DB
func AddTenant(db *bolt.DB, tenant string, port int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tenants"))
		err := b.Put([]byte(tenant), []byte(strconv.Itoa(port)))
		return err
	})
}

// Gets the port for a tenant from the provided DB
func GetTenantPort(db *bolt.DB, tenant string) (int, error) {
	var port []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tenants"))
		port = b.Get([]byte(tenant))
		return nil
	})
	if err != nil {
		return 0, err
	}

	if port == nil {
		return 0, errors.New(fmt.Sprintf("Can't get the port for tenant '%s' from database", tenant))
	}
	return strconv.Atoi(string(port))
}

func InitDB(file string) *bolt.DB {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tenants"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
