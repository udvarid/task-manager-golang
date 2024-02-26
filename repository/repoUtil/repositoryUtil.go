package repoUtil

import (
	"encoding/binary"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func OpenDb() *bolt.DB {
	db, err := bolt.Open("./db/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
