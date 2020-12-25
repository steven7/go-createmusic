package controllers

import (

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/steven7/go-createmusic/models"
)

//
// From methods -- deprecated
//

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}

func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

//
// Get the JSON body and decode into interface
//
func ParseJSONParameters(w http.ResponseWriter, r *http.Request, data interface{}) error {

	fmt.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		errorData := models.Error {
			Title:  "Invalid or missing parameters",
			Detail: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Could not parse parameters")
		WriteJson(w, errorData)
		return err
	}

	return nil // good case
}

