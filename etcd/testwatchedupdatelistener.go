package etcd

import (
	"github.com/baoruxing/caius"
	"log"
)

type TestWatchedUpdatedListener struct {
}

func (t TestWatchedUpdatedListener) UpdateConfiguration(result caius.WatchedUpdateResult) {
	log.Println("added ", result.Added)
	log.Println("changed ", result.Changed)
	log.Println("deleted ", result.Deleted)
}
