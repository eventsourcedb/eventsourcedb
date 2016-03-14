package eventsourcedb

import "github.com/boltdb/bolt"

type Env struct {
	DB *bolt.DB
}
