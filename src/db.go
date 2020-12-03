package main

import (
	"fmt"
	"strconv"

	"github.com/ansel1/merry"
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func init() {
	_db, err := leveldb.OpenFile("path/db", nil)
	merry.Wrap(err)
	db = _db
}

func get(key string) (value string) {
	fmt.Println("get", key)
	data, err := db.Get([]byte(key), nil)
	merry.Wrap(err)
	return string(data)
}

func set(key string, value int) {
	fmt.Println("set", key, value)
	err := db.Put([]byte(key), []byte(strconv.Itoa(value)), nil)
	merry.Wrap(err)
}

func getAll() {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key), string(value))
	}
	iter.Release()
	iter.Error()
}
