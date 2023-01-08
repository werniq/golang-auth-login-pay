package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var maxBytes int = 1048576
	
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)

	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only include single json value")
	}

	return nil
}

func (app *application) BadRequest(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Error bool `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

func (app *application) InvalidCredentials(w http.ResponseWriter) error {
	var payload struct {
		Error bool `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = true
	payload.Message = "invalid authentication credentials"

	err := app.writeJSON(w, http.StatusUnauthorized, payload)

	if err != nil {
		return err
	}

	return nil
}

func (app *application) ValidatePassword(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(request.Password), []byte(user.Password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default: 
			return false, err
		}
	}
	return true, nil
}