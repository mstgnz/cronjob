package config

import (
	"os"
	"sync"

	"github.com/mstgnz/cronjob/pkg"
)

var (
	once     sync.Once
	mu       sync.Mutex
	instance *config
)

// context key type
type CKey string

type config struct {
	DB        *DB
	Mail      *pkg.Mail
	Cache     *pkg.Cache
	Log       *Logger
	SecretKey string
	QUERY     map[string]string
	Running   int
	Shutting  bool
}

func App() *config {
	once.Do(func() {
		instance = &config{
			DB:    &DB{},
			Cache: pkg.NewCache(),
			Log:   &Logger{},
			// the secret key will change every time the application is restarted.
			SecretKey: "asdf1234", //RandomString(8),
			Mail: &pkg.Mail{
				From: os.Getenv("MAIL_FROM"),
				Name: os.Getenv("MAIL_FROM_NAME"),
				Host: os.Getenv("MAIL_HOST"),
				Port: os.Getenv("MAIL_PORT"),
				User: os.Getenv("MAIL_USER"),
				Pass: os.Getenv("MAIL_PASS"),
			},
		}
		instance.DB.ConnectDatabase()
	})
	return instance
}

func ShuttingWrapper(fn func()) {
	if !App().Shutting {
		fn()
	}
}

func IncrementRunning() {
	mu.Lock()
	App().Running++
	mu.Unlock()
}

func DecrementRunning() {
	mu.Lock()
	App().Running--
	mu.Unlock()
}
