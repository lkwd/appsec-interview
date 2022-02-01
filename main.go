package main

import (
	"appsec-interview/fileops"
	"appsec-interview/token"
	"appsec-interview/users"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	authprefix = "Bearer "
	authPrefixLen = 7
)

// this is the secret key used to sign session tokens. It should be 32 bytes
var secret = func()[]byte{
	b, err := ioutil.ReadFile("./secret.txt")
	if err != nil {
		panic(fmt.Sprintf("failed to read secret key: %s", err.Error()))
	}
	return b
}()

// this is the location where the files we serve/move are
var fileDir = func()string{
	f, err := filepath.Abs("./files")
	if err != nil {
		panic(fmt.Sprintf("failed to get absolute path to file dir: %s", err.Error()))
	}
	return f
}()

type RenameRequest struct {
	Old string
	New string
}

// allows the admin user to rename files
func renameHandler(w http.ResponseWriter, r *http.Request) {
	t, ok := r.Context().Value("token").(*token.Token)
	if !ok {
		panic("unauthenticated call to authenticated route")
	}
	if t.User != "admin" {
		w.WriteHeader(403)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Add("X-Error", err.Error())
		w.WriteHeader(400)
		return
	}

	data := RenameRequest{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	o, err := fileops.Rename(fileDir, data.Old, data.New)
	if err != nil {
		w.Header().Add("X-Error", err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write(o)
}


type LoginRequest struct {
	Name string
	Password string
}

// allows an anonymous user to log in
func loginHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Add("X-Error", err.Error())
		w.WriteHeader(400)
		return
	}

	data := LoginRequest{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	user := users.Database().Get(data.Name)
	if user == nil {
		w.WriteHeader(401)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		w.WriteHeader(401)
		return
	}

	t, err := token.Generate(secret, time.Hour * 6, user.Name)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Add("X-Auth-Token", t)
	w.WriteHeader(200)
}

// ensures the user is properly authenticated before granting access to sensitive functionality
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Header.Get("Authorization")
		if len(v) < authPrefixLen {
			w.WriteHeader(401)
			return
		}

		t, err := token.Parse(v[authPrefixLen:])
		if err != nil {
			w.WriteHeader(401)
			return
		}

		err = t.Validate(secret)
		if err != nil {
			w.Header().Add("X-Error", err.Error())
			w.WriteHeader(401)
			return
		}
		req := r.WithContext(context.WithValue(r.Context(), "token", t))
		next.ServeHTTP(w, req)
	})
}

// Route the request to the appropriate handler
func handler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if strings.HasPrefix(p, "/login") {
		loginHandler(w, r)
	} else if strings.HasPrefix(p, "/rename") {
		authenticate(http.HandlerFunc(renameHandler)).ServeHTTP(w, r)
	} else if strings.HasPrefix(p, "/view") {
		h := authenticate(http.FileServer(http.Dir(fileDir)))
		http.StripPrefix("/view", h).ServeHTTP(w, r)
	} else {
		w.WriteHeader(404)
	}
}

func main() {
	println("starting file server on 0.0.0.0:9555")
	if err := http.ListenAndServe(":9555", http.HandlerFunc(handler)); err != nil {
		log.Fatalln(err)
	}
}
