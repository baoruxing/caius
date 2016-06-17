package caius

type WatchedConfigSource interface {
	AddUpdateListener(listener WatchedUpdateListener)

	RemoveUpdateListener(listener WatchedUpdateListener)

	GetCurrentConfigData() (*map[string]interface{}, error)
}
