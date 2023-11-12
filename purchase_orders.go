package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type PurchaseOrder struct {
	gorm.Model
	InventoryID     uint      `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Inventory       Inventory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SupplierID      uint      `gorm:"type:int(10);index;" validate:"alphanum,len=10"`
	Supplier        Supplier  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CompanyID       uint      `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Company         Company   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID          uint      `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	User            User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	QuantityOrdered uint      `gorm:"type:int(10);not null;default:0" validate:"required,alphanum,len=10"`
	OrderDate       time.Time `gorm:"not null" validate:"required,datetime"`
	ReceivedDate    time.Time `gorm:"not null" validate:"required,datetime"`
}

func (c *PurchaseOrder) Decode(data []byte) (PurchaseOrder, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return PurchaseOrder{}, err
	}
	return *c, nil
}

func (c *PurchaseOrder) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func purchaseOrderCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c PurchaseOrder
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

	responseWithJSON(w, http.StatusOK, data, "purchase order created")
	return
}

func purchaseOrderReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []PurchaseOrder
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "")
	return
}

func purchaseOrderReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data PurchaseOrder
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "")
	return
}

func purchaseOrderUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data PurchaseOrder
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

	responseWithJSON(w, http.StatusOK, data, "")
	return
}

func purchaseOrderDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data PurchaseOrder
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithMsg(w, http.StatusOK, "purchase order deleted")
	return
}
