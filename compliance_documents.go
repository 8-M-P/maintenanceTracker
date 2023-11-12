package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ComplianceDocument struct {
	gorm.Model
	EquipmentID  uint      `gorm:"type:int(10);index;not null" validate:"required"`
	Equipment    Equipment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"omitempty"`
	DocumentName string    `gorm:"type:varchar(255);not null" validate:"required"`
	DocumentURL  string    `gorm:"type:varchar(255)"`
	ExpiryDate   time.Time `gorm:"type:date;not null" validate:"required,datetime"`
}

func (c *ComplianceDocument) Decode(data []byte) (ComplianceDocument, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return ComplianceDocument{}, err
	}
	return *c, nil
}

func (c *ComplianceDocument) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func complianceDocumentCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c ComplianceDocument
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

	responseWithJSON(w, http.StatusOK, data, "compliance document created")
	return
}

func complianceDocumentReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []ComplianceDocument
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "compliance document read")
	return
}

func complianceDocumentReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return
	}

	var data ComplianceDocument
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "compliance document read")
	return
}

func complianceDocumentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data ComplianceDocument
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

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("compliance document with id %s updated", id))
	return
}

func complianceDocumentDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return
	}

	var data ComplianceDocument
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithMsg(w, http.StatusOK, fmt.Sprintf("compliance document with id %s deleted", id))
	return
}
