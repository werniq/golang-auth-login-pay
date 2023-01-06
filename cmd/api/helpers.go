package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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