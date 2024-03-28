package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/mstgnz/cronjob/pkg"
)

var (
	once     sync.Once
	mu       sync.Mutex
	instance *config
)

type config struct {
	DB        *DB
	Json      *Json
	Response  *Response
	Mail      *pkg.Mail
	Cache     *pkg.Cache
	SecretKey string
	QUERY     map[string]string
	Running   int
	Shutting  bool
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
}

func App() *config {
	once.Do(func() {
		instance = &config{
			DB:        &DB{},
			Json:      &Json{},
			Response:  &Response{},
			Cache:     pkg.NewCache(),
			SecretKey: GetSecretKey(),
			InfoLog:   log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
			ErrorLog:  log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
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
	//go Logger("info", instance.InfoLog)
	//go Logger("error", instance.ErrorLog)
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

func Logger(fileName string, logger *log.Logger) {
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.Mkdir(logsDir, 0755)
	}

	lastCheckedDay := time.Now().Day()

	for {
		currentDay := time.Now().Day()
		currentTime := time.Now().Format("2006-01-02")
		logFileName := logsDir + "/" + fileName + "-" + currentTime + ".log"

		_, err := os.Stat(fileName)
		if currentDay != lastCheckedDay || os.IsNotExist(err) {

			file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("An error occurred while creating the log file:", err)
			}

			logger.SetOutput(file)

			lastCheckedDay = currentDay

			_ = file.Close()
		}
	}
}
