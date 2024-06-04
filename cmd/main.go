package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mstgnz/cronjob/config"
	"github.com/mstgnz/cronjob/handler"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/schedule"
	"github.com/robfig/cron/v3"
)

var PORT string

func init() {
	// Load Env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Load Env Error: %v", err)
	}
	// init conf
	_ = config.App()

	// Load Sql
	config.App().QUERY = make(map[string]string)
	if query, err := config.LoadSQLQueries(); err != nil {
		log.Fatalf("Load Sql Error: %v", err)
	} else {
		config.App().QUERY = query
	}

	PORT = os.Getenv("PORT")
}

var (
	webHandler handler.Web
	apiHandler handler.Api
)

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func Catch(h HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Info("HTTP handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func main() {

	// test mail with attach
	//err := config.App().Mail.SetSubject("cron").SetContent("mail geldi mi?").SetTo("mesutgenez@hotmail.com").SetAttachment(map[string][]byte{"query.sql": []byte("1. dosya içeriği"), "query2.sql": []byte("2. dosya içeriği")}).SendText()

	// Scheduler Call
	c := cron.New()
	schedule.CallSchedule(c)
	// Start the Cron job scheduler
	//c.Start()

	// Chi Define Routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	workDir, _ := os.Getwd()
	fileServer(r, "/assets", http.Dir(filepath.Join(workDir, "assets")))

	// swagger
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./template/swagger.html")
	})

	// web without auth
	r.Get("/login", Catch(webHandler.LoginHandler))
	r.Post("/login", Catch(webHandler.LoginHandler))
	r.Get("/register", Catch(webHandler.RegisterHandler))
	r.Post("/register", Catch(webHandler.RegisterHandler))

	// web with auth
	r.Group(func(r chi.Router) {
		r.Use(webAuthMiddleware)
		r.Get("/", Catch(webHandler.HomeHandler))
		r.Get("/profile", Catch(webHandler.ProfileHandler))
		r.Get("/schedules", Catch(webHandler.ListHandler))
	})

	// api without auth
	r.With(headerMiddleware).Post("/api/login", Catch(apiHandler.LoginHandler))
	r.With(headerMiddleware).Post("/api/register", Catch(apiHandler.RegisterHandler))

	// api with auth
	r.Group(func(r chi.Router) {
		r.Use(apiAuthMiddleware)
		r.Route("/api", func(r chi.Router) {
			r.Use(headerMiddleware)
			r.Get("/user", Catch(apiHandler.UserHandler))
			r.Put("/user", Catch(apiHandler.UserUpdateHandler))
			r.Get("/schedules", Catch(apiHandler.ScheduleListHandler))
			r.Post("/schedule", Catch(apiHandler.ScheduleCreateHandler))
			r.Put("/schedule", Catch(apiHandler.ScheduleUpdateHandler))
			r.Delete("/schedule/{id}", Catch(apiHandler.ScheduleDeleteHandler))
			r.Get("/schedule/mail/{schedule_id}", Catch(apiHandler.ScheduleMailListHandler))
			r.Post("/schedule/mail", Catch(apiHandler.ScheduleMailCreateHandler))
			r.Put("/schedule/mail", Catch(apiHandler.ScheduleMailUpdateHandler))
			r.Delete("/schedule/mail/{id}", Catch(apiHandler.ScheduleMailDeleteHandler))
		})
	})

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Run your HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	log.Printf("Cron is running on %s", PORT)

	// Block until a signal is received
	<-ctx.Done()

	log.Printf("Cron is shutting on %s", PORT)

	// set Shutting
	config.App().Shutting = true

	// check Running
	for {
		if config.App().Running <= 0 {
			log.Println("BREAK", config.App().Running)
			break
		} else {
			log.Printf("Currently %d active jobs in progress. pending completion...", config.App().Running)
		}
		time.Sleep(time.Second * 5)
	}

	log.Println("Shutting down gracefully...")

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

func webAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		userId, err := config.GetUserIDByToken(tokenString)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user_id, err := strconv.Atoi(userId)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user := &models.User{}
		getUser, err := user.GetUserWithId(user_id)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), getUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func apiAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			_ = config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: "Invalid Token"})
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		userId, err := config.GetUserIDByToken(tokenString)
		if err != nil {
			_ = config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
			return
		}

		user_id, err := strconv.Atoi(userId)
		if err != nil {
			_ = config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
			return
		}

		user := &models.User{}
		getUser, err := user.GetUserWithId(user_id)

		if err != nil {
			_ = config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), getUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			_ = config.WriteJSON(w, http.StatusBadRequest, config.Response{Status: false, Message: "Invalid Content-Type"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
