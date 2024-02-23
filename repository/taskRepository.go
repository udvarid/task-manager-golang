package repository

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/udvarid/task-manager-golang/model"
)

var taskList = []model.MyTask{}

func Init() {
	db := openDb()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Sessions"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	defer db.Close()
}

func GetAllTask(owner string) []model.MyTask {
	db := openDb()
	var result []model.MyTask
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))

		b.ForEach(func(k, v []byte) error {
			var task model.MyTask
			json.Unmarshal([]byte(v), &task)
			if task.Owner == owner {
				result = append(result, task)
			}
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func DeleteTask(taskId int) {
	db := openDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		err := b.Delete(itob(taskId))
		return err
	})
	defer db.Close()
}

func AddTask(task string, owner string) {
	db := openDb()
	newTask := model.MyTask{Task: task, Owner: owner}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		id, _ := b.NextSequence()
		newTask.ID = int(id)
		buf, err := json.Marshal(newTask)
		if err != nil {
			return err
		}
		return b.Put(itob(newTask.ID), buf)
	})

	defer db.Close()
	taskList = append(taskList, newTask)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func openDb() *bolt.DB {
	db, err := bolt.Open("./db/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
