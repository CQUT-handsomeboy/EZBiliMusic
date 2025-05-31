package database

import (
	"fmt"
	"log"
	"testing"

	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

func TestDatabase(t *testing.T) {
	exepath, _ := os.Executable()
	dbpath := filepath.Join(filepath.Dir(filepath.Dir(exepath)), "music.db")
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte("music"))

		if err != nil {
			return err
		}

		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("music"))

		if err != nil {
			return err
		}

		answer := b.Get([]byte("answer"))
		fmt.Println(string(answer))

		answer = b.Get([]byte("anthor"))
		fmt.Println(string(answer) == "")

		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

}

func TestDatabaseFunc(t *testing.T) {
	if _, err := os.Stat(defaultDBPath); err == nil {
		if err := os.Remove(defaultDBPath); err != nil {
			return
		}
	}

	InitDatabase()

	isExist, err := CheckMusicRecord("薛之谦", "喜欢你")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("isExist: %v\n", isExist) // false

	err = InsertMuiscRecord("薛之谦", "喜欢你", "filepath")

	if err != nil {
		fmt.Println(err)
		return
	}

	isExist, err = CheckMusicRecord("薛之谦", "喜欢你")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("isExist: %v\n", isExist) // true

	filename, err := DeleteMusicRecord("薛之谦", "喜欢你")

	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("filename:",filename)

	isExist, err = CheckMusicRecord("薛之谦", "喜欢你")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("isExist: %v\n", isExist) // false

}
