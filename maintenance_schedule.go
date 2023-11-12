package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type MaintenanceSchedule struct {
	gorm.Model
	EquipmentID       uint            `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Equipment         Equipment       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MaintenanceTypeID uint            `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	MaintenanceType   MaintenanceType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReminderSent      bool            `gorm:"type:tinyint(1);default:0;not null" validate:"required,boolean"`
	ScheduledDate     time.Time       `gorm:"not null" validate:"required,datetime"`
	ScheduledTime     time.Time       `gorm:"not null" validate:"required,datetime"`
	Notes             sql.NullString  `gorm:"type:varchar(500)" validate:"max=500"`
}

func (c *MaintenanceSchedule) Decode(data []byte) (MaintenanceSchedule, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return MaintenanceSchedule{}, err
	}
	return *c, nil
}

func (c *MaintenanceSchedule) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func maintenanceScheduleCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c MaintenanceSchedule
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

	responseWithJSON(w, http.StatusOK, data, "maintenance schedule created")
	return
}

func maintenanceScheduleReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []MaintenanceSchedule
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance schedule read")
	return
}

func maintenanceScheduleReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceSchedule
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance schedule read")
	return
}

func maintenanceScheduleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceSchedule
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

	responseWithJSON(w, http.StatusOK, data, "maintenance schedule updated")
	return
}

func maintenanceScheduleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data MaintenanceSchedule
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "maintenance schedule deleted")
	return
}
