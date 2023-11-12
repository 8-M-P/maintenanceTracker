package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type Supplier struct {
	gorm.Model
	SupplierName   string  `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	ContactDetails string  `gorm:"type:varchar(255)" validate:"max=255"`
	Phone          string  `gorm:"type:varchar(255)" validate:"max=255,e164"`
	Address        string  `gorm:"type:varchar(500)" validate:"max=500"`
	Email          string  `gorm:"type:varchar(255)" validate:"max=255,email"`
	IBAN           string  `gorm:"type:varchar(255)" validate:"max=255"`
	Tags           *string `gorm:"type:json" validate:"json"`
}

func (c *Supplier) Decode(data []byte) (Supplier, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return Supplier{}, err
	}
	return *c, nil
}

func (c *Supplier) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func supplierCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c Supplier
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

	responseWithJSON(w, http.StatusOK, data, "supplier created")
	return
}

func supplierReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []Supplier
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "supplier read")
	return
}

func supplierReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Supplier
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "supplier read")
	return
}

func supplierUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Supplier
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

	responseWithJSON(w, http.StatusOK, data, "supplier updated")
	return
}

func supplierDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Supplier
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, nil, "supplier deleted")
	return
}
