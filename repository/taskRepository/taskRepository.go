package taskRepository

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/repository/repoUtil"
)

func Init() {
	db := repoUtil.OpenDb()
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
	db := repoUtil.OpenDb()
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
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		err := b.Delete(repoUtil.Itob(taskId))
		return err
	})
	defer db.Close()
}

func AddTask(task string, owner string) {
	db := repoUtil.OpenDb()
	newTask := model.MyTask{Task: task, Owner: owner}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		id, _ := b.NextSequence()
		newTask.ID = int(id)
		buf, err := json.Marshal(newTask)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(newTask.ID), buf)
	})

	defer db.Close()
}
