package main

import (
	"fmt"
	_ "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
)

type Company struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null;" validate:"required,max=255"`
	Address string `gorm:"type:varchar(500);" validate:"max=500"`
	Email   string `gorm:"type:varchar(255);unique;not null" validate:"required,max=255,email"`
	Phone   string `gorm:"type:varchar(255);unique;not null" validate:"required,max=255,e164"`
}

func companyCreateHandler(w http.ResponseWriter, r *http.Request) {
	if Create(w, r, CompaniesTable) {
		fmt.Println("company created")
		return
	}
	fmt.Println("company not created")
	return
}

func companyReadHandler(w http.ResponseWriter, r *http.Request) {
	if Read(w, r, CompaniesTable) {
		fmt.Println("companies read")
		return
	}
	fmt.Println("companies not read")
	return
}

func companyReadOneHandler(w http.ResponseWriter, r *http.Request) {
	if ReadOne(w, r, CompaniesTable) {
		fmt.Println("company read")
		return
	}
	fmt.Println("company not read")
	return
}

func companyUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if Update(w, r, CompaniesTable) {
		fmt.Println("company updated")
		return
	}
	fmt.Println("company not updated")
	return
}

func companyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if Delete(w, r, CompaniesTable) {
		fmt.Println("company deleted")
		return
	}
	fmt.Println("company not deleted")
	return
}
