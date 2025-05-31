package database

import (
	"fmt"

	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var defaultDBPath string

func InitDatabase() error {
	db, err := bolt.Open(defaultDBPath, 0600, nil)

	if err != nil {
		fmt.Println("open database failed...")
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte("music"))

		if err != nil {
			fmt.Println("create bucket failed...")
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("update failed...")
		return err
	}

	return nil
}

func init() {
	exepath, _ := os.Executable()
	defaultDBPath = filepath.Join(filepath.Dir(exepath), "database.db")

	err := InitDatabase()
	if err != nil {
		fmt.Println("init database failed...")
		panic(err)
	}
}

func InsertMuiscRecord(title string, artist string, filename string) error {
	db, err := bolt.Open(defaultDBPath, 0600, nil)

	if err != nil {
		fmt.Println("open database failed...")
		return err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte("music"))

		if err != nil {
			return err
		}

		recordKey := fmt.Sprintf("%s:%s", artist, title)
		recordValue := filename

		err = b.Put([]byte(recordKey), []byte(recordValue))
		return err
	})

	if err != nil {
		fmt.Println("insert failed...")
		return err
	}

	return nil
}

func DeleteMusicRecord(title string, artist string) (string, error) {
	db, err := bolt.Open(defaultDBPath, 0600, nil)

	if err != nil {
		fmt.Println("open database failed...")
		return "", err
	}

	defer db.Close()

	var filename string

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte("music"))

		if err != nil {
			return err
		}

		recordKey := fmt.Sprintf("%s:%s", artist, title)
		filename = string(b.Get([]byte(recordKey)))
		err = b.Delete([]byte(recordKey))
		return err
	})

	if err != nil {
		fmt.Println("delete failed...")
		return "", err
	}

	return filename, nil
}

func CheckMusicRecord(title string, artist string) (bool, error) {
	db, err := bolt.Open(defaultDBPath, 0600, nil)

	if err != nil {
		fmt.Println("open database failed...")
		return false, err
	}

	defer db.Close()

	var result bool

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("music"))

		if err != nil {
			return err
		}
		recordKey := fmt.Sprintf("%s:%s", artist, title)
		recordValue := string(b.Get([]byte(recordKey)))

		if recordValue != "" {
			result = true
		} else {
			result = false
		}

		return err
	})

	if err != nil {
		fmt.Println("delete failed...")
		return false, err
	}

	return result, nil
}
