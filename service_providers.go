package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type ServiceProvider struct {
	gorm.Model
	Name           string         `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	Contact        string         `gorm:"type:varchar(255)" validate:"max=255"`
	Rating         float32        `gorm:"type:decimal(2,1);default:0" validate:"max=5"`
	ReviewsCount   uint           `gorm:"int(10);default:0" validate:"alphanum,len=10"`
	Specialization sql.NullString `gorm:"type:varchar(500)" validate:"max=500"`
	Tags           string         `gorm:"type:json" validate:"json"`
	Address        string         `gorm:"type:varchar(500)" validate:"max=500"`
	Email          string         `gorm:"type:varchar(255)" validate:"max=255,email"`
}

func (c *ServiceProvider) Decode(data []byte) (ServiceProvider, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return ServiceProvider{}, err
	}
	return *c, nil
}

func (c *ServiceProvider) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func serviceProviderCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c ServiceProvider
	data, err := c.Decode(body)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	validationError := Validate(data)
	if validationError.Message != "" {
		responseWithMsg(w, http.StatusBadRequest, validationError.Message)
		return
	}

	result := db.Create(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "service provider created")
	return
}

func serviceProviderReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []ServiceProvider
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "service provider read")
	return
}

func serviceProviderReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data ServiceProvider
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "service provider read")
	return
}

func serviceProviderUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data ServiceProvider
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err = data.Decode(body)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	validationError := Validate(data)
	if validationError.Message != "" {
		responseWithMsg(w, http.StatusBadRequest, validationError.Message)
		return
	}

	result = db.Save(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "service provider updated")
	return
}

func serviceProviderDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data ServiceProvider
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "service provider deleted")
	return
}
