package caius

import (
	"errors"
	"strings"
)

type WatchedConfigSource interface {
	AddUpdateListener(listener WatchedUpdateListener)

	RemoveUpdateListener(listener WatchedUpdateListener)

	GetCurrentConfigData() (*map[string]interface{}, error)
}
