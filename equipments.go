package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Equipment struct {
	gorm.Model
	CompanyID           uint              `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Company             Company           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EquipmentCategoryID uint              `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	EquipmentCategory   EquipmentCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name                string            `gorm:"varchar(255);not null" validate:"required,max=255"`
	PurchaseDate        time.Time         `gorm:"not null" validate:"required,datetime"`
	WarrantyExpiry      time.Time         `gorm:"not null" validate:"required,datetime"`
	LastMaintenanceDate time.Time         `gorm:"not null" validate:"required,datetime"`
	ImageURL            string            `gorm:"varchar(255);" validate:"max=255"`
	AdditionalNotes     string            `gorm:"varchar(500);" validate:"max=500"`
}

func (c *Equipment) Decode(data []byte) (Equipment, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return Equipment{}, err
	}
	return *c, nil
}

func (c *Equipment) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func equipmentCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c Equipment
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

	responseWithJSON(w, http.StatusOK, data, "equipment created")
	return
}

func equipmentReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []Equipment
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, data, "equipment read")
	return
}

func equipmentReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Equipment
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, data, "equipment read")
	return
}

func equipmentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Equipment
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

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("equipment with id %s updated", id))
	return
}

func equipmentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return
	}

	var data Equipment
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("equipment with id %s deleted", id))
	return
}
