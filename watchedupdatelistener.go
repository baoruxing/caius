package caius

type WatchedUpdateResult struct {
	Complete map[string]interface{}
	Added    map[string]interface{}
	Changed  map[string]interface{}
	Deleted  map[string]interface{}
}

type WatchedUpdateListener interface {
	UpdateConfiguration(result WatchedUpdateResult)
}
