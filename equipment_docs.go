package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type EquipmentDoc struct {
	gorm.Model
	EquipmentID uint      `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Equipment   Equipment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DocName     string    `gorm:"varchar(255);not null" validate:"required,max=255"`
	DocURL      string    `gorm:"varchar(255);not null" validate:"required,max=255"`
	UploadDate  time.Time `gorm:"not null" validate:"required,datetime"`
}

func (c *EquipmentDoc) Decode(data []byte) (EquipmentDoc, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return EquipmentDoc{}, err
	}
	return *c, nil
}

func (c *EquipmentDoc) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func equipmentDocCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c EquipmentDoc
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

	responseWithJSON(w, http.StatusOK, data, "equipment doc created")
	return
}

func equipmentDocReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []EquipmentDoc
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "equipment doc read")
	return
}

func equipmentDocReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data EquipmentDoc
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "equipment doc read")
	return
}

func equipmentDocUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data EquipmentDoc
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

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("equipment doc with id %s updated", id))
	return
}

func equipmentDocDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return
	}

	var data EquipmentDoc
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("equipment doc with id %s deleted", id))
	return
}
