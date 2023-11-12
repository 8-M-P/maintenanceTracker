package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type MaintenanceType struct {
	gorm.Model
	TypeName    string         `gorm:"varchar(255);unique;not null" validate:"required,max=255"`
	Description sql.NullString `gorm:"varchar(255);null" validate:"max=255"`
}

func (c *MaintenanceType) Decode(data []byte) (MaintenanceType, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return MaintenanceType{}, err
	}
	return *c, nil
}

func (c *MaintenanceType) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func maintenanceTypeCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c MaintenanceType
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

	responseWithJSON(w, http.StatusOK, data, "maintenance type created")
	return
}

func maintenanceTypeReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []MaintenanceType
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance type read")
	return
}

func maintenanceTypeReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceType
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance type read")
	return
}

func maintenanceTypeUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceType
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

	responseWithJSON(w, http.StatusOK, data, "maintenance type updated")
	return
}

func maintenanceTypeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceType
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance type deleted")
	return
}
