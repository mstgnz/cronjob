package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
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

var (
	authHandler     handler.Auth
	homeHandler     handler.Home
	scheduleHandler handler.Timing
)

func main() {

	//err := config.App().Mail.SetSubject("tars cron").SetContent("mail geldi mi?").SetTo("mesutgenez@hotmail.com").SetAttachment(map[string][]byte{"query.sql": []byte("1. dosya içeriği"), "query2.sql": []byte("2. dosya içeriği")}).SendText()

	// Scheduler Call
	c := cron.New()
	schedule.CallSchedule(c)
	// Start the Cron job scheduler
	//c.Start()

	// Chi Define Routes
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	workDir, _ := os.Getwd()
	fileServer(r, "/assets", http.Dir(filepath.Join(workDir, "assets")))

	r.Route("/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Get("/", authHandler.LoginHandler)
			r.Post("/", authHandler.LoginHandler)
		})
		r.Route("/register", func(r chi.Router) {
			r.Get("/", authHandler.RegisterHandler)
			r.Post("/", authHandler.RegisterHandler)
		})
		r.With(authMiddleware).Route("/profile", func(r chi.Router) {
			r.Get("/", authHandler.UpdateHandler)
			r.Post("/", authHandler.DeleteHandler)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", homeHandler.HomeHandler)
		r.Route("/schedule", func(r chi.Router) {
			r.Get("/", scheduleHandler.ScheduleHandler)
			r.Post("/", scheduleHandler.ScheduleHandler)
		})
	})

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

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func CustomMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Burada middleware işlemleri yapabilirsiniz.
		fmt.Println("Middleware çalıştı!")

		// Sonraki middleware'e veya ana işleyiciye (handler) talebi iletmek için:
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwtToken")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		userId, err := config.GetUserIDByToken(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// get user
		var val config.WithValueVal = map[string]any{}

		// Kullanıcı ID'sini talep içerisine ekleyerek devam et
		ctx := context.WithValue(r.Context(), config.WithValueKey(userId), val)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
