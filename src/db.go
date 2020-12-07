package main

import (
	"fmt"
	"sort"
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

func get(key string) (value int) {
	data, err := db.Get([]byte(key), nil)
	merry.Wrap(err)
	count, err := strconv.Atoi(string(data))
	merry.Wrap(err)
	return count

}

func set(key string, value int, domain string) {
	fmt.Println("set", key, value, domain)
	err := db.Put([]byte(key), []byte(strconv.Itoa(value)), nil)
	merry.Wrap(err)
}

// 返回已经按照出现频率排过序的标签数据
func getAll() []Tag {
	var tags []Tag
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()
		count, err := strconv.Atoi(string(value))
		merry.Wrap(err)
		tags = append(tags, Tag{
			name:  key,
			count: count,
		})
		// 记录完后删除之
		err = db.Delete(iter.Key(), nil)
		merry.Wrap(err)
	}
	iter.Release()
	iter.Error()

	sort.Slice(tags, func(i, j int) bool { return tags[i].count > tags[j].count })
	return tags
}
