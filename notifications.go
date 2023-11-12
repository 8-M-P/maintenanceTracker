package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type Notification struct {
	gorm.Model
	UserID           uint    `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	User             User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RelatedID        uint    `gorm:"type:int(10)" validate:"alphanum,len=10"`
	RelatedType      string  `gorm:"type:ENUM('inventory','equipments','schedule','role','providers','parts_usage','documents');not null;default:'inventory';column:related_type" validate:"required,oneof=inventory equipments schedule role providers parts_usage documents"`
	NotificationType string  `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	Message          *string `gorm:"type:text;not null" validate:"required,max=65535"`
	Status           string  `gorm:"type:ENUM('Unread','Read','Dismissed');default:'Unread';column:status" validate:"oneof=Unread Read Dismissed"`
}

func (c *Notification) Decode(data []byte) (Notification, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return Notification{}, err
	}
	return *c, nil
}

func (c *Notification) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func notificationCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c Notification
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

	responseWithJSON(w, http.StatusOK, data, "notification created")
	return
}

func notificationReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []Notification
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "notification read")
	return
}

func notificationReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Notification
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "notification read")
	return
}

func notificationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Notification
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

	responseWithJSON(w, http.StatusOK, data, "notification updated")
	return
}

func notificationDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Notification
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, nil, "notification deleted")
	return
}
