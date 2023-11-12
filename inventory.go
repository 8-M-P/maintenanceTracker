package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Inventory struct {
	gorm.Model
	CompanyID           uint      `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Company             Company   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name                string    `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	CurrentStock        uint      `gorm:"type:int(10);default:0" validate:"alphanum,len=10"`
	MinRequiredQuantity uint      `gorm:"type:int(10);default:0" validate:"alphanum,len=10"`
	LastOrderDate       time.Time `gorm:"not null" validate:"required,datetime"`
	Tags                string    `gorm:"type:varchar(500)" validate:"max=500"`
	Location            string    `gorm:"type:varchar(255)" validate:"max=255"`
}

func (c *Inventory) Decode(data []byte) (Inventory, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return Inventory{}, err
	}
	return *c, nil
}

func (c *Inventory) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func inventoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c Inventory
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

	responseWithJSON(w, http.StatusOK, data, "inventory created")
	return
}

func inventoryReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []Inventory
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "inventory read")
	return
}

func inventoryReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Inventory
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, data, "inventory read")
	return
}

func inventoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Inventory
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

	responseWithJSON(w, http.StatusOK, data, "inventory updated")
	return
}

func inventoryDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Inventory
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	responseWithMsg(w, http.StatusOK, "inventory deleted")
	return
}
