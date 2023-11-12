package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"reflect"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func responseWithMsg(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	res := Response{
		Message: msg,
		Code:    statusCode,
	}
	_ = jsoniter.NewEncoder(w).Encode(res)
	w.Header().Set("Content-Type", "application/json")
}

func responseWithJSON(w http.ResponseWriter, statusCode int, data interface{}, msg string) {
	w.WriteHeader(statusCode)
	res := Response{
		Message: msg,
		Code:    statusCode,
		Data:    data,
	}
	json.NewEncoder(w).Encode(res)
	w.Header().Set("Content-Type", "application/json")
}

func Reader(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	return io.ReadAll(r.Body)
}

func EmptyFields(data interface{}) (emptyFields []string) {
	for i := 0; i < reflect.ValueOf(data).NumField(); i++ {
		field := reflect.ValueOf(data).Field(i)
		kind := field.Kind()
		switch kind {
		case reflect.String:
			if field.String() == "" {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Int:
			if field.Int() == 0 {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Float64:
			if field.Float() == 0 {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Bool:
			if !field.Bool() {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Struct:
			if field.IsZero() {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Ptr:
			if field.IsNil() {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Slice:
			if field.IsNil() {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		case reflect.Map:
			if field.IsNil() {
				emptyFields = append(emptyFields, reflect.TypeOf(data).Field(i).Name)
			}
		}
	}
	return
}

func Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func Decode(data []byte, c interface{}) (interface{}, error) {
	fmt.Println("Decode")
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Update(w http.ResponseWriter, r *http.Request, t Tables) bool {
	id := chi.URLParam(r, "id")
	if id == "" || id == "0" || id == " " || len(id) == 0 || id == "null" || id == "undefined" || id == "NaN" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return false
	}

	var data = t.Struct()
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return false
	}

	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return false
	}

	if err = json.Unmarshal(body, &data); err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return false
	}

	validationError := Validate(data)
	if validationError.Message != "" {
		responseWithMsg(w, http.StatusBadRequest, validationError.Message)
		return false
	}

	result = db.Table(t.String()).Updates(data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return false
	}

	responseWithJSON(w, http.StatusOK, data, fmt.Sprintf("ID %s updated from %s", id, t.String()))
	return true
}

func Create(w http.ResponseWriter, r *http.Request, t Tables) bool {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return false
	}

	var data = t.Struct()
	if err = json.Unmarshal(body, &data); err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return false
	}

	validationError := Validate(data)
	if validationError.Message != "" {
		responseWithMsg(w, http.StatusBadRequest, validationError.Message)
		return false
	}

	result := db.Table(t.String()).Create(data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return false
	}

	responseWithJSON(w, http.StatusOK, data, fmt.Sprintf("New %s created", t.String()))
	return true
}

func Read(w http.ResponseWriter, r *http.Request, t Tables) bool {
	var data = t.Slice()
	result := db.Table(t.String()).Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return false
	}

	responseWithJSON(w, http.StatusOK, data, fmt.Sprintf("%s read", t.String()))
	return true
}

func ReadOne(w http.ResponseWriter, r *http.Request, t Tables) bool {
	id := chi.URLParam(r, "id")
	if id == "" || id == "0" || id == " " || len(id) == 0 || id == "null" || id == "undefined" || id == "NaN" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return false
	}

	var data = t.Struct()
	result := db.Table(t.String()).First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return false
	}

	responseWithJSON(w, http.StatusOK, data, fmt.Sprintf("ID %s read from %s", id, t.String()))
	return true
}

func Delete(w http.ResponseWriter, r *http.Request, t Tables) bool {
	id := chi.URLParam(r, "id")
	if id == "" || id == "0" || id == " " || len(id) == 0 || id == "null" || id == "undefined" || id == "NaN" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return false
	}

	var data = t.Struct()
	result := db.Table(t.String()).Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return false
	}

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("ID %s deleted from %s", id, t.String()))
	return true
}
