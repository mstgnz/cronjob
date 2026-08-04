package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/cronjob/models"
	"github.com/mstgnz/cronjob/pkg/auth"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/mstgnz/cronjob/pkg/logger"
	"github.com/mstgnz/cronjob/pkg/response"
	"github.com/mstgnz/cronjob/pkg/validate"
)

type UserService struct{}

// Login attempts are limited per address and per account. Both are needed: the
// address bucket stops one host walking a password list, the account bucket stops
// a distributed attempt against a single user.
var (
	loginIPLimiter    = auth.NewLimiter(10, 15*time.Minute)
	loginEmailLimiter = auth.NewLimiter(5, 15*time.Minute)
)

// invalidCredentials is returned for every failed login. A distinct message or
// status for an unknown address would tell an attacker which accounts exist.
const invalidCredentials = "Invalid credentials"

func (s *UserService) LoginService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	login := &models.Login{}
	if err := response.ReadJSON(w, r, login); err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Invalid request body"}
	}

	err := validate.Validate(login)
	if err != nil {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Content validation invalid", Data: map[string]any{"error": err.Error()}}
	}

	clientIP := auth.ClientIP(r)
	emailKey := strings.ToLower(strings.TrimSpace(login.Email))
	if !loginIPLimiter.Allow(clientIP) || !loginEmailLimiter.Allow(emailKey) {
		logger.Warn("Login rate limited", fmt.Sprintf("ip=%s", clientIP))
		return http.StatusTooManyRequests, response.Response{Status: false, Message: "Too many attempts, try again later"}
	}

	user := &models.User{}
	err = user.GetWithMail(login.Email)
	if err != nil {
		// The password is still hashed against a dummy value so an unknown address
		// does not answer measurably faster than a known one.
		auth.ComparePassword(auth.DummyHash, login.Password)
		return http.StatusBadRequest, response.Response{Status: false, Message: invalidCredentials}
	}

	if !auth.ComparePassword(user.Password, login.Password) {
		return http.StatusBadRequest, response.Response{Status: false, Message: invalidCredentials}
	}

	if !user.Active {
		return http.StatusBadRequest, response.Response{Status: false, Message: invalidCredentials}
	}

	loginIPLimiter.Reset(clientIP)
	loginEmailLimiter.Reset(emailKey)

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		logger.Warn("Token Generate Error", err.Error())
		return http.StatusInternalServerError, response.Response{Status: false, Message: "Failed to generate token"}
	}

	// update last_login
	user.LastLoginUpdate()

	data := make(map[string]any)
	data["token"] = token
	data["user"] = user
	return http.StatusOK, response.Response{Status: true, Message: "Login successful", Data: data}
}

func (s *UserService) RegisterService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	register := &models.Register{}
	if err := response.ReadJSON(w, r, register); err != nil {
		return http.StatusBadRequest, badRequestError(err)
	}

	err := validate.Validate(register)
	if err != nil {
		return http.StatusBadRequest, serverError(err)
	}

	user := &models.User{}
	exists, err := user.Exists(register.Email)
	if err != nil {
		return http.StatusInternalServerError, serverError(err)
	}
	if exists {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Email already exists"}
	}

	err = user.Create(register)
	if err != nil {
		return http.StatusInternalServerError, serverError(err)
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return http.StatusInternalServerError, response.Response{Status: false, Message: "Failed to generate token"}
	}

	data := make(map[string]any)
	data["token"] = token
	data["user"] = user
	return http.StatusCreated, response.Response{Status: true, Message: "User created", Data: data}
}

func (s *UserService) ProfileService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	user := r.Context().Value(config.CKey("user"))
	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"user": user}}
}

func (s *UserService) Users(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	user := &models.User{}
	count := user.Count()
	users := user.Get(0, 20, "")
	return http.StatusOK, response.Response{Status: count > 0, Message: "Success", Data: map[string]any{"count": count, "users": users}}
}

