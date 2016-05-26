package etcd

import (
	"container/list"
	"github.com/baoruxing/caius/"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"strings"
)

type EtcdWatchedConfigSource struct {
	configPath string
	valueCache map[string]interface{}
	listeners  List
	KeysAPI    client.KeysAPI
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

func NewEtcdWatchedConfigSource(endpoints []string, configPath string) *EtcdWatchedConfigSource {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	configSource := &EtcdWatchedConfigSource{
		configPath: configPath,
		valueCache: make(map[string]interface{}),
		listeners:  list.New(),
		KeysAPI:    client.NewKeysAPI(etcdClient),
	}
	go configSource.updateHandler()
	return configSource
}

func (e *EtcdWatchedConfigSource) updateHandler() {
	api := e.KeysAPI
	watcher := api.Watcher("configPath", &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background)
		if err != nil {
			log.Println("Error watch configPath:", err)
			break
		}
		etcdKey := res.Node.Key
		splitedKey := trings.Split(etcdKey, "/")
		sourceKey := splitedKey[len(splitedKey)-1]
		value := res.Node.Value
		if res.Action == "expire" {

		} else if res.Action == "create" || res.Action == "set" || res.Action == "update" {
			e.valueCache[sourceKey] = value
		} else if res.Action == "delete" {
			delete(e.valueCache, value)
		}
	}
}
