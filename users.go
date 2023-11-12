package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type User struct {
	gorm.Model
	CompanyID    uint    `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Company      Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoleID       uint    `gorm:"type:int(10);index;not null" validate:"required,alphanum,len=10"`
	Role         Role    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Username     string  `gorm:"type:varchar(50);unique;not null" validate:"required,max=50"`
	PasswordHash string  `gorm:"type:varchar(255);not null" validate:"required,max=255,sha256"`
	Email        string  `gorm:"type:varchar(255);unique;not null" validate:"required,max=255,email"`
	FirstName    string  `gorm:"type:varchar(50)" validate:"max=50"`
	LastName     string  `gorm:"type:varchar(50)" validate:"max=50"`
	Phone        string  `gorm:"type:varchar(50);unique" validate:"max=50,e164"`
}

func (c *User) Decode(data []byte) (User, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return User{}, err
	}
	return *c, nil
}

func (c *User) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c User
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

	responseWithJSON(w, http.StatusOK, data, "user created")
}

func userReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []User
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "")
}

func userReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data User
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "")
}

func userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data User
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

	responseWithJSON(w, http.StatusOK, data, "user updated")
}

func userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data User
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithMsg(w, http.StatusOK, "user deleted")
}