func (s *UserService) UpdateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	updateData := &models.ProfileUpdate{}
	if err := response.ReadJSON(w, r, updateData); err != nil {
		return http.StatusBadRequest, badRequestError(err)
	}

	err := validate.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, serverError(err)
	}

	user := &models.User{}
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)
	user.ID = cUser.ID

	// If the administrator wants to update a user.
	if updateData.ID > 0 && cUser.IsAdmin {
		user.ID = updateData.ID
	}

	queryParts := []string{"UPDATE users SET"}
	params := []any{}
	paramCount := 1

	if updateData.Fullname != "" {
		queryParts = append(queryParts, fmt.Sprintf("fullname=$%d,", paramCount))
		params = append(params, updateData.Fullname)
		paramCount++
	}
	if updateData.Email != "" {
		// check email if not same email
		if err := user.GetWithId(user.ID); err != nil {
			return http.StatusInternalServerError, serverError(err)
		}
		if user.Email != updateData.Email {
			exists, err := user.Exists(updateData.Email)
			if err != nil {
				return http.StatusInternalServerError, serverError(err)
			}
			if exists {
				return http.StatusBadRequest, response.Response{Status: false, Message: "Email already exists"}
			}
		}
		queryParts = append(queryParts, fmt.Sprintf("email=$%d,", paramCount))
		params = append(params, updateData.Email)
		paramCount++
	}
	if updateData.Phone != "" {
		queryParts = append(queryParts, fmt.Sprintf("phone=$%d,", paramCount))
		params = append(params, updateData.Phone)
		paramCount++
	}

	if len(params) == 0 {
		return http.StatusBadRequest, response.Response{Status: false, Message: "No fields to update"}
	}

	// update at
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	queryParts = append(queryParts, fmt.Sprintf("updated_at=$%d", paramCount))
	params = append(params, updatedAt)
	paramCount++

	// remove last comma
	// size := len(queryParts) - 1
	// queryParts[size] = strings.TrimSuffix(queryParts[size], ",")

	queryParts = append(queryParts, fmt.Sprintf("WHERE id=$%d", paramCount))
	params = append(params, user.ID)
	query := strings.Join(queryParts, " ")

	err = user.ProfileUpdate(query, params)

	if err != nil {
		return http.StatusInternalServerError, serverError(err)
	}

	return http.StatusOK, response.Response{Status: true, Message: "Success", Data: map[string]any{"update": updateData}}
}

func (s *UserService) PassUpdateService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	updateData := &models.PasswordUpdate{}
	if err := response.ReadJSON(w, r, updateData); err != nil {
		return http.StatusBadRequest, badRequestError(err)
	}

	err := validate.Validate(updateData)
	if err != nil {
		return http.StatusBadRequest, serverError(err)
	}

	if updateData.Password != updateData.RePassword {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Passwords do not match"}
	}

	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	user := &models.User{}
	user.ID = cUser.ID

	// If the administrator wants to update a user.
	adminActingOnSomeoneElse := updateData.ID > 0 && cUser.IsAdmin && updateData.ID != cUser.ID
	if adminActingOnSomeoneElse {
		user.ID = updateData.ID
	} else {
		// Changing your own password requires proving you know the current one:
		// otherwise a momentarily borrowed session becomes permanent access.
		if updateData.CurrentPassword == "" {
			return http.StatusBadRequest, response.Response{Status: false, Message: "Current password is required"}
		}
		if !auth.ComparePassword(cUser.Password, updateData.CurrentPassword) {
			return http.StatusBadRequest, response.Response{Status: false, Message: "Current password is incorrect"}
		}
	}

	// the update also stamps tokens_valid_after, retiring sessions opened with the old password
	err = user.PasswordUpdate(updateData.Password)

	if err != nil {
		logger.Warn("Password Update Error", err.Error())
		return http.StatusInternalServerError, response.Response{Status: false, Message: "Password could not be updated"}
	}

	// the submitted passwords are deliberately not echoed back
	return http.StatusOK, response.Response{Status: true, Message: "Password updated"}
}

func (s *UserService) DeleteService(w http.ResponseWriter, r *http.Request) (int, response.Response) {
	// get auth user in context
	cUser, _ := r.Context().Value(config.CKey("user")).(*models.User)

	if !cUser.IsAdmin {
		return http.StatusForbidden, response.Response{Status: false, Message: "You're not a admin!"}
	}
	user := &models.User{}

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	// the row has to be loaded, not just probed: the guards below read its fields,
	// and they have to run before the delete, not after it
	if err := user.GetWithId(id); err != nil {
		return http.StatusNotFound, response.Response{Status: false, Message: "User not found"}
	}

	if user.ID == cUser.ID {
		return http.StatusBadRequest, response.Response{Status: false, Message: "You can't erase yourself!"}
	}
	if user.IsAdmin {
		return http.StatusBadRequest, response.Response{Status: false, Message: "Admin cannot delete admin!"}
	}

	if err := user.Delete(id); err != nil {
		return http.StatusInternalServerError, serverError(err)
	}

	return http.StatusOK, response.Response{Status: true, Message: "Soft delete success"}
}
