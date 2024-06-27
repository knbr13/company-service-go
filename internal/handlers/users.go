package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/knbr13/company-service-go/config"
	jwtUtil "github.com/knbr13/company-service-go/internal/jwt"
	"github.com/knbr13/company-service-go/internal/repositories"
	"github.com/knbr13/company-service-go/internal/services"
	"github.com/knbr13/company-service-go/internal/validator"
	"github.com/knbr13/company-service-go/pkg/util"
)

type UserHandler struct {
	Services *services.Services
	Cfg      *config.Config
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrJsonResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := &repositories.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err = uh.Services.Users.Register(ctx, user)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateEmail) {
			util.ErrJsonResponse(w, http.StatusConflict, fmt.Sprintf("email: %s", err.Error()))
			return
		}
		if e, ok := err.(validator.ValidationError); ok {
			util.ErrJsonResponse(w, http.StatusBadRequest, e.Error())
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error registering user: %s\n", err.Error())
		return
	}

	token, err := jwtUtil.GenerateToken(jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{ // may need to add user id into the claims in the future
			Time: time.Now().Add(time.Hour * 24),
		},
	}, []byte(uh.Cfg.JWTKey))
	if err != nil {
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error generating token: %s\n", err.Error())
		return
	}

	util.JsonResponse(w, http.StatusCreated, map[string]any{
		"token": token,
	})
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.ErrJsonResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := &repositories.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err = uh.Services.Users.Login(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.ErrJsonResponse(w, http.StatusNotFound, "email: not found")
			return
		}
		if errors.Is(err, repositories.ErrInvalidPassword) {
			util.ErrJsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("password: %s", err.Error()))
			return
		}
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error logging in user: %s\n", err.Error())
		return
	}

	token, err := jwtUtil.GenerateToken(jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{ // may need to add user id into the claims in the future
			Time: time.Now().Add(time.Hour * 24),
		},
	}, []byte(uh.Cfg.JWTKey))
	if err != nil {
		util.ErrJsonResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		log.Printf("error generating token: %s\n", err.Error())
		return
	}

	util.JsonResponse(w, http.StatusCreated, map[string]any{
		"token": token,
	})
}
