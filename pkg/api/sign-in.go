package api

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ionov-egor/go_todo_v2/pkg/config"
	"net/http"
)

type signInInf struct {
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var signInInf signInInf
	var buf bytes.Buffer
	var err error

	if _, err := buf.ReadFrom(r.Body); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &signInInf); err != nil {
		HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	envPassword := config.GetEnv("TODO_PASSWORD", "12345")
	if signInInf.Password != envPassword {
		HandleError(w, r, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	signedToken, err := jwtToken.SignedString([]byte("my_secret_key"))
	if err != nil {
		HandleError(w, r, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	WriteJson(w, r, http.StatusOK, struct {
		Token string `json:"token"`
	}{signedToken})
	return
	
}
