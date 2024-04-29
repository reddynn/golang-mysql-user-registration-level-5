package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"reddynn/config"
	"reddynn/models"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

type session struct {
	username string
	expiry   time.Time
}

var sessions = map[string]session{}

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to homepage"))

}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&user)

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("Syntax error: %s", err), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		// Validation failed, handle the error
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	dbs, err := config.Dbconnect()
	if err != nil {
		logger.Error(err.Error())
	}
	defer dbs.Close()
	var newuser string
	err = dbs.QueryRow("select username from users where username=?", user.Username).Scan(&newuser)
	switch {
	case err == sql.ErrNoRows:
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "unable to create account", http.StatusInternalServerError)
			return
		}
		_, err = dbs.Query("insert into users(username,password,email) values(?,?,?)", user.Username, hashedpassword, user.Email)
		if err != nil {
			http.Error(w, "server unble to create user", http.StatusInternalServerError)
			logger.Error(err.Error())
			return

		}
		w.Write([]byte("user has been created"))
		logger.Info("?: user has been created", user.Username, "users")
	case err != nil:
		http.Error(w, "server db query error", http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	default:
		http.Error(w, "user already exsited", http.StatusBadRequest)

		return
	}

}

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user models.User
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("Syntax error: %s", err), http.StatusBadRequest)
		return
	}

	dbs, err := config.Dbconnect()
	if err != nil {
		logger.Error(err.Error())
	}
	defer dbs.Close()
	var newpassword string
	err = dbs.QueryRow("select password from users where username=?", user.Username).Scan(&newpassword)
	switch {
	case err != nil:
		http.Error(w, "unauthorized no user with this username", http.StatusUnauthorized)
		return
	default:
		err := bcrypt.CompareHashAndPassword([]byte(newpassword), []byte(user.Password))
		if err != nil {
			http.Error(w, "unathorized password wrong", http.StatusUnauthorized)
			return
		}
		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)
		sessions[sessionToken] = session{
			username: user.Username,
			expiry:   expiresAt,
		}
		
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		
		w.Write([]byte("welcome" + user.Username))
	}

}
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func Profile(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {

			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.username)))
}

func Signout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}
