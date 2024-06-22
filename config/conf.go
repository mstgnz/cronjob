package config

import (
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mstgnz/cronjob/pkg"
	"github.com/robfig/cron/v3"
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
	ES        *ES
	Kraft     *Kraft
	Mail      *pkg.Mail
	Cron      *cron.Cron
	Cache     *pkg.Cache
	Log       *Logger
	Validador *validator.Validate
	SecretKey string
	QUERY     map[string]string
	Running   int
	Shutting  bool
}

func App() *config {
	once.Do(func() {
		instance = &config{
			DB:        &DB{},
			Cron:      cron.New(),
			Cache:     pkg.NewCache(),
			Log:       &Logger{},
			Validador: validator.New(),
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
		// Connect to Postgres DB
		instance.DB.ConnectDatabase()
		// Connect to Kafka Kraft
		if kraft, err := newKafkaClient(); err != nil {
			log.Println(err)
		} else {
			instance.Kraft = kraft
		}
		// Connect to Elastic Search
		if es, err := newESClient(); err != nil {
			log.Println(err)
		} else {
			instance.ES = es
		}
	})
	return instance
}

func StructToMap(obj interface{}) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result
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

func GetIntQuery(r *http.Request, name string) int {
	pageStr := r.URL.Query().Get(name)
	if page, err := strconv.Atoi(pageStr); err == nil {
		return int(math.Abs(float64(page)))
	}
	return 1
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
