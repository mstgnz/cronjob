package main

import (
	"context"
	"fmt"
	"log"
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
	"github.com/mstgnz/cronjob/handler/api"
	"github.com/mstgnz/cronjob/handler/web"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/schedule"
	"github.com/robfig/cron/v3"
)

var (
	PORT       string
	webHandler web.Web

	apiUserHandler         api.UserHandler
	apiGroupHandler        api.GroupHandler
	apiRequestHandler      api.RequestHandler
	apiNotificationHandler api.NotificationHandler
	apiScheduleHandler     api.ScheduleHandler
	apiWebhookHandler      api.WebhookHandler
)

func init() {
	// Load Env
	if err := godotenv.Load(".env"); err != nil {
		config.App().Log.Warn(fmt.Sprintf("Load Env Error: %v", err))
		log.Fatalf("Load Env Error: %v", err)
	}
	// init conf
	_ = config.App()
	config.CustomValidate()

	// Load Sql
	config.App().QUERY = make(map[string]string)
	if query, err := config.LoadSQLQueries(); err != nil {
		config.App().Log.Warn(fmt.Sprintf("Load Sql Error: %v", err))
		log.Fatalf("Load Sql Error: %v", err)
	} else {
		config.App().QUERY = query
	}

	PORT = os.Getenv("PORT")
}

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func Catch(h HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			config.App().Log.Info("HTTP Handler Error", "err", err.Error(), "path", r.URL.Path)
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
	r.With(headerMiddleware).Post("/api/login", Catch(apiUserHandler.LoginHandler))
	r.With(headerMiddleware).Post("/api/register", Catch(apiUserHandler.RegisterHandler))

	// api with auth
	r.Group(func(r chi.Router) {
		r.Use(apiAuthMiddleware)
		r.Route("/api", func(r chi.Router) {
			r.Use(headerMiddleware)
			// users
			r.Get("/user", Catch(apiUserHandler.UserHandler))
			r.Put("/user", Catch(apiUserHandler.UserUpdateHandler))
			r.Delete("/user", Catch(apiUserHandler.UserDeleteHandler))
			r.Put("/user-change-pass", Catch(apiUserHandler.UserPassUpdateHandler))
			// groups
			r.Get("/groups", Catch(apiGroupHandler.GroupListHandler))
			r.Post("/groups", Catch(apiGroupHandler.GroupCreateHandler))
			r.Put("/groups/{id}", Catch(apiGroupHandler.GroupUpdateHandler))
			r.Delete("/groups/{id}", Catch(apiGroupHandler.GroupDeleteHandler))
			// requests
			r.Get("/requests", Catch(apiRequestHandler.RequestListHandler))
			r.Post("/requests", Catch(apiRequestHandler.RequestCreateHandler))
			r.Put("/requests/{id}", Catch(apiRequestHandler.RequestUpdateHandler))
			r.Delete("/requests/{id}", Catch(apiRequestHandler.RequestDeleteHandler))
			// request headers
			r.Get("/request-headers", Catch(apiRequestHandler.RequestHeaderListHandler))
			r.Post("/request-headers", Catch(apiRequestHandler.RequestHeaderCreateHandler))
			r.Put("/request-headers/{id}", Catch(apiRequestHandler.RequestHeaderUpdateHandler))
			r.Delete("/request-headers/{id}", Catch(apiRequestHandler.RequestHeaderDeleteHandler))
			// notifications
			r.Get("/notifications", Catch(apiNotificationHandler.NotificationListHandler))
			r.Post("/notifications", Catch(apiNotificationHandler.NotificationCreateHandler))
			r.Put("/notifications/{id}", Catch(apiNotificationHandler.NotificationUpdateHandler))
			r.Delete("/notifications/{id}", Catch(apiNotificationHandler.NotificationDeleteHandler))
			// notification emails
			r.Get("/notify-emails", Catch(apiNotificationHandler.NotifyEmailListHandler))
			r.Post("/notify-emails", Catch(apiNotificationHandler.NotifyEmailCreateHandler))
			r.Put("/notify-emails/{id}", Catch(apiNotificationHandler.NotifyEmailUpdateHandler))
			r.Delete("/notify-emails/{id}", Catch(apiNotificationHandler.NotifyEmailDeleteHandler))
			// notification sms
			r.Get("/notify-sms", Catch(apiNotificationHandler.NotifySmsListHandler))
			r.Post("/notify-sms", Catch(apiNotificationHandler.NotifySmsCreateHandler))
			r.Put("/notify-sms/{id}", Catch(apiNotificationHandler.NotifySmsUpdateHandler))
			r.Delete("/notify-sms/{id}", Catch(apiNotificationHandler.NotifySmsDeleteHandler))
			// webhooks
			r.Get("/webhooks", Catch(apiWebhookHandler.WebhookListHandler))
			r.Post("/webhooks", Catch(apiWebhookHandler.WebhookCreateHandler))
			r.Put("/webhooks/{id}", Catch(apiWebhookHandler.WebhookUpdateHandler))
			r.Delete("/webhooks/{id}", Catch(apiWebhookHandler.WebhookDeleteHandler))
			// schedules
			r.Get("/schedules", Catch(apiScheduleHandler.ScheduleListHandler))
			r.Post("/schedules", Catch(apiScheduleHandler.ScheduleCreateHandler))
			r.Put("/schedules/{id}", Catch(apiScheduleHandler.ScheduleUpdateHandler))
			r.Delete("/schedules/{id}", Catch(apiScheduleHandler.ScheduleDeleteHandler))
			// schedule logs
			r.Get("/schedule-logs", Catch(apiScheduleHandler.ScheduleLogListHandler))
			r.Post("/schedule-logs", Catch(apiScheduleHandler.ScheduleLogCreateHandler))
			r.Put("/schedule-logs/{id}", Catch(apiScheduleHandler.ScheduleLogUpdateHandler))
			r.Delete("/schedule-logs/{id}", Catch(apiScheduleHandler.ScheduleLogDeleteHandler))
		})
	})

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Run your HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
		if err != nil && err != http.ErrServerClosed {
			config.App().Log.Warn("Fatal Error", "err", err.Error())
			log.Fatal(err.Error())
		}
	}()

	config.App().Log.Info("Cron is running on", PORT)

	// Block until a signal is received
	<-ctx.Done()

	config.App().Log.Info("Cron is shutting on", PORT)

	// set Shutting
	config.App().Shutting = true

	// check Running
	for {
		if config.App().Running <= 0 {
			config.App().Log.Info("Cronjobs all done")
			break
		} else {
			config.App().Log.Info(fmt.Sprintf("Currently %d active jobs in progress. pending completion...", config.App().Running))
		}
		time.Sleep(time.Second * 5)
	}

	config.App().Log.Info("Shutting down gracefully...")

	config.App().DB.CloseDatabase()
	c.Stop()
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
		if err != nil && user_id == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user := &models.User{}
		err = user.GetWithId(user_id)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
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
		if err != nil && user_id == 0 {
			_ = config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
			return
		}

		user := &models.User{}
		err = user.GetWithId(user_id)

		if err != nil {
			_ = config.WriteJSON(w, http.StatusUnauthorized, config.Response{Status: false, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
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
