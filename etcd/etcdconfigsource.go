package etcd

import (
	"container/list"
	"github.com/baoruxing/caius/"
	"strings"
)

type EtcdWatchedConfigSource struct {
	valueCache map[string]interface{}
	connPath   string
	listeners  List
}

func (e *EtcdWatchedConfigSource) AddUpdateListener(listener WatchedUpdateListener) {
	if listener != nil {
		e.listeners.PushBack(listener)
	}

}

func (e *EtcdWatchedConfigSource) RemoveUpdateListener(listener WatchedUpdateListener) {
	if listener != nil {
		e.listeners.Remove(listener)
	}

}

func (e *EtcdWatchedConfigSource) GetCurrentConfigData() (*map[string]interface{}, error) {
	return e.valueCache
}
