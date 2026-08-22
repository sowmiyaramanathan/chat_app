package controllers

import (
	"backend/auth"
	e "backend/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user e.User
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Println("body : ", r)
	if err != nil {
		http.Error(w, "Cannot Process your Request", http.StatusUnprocessableEntity)
		return
	}

	err = c.s.CreateUser(&user)
	if err != nil {
		if err.Error() == "username" {
			http.Error(w, `{"message": "Username"}`, http.StatusConflict)
			return
		} else if err.Error() == "mobile_number" {
			http.Error(w, `{"message": "Number"}`, http.StatusConflict)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("User Created Successfully"))
}

func (c *controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user e.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Cannot Process your Request", http.StatusUnprocessableEntity)
		return
	}

	userId, err := c.s.LoginUser(user.Username, user.Password)
	if err != nil {
		if err.Error() == "username" {
			http.Error(w, `{"message": "Username"}`, http.StatusConflict)
			return
		} else if err.Error() == "wrong_password" {
			http.Error(w, `{"message": "Password"}`, http.StatusConflict)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	tokenString := auth.CreateToken(userId, user.Username)

	if tokenString == "" {
		http.Error(w, "Precondition Failed", http.StatusPreconditionFailed)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Logged in successfully",
		"token":  tokenString,
	})
}

func (c *controller) Profile(w http.ResponseWriter, r *http.Request) {
	claims := auth.ExtractToken(r)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name": claims.Username,
	})
}

func (c *controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	claims := auth.ExtractToken(r)

	users, err := c.s.GetAllUsers(claims.Username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := users
	json.NewEncoder(w).Encode(data)
}

func (c *controller) GetPublicKey(w http.ResponseWriter, r *http.Request) {
	userIDString := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	publicKey, err := c.s.GetPublicKey(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := publicKey
	json.NewEncoder(w).Encode(data)
}
