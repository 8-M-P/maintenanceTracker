package main

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type Role struct {
	gorm.Model
	CompanyID            uint    `gorm:"type:int(10) unsigned;not null;default:0;index:idx_company_id;column:company_id" validate:"required,alphanum,len=10"`
	Company              Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentRoleID         *uint   `gorm:"type:int(10) unsigned;default:NULL;column:parent_role_id" validate:"alphanum,len=10"`
	ParentRole           *Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoleOrDepartmentName string  `gorm:"type:varchar(255);not null" validate:"required,max=255"`
	IsDepartment         bool    `gorm:"type:tinyint(1);not null;default:0" validate:"required,boolean"`
}

func (c *Role) Decode(data []byte) (Role, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return Role{}, err
	}
	return *c, nil
}

func (c *Role) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func roleCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := Reader(r)
	if err != nil {
		responseWithMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var c Role
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

	responseWithJSON(w, http.StatusOK, data, "role created")
}

func roleReadHandler(w http.ResponseWriter, r *http.Request) {
	var data []Role
	result := db.Find(&data)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "roles retrieved")
}

func roleReadOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Role
	result := db.First(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "")
}

func roleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Role
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

	responseWithJSON(w, http.StatusOK, data, "role updated")
}

func roleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var data Role
	result := db.Delete(&data, id)
	if result.Error != nil {
		responseWithMsg(w, http.StatusBadRequest, result.Error.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, data, "role deleted")
}
