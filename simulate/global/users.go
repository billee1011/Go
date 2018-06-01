package global

import (
	"fmt"
	"log"
	"steve/simulate/flag"
	"sync"

	bolt "github.com/coreos/bbolt"
)

var useridDB *bolt.DB
var once sync.Once

// AllocUserName 分配用户名
func AllocUserName() string {
	once.Do(setupDB)
	var userID uint64
	if err := useridDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("max_user_id"))
		if err != nil {
			return err
		}
		userID, err = b.NextSequence()
		return err
	}); err != nil {
		log.Panicln("分配用户ID失败")
	}
	return fmt.Sprintf("user_%v", userID)
}

func setupDB() {
	var err error
	file := fmt.Sprintf("%s/userid.db", flag.Flags.DBPath)
	useridDB, err = bolt.Open(file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}
