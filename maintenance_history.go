package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type MaintenanceHistory struct {
	gorm.Model
	EquipmentID           uint                `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Equipment             Equipment           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ServiceProviderID     uint                `gorm:"type:int(10);index;" validate:"alphanum,len=10"`
	ServiceProvider       ServiceProvider     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID                uint                `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	User                  User                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MaintenanceScheduleID uint                `gorm:"type:int(10);index;" validate:"alphanum,len=10"`
	MaintenanceSchedule   MaintenanceSchedule `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MaintenanceDate       time.Time           `gorm:"not null" validate:"required,datetime"`
	MaintenanceTime       time.Time           `gorm:"not null" validate:"required,datetime"`
	AdditionalNotes       string              `gorm:"type:varchar(500)" validate:"max=500"`
}

func (c *MaintenanceHistory) Decode(data []byte) (MaintenanceHistory, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return MaintenanceHistory{}, err
	}
	return *c, nil
}

func (c *MaintenanceHistory) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func maintenanceHistoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c MaintenanceHistory
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

	responseWithJSON(w, http.StatusOK, data, "maintenance history created")
	return
}

func maintenanceHistoryReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []MaintenanceHistory
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance history read")
	return
}

func maintenanceHistoryReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceHistory
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, data, "maintenance history read")
	return
}

func maintenanceHistoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceHistory
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

	responseWithJSON(w, http.StatusOK, data, "maintenance history updated")
	return
}

func maintenanceHistoryDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceHistory
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}
	responseWithJSON(w, http.StatusOK, data, "maintenance history deleted")
	return
}
