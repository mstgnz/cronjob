package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/handler"
	"github.com/mstgnz/cronjob/schedule"
	"github.com/robfig/cron/v3"
)

var PORT string

func init() {
	// Load Env
	if err := godotenv.Load(".env"); err != nil {
		config.App().ErrorLog.Fatalf("Load Env Error: %v", err)
	}
	// init conf
	_ = config.App()

	// Load Sql
	config.App().QUERY = make(map[string]string)
	if query, err := config.LoadSQLQueries(); err != nil {
		config.App().ErrorLog.Fatalf("Load Sql Error: %v", err)
	} else {
		config.App().QUERY = query
	}

	PORT = os.Getenv("PORT")
}

func main() {

	//err := config.App().Mail.SetSubject("tars cron").SetContent("mail geldi mi?").SetTo("mesutgenez@hotmail.com").SetAttachment(map[string][]byte{"query.sql": []byte("1. dosya içeriği"), "query2.sql": []byte("2. dosya içeriği")}).SendText()

	// Scheduler Call
	c := cron.New()
	schedule.CallSchedule(c)
	// Start the Cron job scheduler
	c.Start()

	// Chi Define Routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fileServer(r, "/assets", http.Dir("./assets"))

	r.Get("/", handler.HomeHandler)
	r.Get("/post", handler.PostHandler)
	r.Get("/create", handler.CreateHandler)
	r.Post("/create", handler.CreateHandler)

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Run your HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
		if err != nil && err != http.ErrServerClosed {
			config.App().ErrorLog.Fatal(err.Error())
		}
	}()

	config.App().InfoLog.Printf("Tars Cron is running on %s", PORT)

	// Block until a signal is received
	<-ctx.Done()

	config.App().InfoLog.Printf("Tars Cron is shutting on %s", PORT)

	// set Shutting
	config.App().Shutting = true

	// check Running
	for {
		if config.App().Running <= 0 {
			config.App().InfoLog.Println("BREAK", config.App().Running)
			break
		} else {
			config.App().InfoLog.Printf("Currently %d active jobs in progress. pending completion...", config.App().Running)
		}
		time.Sleep(time.Second * 5)
	}

	config.App().InfoLog.Println("Shutting down gracefully...")

	config.App().DB.CloseDatabase()
	c.Stop()
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	// chi.StripPrefix ile servis edilecek route'u belirleyin
	r.Get(path+"*", http.StripPrefix(path, http.FileServer(root)).ServeHTTP)
}
