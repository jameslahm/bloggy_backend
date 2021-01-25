package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jameslahm/bloggy_backend/api/auth"
	"github.com/jameslahm/bloggy_backend/api/responses"
	"github.com/jameslahm/bloggy_backend/api/utils"
	"github.com/jameslahm/bloggy_backend/models"
)

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// TODO: validate form
	err = models.CreateUser(server.DB, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, utils.FormatError(err))
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user, err := models.FindUserById(server.DB, int(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, *user)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	authId := r.Context().Value("id")
	if authId == nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if authId != id {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Internal Error"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var obj map[string]interface{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = models.UpdateUser(server.DB, int(id), obj)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, "ok")
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	authId := r.Context().Value("id")
	if authId == nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Internal Error"))
		return
	}

	if authId != id {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	err = models.DeleteUser(server.DB, int(id))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, "ok")
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var obj map[string]string
	var user = models.User{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = server.DB.Debug().Model(&models.User{}).Where("email=?", obj["email"]).Take(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = models.VerifyPassword(user.Password, obj["password"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user.Token, err = auth.CreateToken(int(user.ID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, user)
}
