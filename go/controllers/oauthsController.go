package controllers

import (
	"context"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	llctx "github.com/steven7/go-createmusic/go/context"
	//"go-createmusic/dbx"
	"github.com/steven7/go-createmusic/models"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func NewOAuths(os models.OAuthService, configs map[string]*oauth2.Config) *OAuths {
	return &OAuths{
		os:        os,
		configs: configs,
	}
}

type OAuths struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}


func (o *OAuths) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}


// dbxRedirect
func (o *OAuths) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", http.StatusBadRequest)
		return
	}
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, service)
	if err == models.ErrNotFound {
		// noop
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}
	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: service,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", token)
}

/*
func (o *OAuths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOAuth, error := o.os.Find(user.ID, service)
	if error != nil {
		panic(error)
	}
	token := userOAuth.Token

	folders, files, err := dbx.List(token.AccessToken, path)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Folders=", folders)
	fmt.Fprintln(w, "Files=", files)
}
 */