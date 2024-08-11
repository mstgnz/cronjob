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
	"github.com/mstgnz/cronjob/handler/api"
	"github.com/mstgnz/cronjob/handler/web"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/auth"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/load"
	"github.com/mstgnz/cronjob/pkg/logger"
	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/pkg/validate"
	"github.com/mstgnz/cronjob/schedule"
)

var (
	PORT                   string
	webUserHandler         web.UserHandler
	webHomeHandler         web.HomeHandler
	webScheduleHandler     web.ScheduleHandler
	webRequestHandler      web.RequestHandler
	webGroupHandler        web.GroupHandler
	webWebhookHandler      web.WebhookHandler
	webNotificationHandler web.NotificationHandler
	webSettingHandler      web.SettingHandler

	apiUserHandler           api.UserHandler
	apiGroupHandler          api.GroupHandler
	apiRequestHandler        api.RequestHandler
	apiRequestHeadderHandler api.RequestHeaderHandler
	apiNotificationHandler   api.NotificationHandler
	apiNotifyEmailHandler    api.NotifyEmailHandler
	apiNotifyMessageHandler  api.NotifyMessageHandler
	apiScheduleHandler       api.ScheduleHandler
	apiWebhookHandler        api.WebhookHandler
)

func init() {
	// Load Env
	if err := godotenv.Load(".env"); err != nil {
		logger.Warn(fmt.Sprintf("Load Env Error: %v", err))
		log.Fatalf("Load Env Error: %v", err)
	}
	// init conf
	_ = config.App()
	validate.CustomValidate()

	// Load Sql
	config.App().QUERY = make(map[string]string)
	if query, err := load.LoadSQLQueries(); err != nil {
		logger.Warn(fmt.Sprintf("Load Sql Error: %v", err))
		log.Fatalf("Load Sql Error: %v", err)
	} else {
		config.App().QUERY = query
	}

	PORT = os.Getenv("APP_PORT")
}

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func Catch(h HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			logger.Info("HTTP Handler Error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func main() {

	// test mail with attach
	//err := config.App().Mail.SetSubject("cron").SetContent("test mail").SetTo("mesutgenez@hotmail.com").SetAttachment(map[string][]byte{"query.sql": []byte("1. folder content"), "query2.sql": []byte("2. folder content")}).SendText()

	// Scheduler Call
	schedule.CallSchedule(config.App().Cron)
	config.App().Cron.Start()

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
		http.ServeFile(w, r, "./views/swagger.html")
	})

	// test route
	r.Get("/trigger", Catch(webHomeHandler.TriggerHandler))

	// web without auth
	r.Group(func(r chi.Router) {
		r.Use(isAuthMiddleware)
		r.Get("/login", Catch(webUserHandler.LoginHandler))
		r.Post("/login", Catch(webUserHandler.LoginHandler))
		r.Get("/register", Catch(webUserHandler.RegisterHandler))
		r.Post("/register", Catch(webUserHandler.RegisterHandler))
	})

	// web with auth
	r.Group(func(r chi.Router) {
		r.Use(webAuthMiddleware)
		r.Get("/", Catch(webHomeHandler.HomeHandler))
		r.Get("/logout", Catch(webUserHandler.LogoutHandler))
		r.Get("/profile", Catch(webUserHandler.ProfileHandler))
		// schedule
		r.Get("/schedules", Catch(webScheduleHandler.HomeHandler))
		r.Get("/schedules-pagination", Catch(webScheduleHandler.PaginationHandler))
		r.Post("/schedules", Catch(webScheduleHandler.CreateHandler))
		r.Get("/schedules/{id}", Catch(webScheduleHandler.EditHandler))
		r.Put("/schedules/{id}", Catch(webScheduleHandler.UpdateHandler))
		r.Delete("/schedules/{id}", Catch(webScheduleHandler.DeleteHandler))
		r.Get("/schedules/logs", Catch(webScheduleHandler.LogPaginationHandler))
		// request
		r.Get("/requests", Catch(webRequestHandler.HomeHandler))
		r.Get("/requests-pagination", Catch(webRequestHandler.PaginationHandler))
		r.Post("/requests", Catch(webRequestHandler.CreateHandler))
		r.Get("/requests/{id}", Catch(webRequestHandler.EditHandler))
		r.Put("/requests/{id}", Catch(webRequestHandler.UpdateHandler))
		r.Delete("/requests/{id}", Catch(webRequestHandler.DeleteHandler))
		r.Get("/requests/headers-pagination", Catch(webRequestHandler.HeaderPaginationHandler))
		r.Post("/requests/headers", Catch(webRequestHandler.HeaderCreateHandler))
		r.Get("/requests/headers/{id}", Catch(webRequestHandler.HeaderEditHandler))
		r.Put("/requests/headers/{id}", Catch(webRequestHandler.HeaderUpdateHandler))
		r.Delete("/requests/headers/{id}", Catch(webRequestHandler.HeaderDeleteHandler))
		// group
		r.Get("/groups", Catch(webGroupHandler.HomeHandler))
		r.Get("/groups-pagination", Catch(webGroupHandler.PaginationHandler))
		r.Post("/groups", Catch(webGroupHandler.CreateHandler))
		r.Get("/groups/{id}", Catch(webGroupHandler.EditHandler))
		r.Put("/groups/{id}", Catch(webGroupHandler.UpdateHandler))
		r.Delete("/groups/{id}", Catch(webGroupHandler.DeleteHandler))
		// webhook
		r.Get("/webhooks", Catch(webWebhookHandler.HomeHandler))
		r.Get("/webhooks-pagination", Catch(webWebhookHandler.PaginationHandler))
		r.Post("/webhooks", Catch(webWebhookHandler.CreateHandler))
		r.Get("/webhooks/{id}", Catch(webWebhookHandler.EditHandler))
		r.Put("/webhooks/{id}", Catch(webWebhookHandler.UpdateHandler))
		r.Delete("/webhooks/{id}", Catch(webWebhookHandler.DeleteHandler))
		// notification
		r.Get("/notifications", Catch(webNotificationHandler.HomeHandler))
		r.Get("/notifications-pagination", Catch(webNotificationHandler.PaginationHandler))
		r.Post("/notifications", Catch(webNotificationHandler.CreateHandler))
		r.Get("/notifications/{id}", Catch(webNotificationHandler.EditHandler))
		r.Put("/notifications/{id}", Catch(webNotificationHandler.UpdateHandler))
		r.Delete("/notifications/{id}", Catch(webNotificationHandler.DeleteHandler))
		r.Get("/notifications/email-pagination", Catch(webNotificationHandler.EmailPaginationHandler))
		r.Post("/notifications/email", Catch(webNotificationHandler.EmailCreateHandler))
		r.Get("/notifications/email/{id}", Catch(webNotificationHandler.EmailEditHandler))
		r.Put("/notifications/email/{id}", Catch(webNotificationHandler.EmailUpdateHandler))
		r.Delete("/notifications/email/{id}", Catch(webNotificationHandler.EmailDeleteHandler))
		r.Get("/notifications/message-pagination", Catch(webNotificationHandler.MessagePaginationHandler))
		r.Post("/notifications/message", Catch(webNotificationHandler.MessageCreateHandler))
		r.Get("/notifications/message/{id}", Catch(webNotificationHandler.MessageEditHandler))
		r.Put("/notifications/message/{id}", Catch(webNotificationHandler.MessageUpdateHandler))
		r.Delete("/notifications/message/{id}", Catch(webNotificationHandler.MessageDeleteHandler))
		// setting
		r.Route("/settings", func(r chi.Router) {
			r.Use(isAdminMiddleware)
			r.Get("/", Catch(webSettingHandler.HomeHandler))
			r.Get("/users", Catch(webSettingHandler.UsersHandler))
			r.Get("/app-logs", Catch(webSettingHandler.AppLogHandler))
			r.Post("/users/signup", Catch(webSettingHandler.UserSignUpHandler))
			r.Put("/users/change-profile", Catch(webSettingHandler.UserChangeProfileHandler))
			r.Put("/users/change-password", Catch(webSettingHandler.UserChangePasswordHandler))
			r.Delete("/users/{id}", Catch(webSettingHandler.UserDeleteHandler))
			r.Delete("/app-logs/{id}", Catch(webSettingHandler.AppLogDeleteHandler))
		})
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
			r.Get("/user", Catch(apiUserHandler.ProfileHandler))
			r.Put("/user", Catch(apiUserHandler.UpdateHandler))
			r.Delete("/user/{id}", Catch(apiUserHandler.DeleteHandler))
			r.Put("/user-change-pass", Catch(apiUserHandler.PassUpdateHandler))
			// groups
			r.Get("/groups", Catch(apiGroupHandler.ListHandler))
			r.Post("/groups", Catch(apiGroupHandler.CreateHandler))
			r.Put("/groups/{id}", Catch(apiGroupHandler.UpdateHandler))
			r.Delete("/groups/{id}", Catch(apiGroupHandler.DeleteHandler))
			// requests
			r.Get("/requests", Catch(apiRequestHandler.ListHandler))
			r.Post("/requests", Catch(apiRequestHandler.CreateHandler))
			r.Post("/requests/bulk", Catch(apiRequestHandler.BulkHandler))
			r.Put("/requests/{id}", Catch(apiRequestHandler.UpdateHandler))
			r.Delete("/requests/{id}", Catch(apiRequestHandler.DeleteHandler))
			// request headers
			r.Get("/request-headers", Catch(apiRequestHeadderHandler.ListHandler))
			r.Post("/request-headers", Catch(apiRequestHeadderHandler.CreateHandler))
			r.Put("/request-headers/{id}", Catch(apiRequestHeadderHandler.UpdateHandler))
			r.Delete("/request-headers/{id}", Catch(apiRequestHeadderHandler.DeleteHandler))
			// notifications
			r.Get("/notifications", Catch(apiNotificationHandler.ListHandler))
			r.Post("/notifications", Catch(apiNotificationHandler.CreateHandler))
			r.Post("/notifications/bulk", Catch(apiNotificationHandler.BulkHandler))
			r.Put("/notifications/{id}", Catch(apiNotificationHandler.UpdateHandler))
			r.Delete("/notifications/{id}", Catch(apiNotificationHandler.DeleteHandler))
			// notification emails
			r.Get("/notify-emails", Catch(apiNotifyEmailHandler.ListHandler))
			r.Post("/notify-emails", Catch(apiNotifyEmailHandler.CreateHandler))
			r.Put("/notify-emails/{id}", Catch(apiNotifyEmailHandler.UpdateHandler))
			r.Delete("/notify-emails/{id}", Catch(apiNotifyEmailHandler.DeleteHandler))
			// notification message
			r.Get("/notify-messages", Catch(apiNotifyMessageHandler.ListHandler))
			r.Post("/notify-messages", Catch(apiNotifyMessageHandler.CreateHandler))
			r.Put("/notify-messages/{id}", Catch(apiNotifyMessageHandler.UpdateHandler))
			r.Delete("/notify-messages/{id}", Catch(apiNotifyMessageHandler.DeleteHandler))
			// webhooks
			r.Get("/webhooks", Catch(apiWebhookHandler.ListHandler))
			r.Post("/webhooks", Catch(apiWebhookHandler.CreateHandler))
			r.Put("/webhooks/{id}", Catch(apiWebhookHandler.UpdateHandler))
			r.Delete("/webhooks/{id}", Catch(apiWebhookHandler.DeleteHandler))
			// schedules
			r.Get("/schedules", Catch(apiScheduleHandler.ListHandler))
			r.Post("/schedules", Catch(apiScheduleHandler.CreateHandler))
			r.Post("/schedules/bulk", Catch(apiScheduleHandler.BulkHandler))
			r.Put("/schedules/{id}", Catch(apiScheduleHandler.UpdateHandler))
			r.Delete("/schedules/{id}", Catch(apiScheduleHandler.DeleteHandler))
			// schedule logs
			r.Get("/schedule-logs", Catch(apiScheduleHandler.LogListHandler))
		})
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "api") {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: "Not Found"})
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Run your HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
		if err != nil && err != http.ErrServerClosed {
			logger.Warn("Fatal Error", "err", err.Error())
			log.Fatal(err.Error())
		}
	}()

	logger.Info("Cron is running on", PORT)

	// Block until a signal is received
	<-ctx.Done()

	logger.Info("Cron is shutting on", PORT)

	// set Shutting
	config.App().Shutting = true

	// check Running
	for {
		if config.App().Running <= 0 {
			logger.Info("Cronjobs all done")
			break
		} else {
			logger.Info(fmt.Sprintf("Currently %d active jobs in progress. pending completion...", config.App().Running))
		}
		time.Sleep(time.Second * 5)
	}

	config.App().Cron.Stop()
	logger.Info("Shutting down gracefully...")
	config.App().DB.CloseDatabase()
}

func webAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		token := strings.Replace(cookie.Value, "Bearer ", "", 1)

		userId, err := auth.GetUserIDByToken(token)
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
		token := r.Header.Get("Authorization")
		if token == "" {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: "Invalid Token"})
			return
		}
		token = strings.Replace(token, "Bearer ", "", 1)

		userId, err := auth.GetUserIDByToken(token)
		if err != nil {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: err.Error()})
			return
		}

		user_id, err := strconv.Atoi(userId)
		if err != nil && user_id == 0 {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: err.Error()})
			return
		}

		user := &models.User{}
		err = user.GetWithId(user_id)

		if err != nil {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")

		if err == nil {
			token := strings.Replace(cookie.Value, "Bearer ", "", 1)
			_, err = auth.GetUserIDByToken(token)
			if err == nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cUser, ok := r.Context().Value(config.CKey("user")).(*models.User)

		if !cUser.IsAdmin || !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkMethod := r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH"
		if checkMethod && r.Header.Get("Content-Type") != "application/json" {
			_ = response.WriteJSON(w, http.StatusBadRequest, response.Response{Status: false, Message: "Invalid Content-Type"})
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
