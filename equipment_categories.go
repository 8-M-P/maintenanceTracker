package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type EquipmentCategory struct {
	gorm.Model
	CompanyID        uint               `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Company          Company            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentCategoryID uint               `gorm:"default:null" validate:"alphanum"`
	ParentCategory   *EquipmentCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CategoryName     string             `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	IsMainCategory   bool               `gorm:"default:false;not null" validate:"required,boolean"`
}

func (c *EquipmentCategory) Decode(data []byte) (EquipmentCategory, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return EquipmentCategory{}, err
	}
	return *c, nil
}

func (c *EquipmentCategory) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func equipmentCategoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c EquipmentCategory
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

	responseWithJSON(w, http.StatusOK, data, "equipment category created")
	return
}

func equipmentCategoryReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []EquipmentCategory
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "equipment categories retrieved")
	return
}

func equipmentCategoryReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return
	}

	var data EquipmentCategory
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "equipment category retrieved")
	return

}

func equipmentCategoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data EquipmentCategory
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

	responseWithJSON(w, http.StatusOK, data, "equipment category updated")
	return
}

func equipmentCategoryDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		responseWithMsg(w, http.StatusBadRequest, "id is required")
		return
	}

	result := db.Delete(&EquipmentCategory{}, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithMsg(w, http.StatusOK, "equipment category deleted")
	return
}
