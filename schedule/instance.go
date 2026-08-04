package schedule

import (
	"fmt"
	"os"
	"sync"

	"github.com/mstgnz/cronjob/pkg/auth"
)

var (
	instanceOnce sync.Once
	instanceName string
)

// InstanceID names this process in the triggered table. It has to be unique per
// process: two replicas sharing an id could release each other's locks.
// The hostname is included because it is what identifies a pod in Kubernetes and
// makes a stuck lock traceable to a specific instance; the random suffix keeps the
// id unique when a pod restarts and reuses its name.
func InstanceID() string {
	instanceOnce.Do(func() {
		host, err := os.Hostname()
		if err != nil || host == "" {
			host = "unknown"
		}
		instanceName = fmt.Sprintf("%s-%d-%s", host, os.Getpid(), auth.RandomHex(4))
	})
	return instanceName
}
