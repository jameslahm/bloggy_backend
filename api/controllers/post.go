package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jameslahm/bloggy_backend/api/responses"
	"github.com/jameslahm/bloggy_backend/models"
)

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	authId := r.Context().Value("id")
	if authId == nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Internal Error"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var post = models.Post{}
	err = json.Unmarshal(body, &post)
	post.AuthorID = authId.(uint32)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = models.CreatePost(server.DB, &post)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, post)
}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := models.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post, err := models.FindPostById(server.DB, int(id))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, post)
}

func (server *Server) UpdataPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	authId := r.Context().Value("id")
	post, err := models.FindPostById(server.DB, int(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if post.AuthorID != authId {
		responses.ERROR(w, http.StatusBadRequest, err)
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
	err = models.UpdatePost(server.DB, int(id), obj)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, "ok")
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	authId := r.Context().Value("id")
	post, err := models.FindPostById(server.DB, int(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if post.AuthorID != authId {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = models.DeletePost(server.DB, int(id))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, "ok")
}
