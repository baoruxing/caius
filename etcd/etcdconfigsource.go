package etcd

import (
	"container/list"
	"github.com/baoruxing/caius"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

type EtcdWatchedConfigSource struct {
	configPath string
	valueCache map[string]interface{}
	listeners  *list.List
	KeysAPI    client.KeysAPI
}

func (e *EtcdWatchedConfigSource) AddUpdateListener(listener caius.WatchedUpdateListener) {
	if listener != nil {
		e.listeners.PushBack(listener)
	}

}

func (e *EtcdWatchedConfigSource) RemoveUpdateListener(listener caius.WatchedUpdateListener) {
	if listener != nil {
		for el := e.listeners.Front(); el != nil; el = el.Next() {
			if el.Value == listener {
				e.listeners.Remove(el)
			}
		}

	}

}

func (e *EtcdWatchedConfigSource) updateConfiguration(result caius.WatchedUpdateResult) {

	for el := e.listeners.Front(); el != nil; el = el.Next() {

		el.Value.(caius.WatchedUpdateListener).UpdateConfiguration(result)
	}
}

func (e *EtcdWatchedConfigSource) GetCurrentConfigData() (map[string]interface{}, error) {
	return e.valueCache, nil
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
	configSource.cacheValues()
	go configSource.updateHandler()
	return configSource
}

func (e *EtcdWatchedConfigSource) cacheValues() {
	api := e.KeysAPI
	res, err := api.Get(context.Background(), e.configPath, &client.GetOptions{Recursive: true})
	if err != nil {
		log.Println("Error get configPath:", err)
		return
	}
	nodes := res.Node.Nodes
	for _, node := range nodes {
		etcdKey := node.Key
		splitedKey := strings.Split(etcdKey, "/")
		sourceKey := splitedKey[len(splitedKey)-1]
		value := node.Value
		e.valueCache[sourceKey] = value
	}

}

func (e *EtcdWatchedConfigSource) updateHandler() {
	api := e.KeysAPI
	watcher := api.Watcher(e.configPath, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch configPath:", err)
			break
		}
		etcdKey := res.Node.Key
		splitedKey := strings.Split(etcdKey, "/")
		sourceKey := splitedKey[len(splitedKey)-1]
		value := res.Node.Value
		log.Println("res.Action = ", res.Action)
		if res.Action == "expire" {

		} else if res.Action == "create" {
			e.valueCache[sourceKey] = value
			result := caius.WatchedUpdateResult{
				Added: map[string]interface{}{
					sourceKey: value,
				},
			}
			e.updateConfiguration(result)
		} else if res.Action == "set" || res.Action == "update" {
			e.valueCache[sourceKey] = value
			result := caius.WatchedUpdateResult{
				Changed: map[string]interface{}{
					sourceKey: value,
				},
			}
			e.updateConfiguration(result)
		} else if res.Action == "delete" {
			delete(e.valueCache, sourceKey)
			result := caius.WatchedUpdateResult{
				Deleted: map[string]interface{}{
					sourceKey: "",
				},
			}
			e.updateConfiguration(result)
		}
	}
}
