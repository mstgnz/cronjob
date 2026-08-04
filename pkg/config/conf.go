package config

import (
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/go-playground/validator/v10"
	"github.com/mstgnz/cronjob/pkg/cache"
	"github.com/mstgnz/cronjob/pkg/conn"
	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/robfig/cron/v3"
)

type CKey string

type Manager struct {
	DB        *conn.DB
	Mail      *response.Mail
	ES        *conn.ES
	Cache     *cache.Cache
	Kafka     *conn.Kraft
	Redis     *conn.Redis
	Cron      *cron.Cron
	Validator *validator.Validate
	SecretKey string
	QUERY     map[string]string
	Running   int
	Shutting  bool
}

var (
	instance *Manager
	mu       sync.Mutex
)

// newCron builds the scheduler used by the whole application.
// Location comes from DB_ZONE; without it robfig/cron falls back to time.Local,
// which is UTC inside a container and silently shifts every schedule.
// The job chain is explicit because cron.New() ships an empty chain: without
// Recover a panic inside a job kills the process, and without SkipIfStillRunning
// a slow endpoint lets executions of the same schedule pile up.
func newCron() *cron.Cron {
	loc := time.Local
	if zone := os.Getenv("DB_ZONE"); zone != "" {
		if parsed, err := time.LoadLocation(zone); err != nil {
			log.Printf("cron: unknown time zone %q, falling back to %s: %v", zone, loc, err)
		} else {
			loc = parsed
		}
	}
	return cron.New(
		cron.WithLocation(loc),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)
}

func App() *Manager {
	if instance == nil {
		instance = &Manager{
			DB:        &conn.DB{},
			Cache:     &cache.Cache{},
			Kafka:     &conn.Kraft{},
			Redis:     &conn.Redis{},
			Cron:      newCron(),
			Validator: validator.New(),
			// the secret key will change every time the application is restarted.
			SecretKey: os.Getenv("JWT_SECRET"), //RandomString(8),
			Mail: &response.Mail{
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
		/* if kraft, err := conn.NewKafkaClient(); err != nil {
			log.Println(err)
		} else {
			instance.Kafka = kraft
		} */
		// Connect to Elastic Search
		/* if es, err := conn.NewESClient(); err != nil {
			log.Println(err)
		} else {
			instance.ES = es
		} */
	}
	return instance
}

func StructToMap(obj any) map[string]any {
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
	if !IsShutting() {
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

// RunningCount reports how many schedule executions are in flight. Reading the
// field directly would race with the counter updates above.
func RunningCount() int {
	mu.Lock()
	defer mu.Unlock()
	return App().Running
}

// SetShutting marks the application as shutting down so no further jobs start.
func SetShutting() {
	mu.Lock()
	App().Shutting = true
	mu.Unlock()
}

// IsShutting reports whether shutdown has started.
func IsShutting() bool {
	mu.Lock()
	defer mu.Unlock()
	return App().Shutting
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

func ActiveClass(a, b int) string {
	active := ""
	if a == b {
		active = "active"
	}
	return active
}

func WriteBody(r *http.Request) {
	if body, err := io.ReadAll(r.Body); err != nil {
		log.Println("WriteBody: ", err)
	} else {
		log.Println("WriteBody: ", string(body))
	}
}
