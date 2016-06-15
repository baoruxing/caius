package etcd

import (
	"fmt"
	"testing"
	"time"
)

func TestNewEtcdWatchedConfigSource(t *testing.T) {
	endpoints := []string{"http://127.0.0.1:4001"}
	configPath := "/member/register"
	source := NewEtcdWatchedConfigSource(endpoints, configPath)
	for {
		data, _ := source.GetCurrentConfigData()
		fmt.Println(data)
		time.Sleep(1000 * time.Millisecond)
	}

	/*for {
		for _, value := range data {
			if value != "bar" {
				t.Errorf("Want=bar, got=%s ", value)
			} else {
				t.Logf("got=%s", value)
			}
		}
	}*/

}
