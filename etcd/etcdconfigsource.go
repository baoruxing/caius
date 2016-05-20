package etcd

type EtcdWatchedConfigSource struct {

}

func (e *EtcdWatchedConfigSource) AddUpdateListener() {
    ValueCache  map[string] interface{}
}