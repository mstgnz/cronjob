package conn

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ES struct {
	*elasticsearch.Client
}

func NewESClient() (*ES, error) {
	connStr := fmt.Sprintf("%s:%s", os.Getenv("ES_HOST"), os.Getenv("ES_PORT"))
	cfg := elasticsearch.Config{
		Addresses: []string{connStr},
		Username:  os.Getenv("ES_USER"),
		Password:  os.Getenv("ES_PASS"),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Millisecond,
			DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				// ...
			},
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	es := &ES{client}

	return es, nil
}
