package main

import "encoding/json"

type InitiumHandler struct {
  Success bool `json:"success"`
  Redirect string `json:"redirect,omitempty"`
  Error string `json:"error,omitempty"`
  Fields FieldErrors `json:"fields,omitempty"`
}

type FieldErrors map[string]string

func (fields FieldErrors) Add(id string, err string) {
  if fields == nil {
    fields = make(map[string]string)
  }
  fields[id] = err
}

func (app *InitiumApp) RenderData(request *InitiumRequest, data interface{}) error {
  request.Writer.Header().Set("Content-Type", "application/json")
  json.NewEncoder(request.Writer).Encode(data)

  return nil
}