package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
)

type Tables int // enum

const (
	CompaniesTable Tables = iota
	ComplianceDocumentsTable
	EquipmentCategoriesTable
	EquipmentDocsTable
	EquipmentTable
	InventoryTable
	MaintenanceHistoryTable
	MaintenancePartsUsageTable
	MaintenanceScheduleTable
	MaintenanceTypesTable
	NotificationsTable
	PurchaseOrdersTable
	RolesTable
	ServiceProvidersTable
	SuppliersTable
	UsersTable
)

func (t Tables) String() string {
	return [...]string{
		"companies",
		"compliance_documents",
		"equipment_categories",
		"equipment_docs",
		"equipment",
		"inventory",
		"maintenance_history",
		"maintenance_parts_usage",
		"maintenance_schedule",
		"maintenance_types",
		"notifications",
		"purchase_orders",
		"roles",
		"service_providers",
		"suppliers",
		"users",
	}[t]
}

func (t Tables) Struct() interface{} {
	switch t {
	case CompaniesTable:
		return &Company{}
	case ComplianceDocumentsTable:
		return &ComplianceDocument{}
	case EquipmentCategoriesTable:
		return &EquipmentCategory{}
	case EquipmentDocsTable:
		return &EquipmentDoc{}
	case EquipmentTable:
		return &Equipment{}
	case InventoryTable:
		return &Inventory{}
	case MaintenanceHistoryTable:
		return &MaintenanceHistory{}
	case MaintenancePartsUsageTable:
		return &MaintenancePartsUsage{}
	case MaintenanceScheduleTable:
		return &MaintenanceSchedule{}
	case MaintenanceTypesTable:
		return &MaintenanceType{}
	case NotificationsTable:
		return &Notification{}
	case PurchaseOrdersTable:
		return &PurchaseOrder{}
	case RolesTable:
		return &Role{}
	case ServiceProvidersTable:
		return &ServiceProvider{}
	case SuppliersTable:
		return &Supplier{}
	case UsersTable:
		return &User{}
	default:
		return nil
	}
}

func (t Tables) Slice() interface{} {
	switch t {
	case CompaniesTable:
		return []Company{}
	case ComplianceDocumentsTable:
		return []ComplianceDocument{}
	case EquipmentCategoriesTable:
		return []EquipmentCategory{}
	case EquipmentDocsTable:
		return []EquipmentDoc{}
	case EquipmentTable:
		return []Equipment{}
	case InventoryTable:
		return []Inventory{}
	case MaintenanceHistoryTable:
		return []MaintenanceHistory{}
	case MaintenancePartsUsageTable:
		return []MaintenancePartsUsage{}
	case MaintenanceScheduleTable:
		return []MaintenanceSchedule{}
	case MaintenanceTypesTable:
		return []MaintenanceType{}
	case NotificationsTable:
		return []Notification{}
	case PurchaseOrdersTable:
		return []PurchaseOrder{}
	case RolesTable:
		return []Role{}
	case ServiceProvidersTable:
		return []ServiceProvider{}
	case SuppliersTable:
		return []Supplier{}
	case UsersTable:
		return []User{}
	default:
		return nil
	}
}

type ValidationError struct {
	Namespace       string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field           string `json:"field"`     // by passing alt name to ReportError like below
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}

func dbDSN() string {
	dbHost := os.Getenv("DBDEVHOST")
	dbPort := os.Getenv("DBDEVPORT")
	dbUser := os.Getenv("DBDEVUSER")
	dbPass := os.Getenv("DBDEVPASSWORD")
	dbName := os.Getenv("DBDEVDATABASE")
	return dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

var (
	db       *gorm.DB
	validate *validator.Validate
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
)

func Validate(c interface{}) ValidationError {
	err := validate.Struct(c)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return ValidationError{}
		}

		for _, err := range err.(validator.ValidationErrors) {
			return ValidationError{
				Namespace:       err.Namespace(),
				Field:           err.Field(),
				StructNamespace: err.StructNamespace(),
				StructField:     err.StructField(),
				Tag:             err.Tag(),
				ActualTag:       err.ActualTag(),
				Kind:            fmt.Sprintf("%v", err.Kind()),
				Type:            fmt.Sprintf("%v", err.Type()),
				Value:           fmt.Sprintf("%v", err.Value()),
				Param:           err.Param(),
				Message:         fmt.Sprintf("%s to be compatible with rule %s", err.StructField(), err.Tag()),
			}
		}
	}
	return ValidationError{}
}

func ValidateExcept(c interface{}, exp []string) ValidationError {
	err := validate.StructExcept(c, exp...)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return ValidationError{}
		}

		for _, err := range err.(validator.ValidationErrors) {
			return ValidationError{
				Namespace:       err.Namespace(),
				Field:           err.Field(),
				StructNamespace: err.StructNamespace(),
				StructField:     err.StructField(),
				Tag:             err.Tag(),
				ActualTag:       err.ActualTag(),
				Kind:            fmt.Sprintf("%v", err.Kind()),
				Type:            fmt.Sprintf("%v", err.Type()),
				Value:           fmt.Sprintf("%v", err.Value()),
				Param:           err.Param(),
				Message:         fmt.Sprintf("%s to be compatible with rule %s", err.StructField(), err.Tag()),
			}
		}
	}
	return ValidationError{}
}

func init() {
	var err error
	db, err = gorm.Open(mysql.Open(dbDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/companies", func(r chi.Router) {
		r.Post("/", companyCreateHandler)
		r.Get("/", companyReadHandler)
		r.Get("/{id}", companyReadOneHandler)
		r.Put("/{id}", companyUpdateHandler)
		r.Delete("/{id}", companyDeleteHandler)
	})

	r.Route("/compliance-documents", func(r chi.Router) {
		r.Post("/", complianceDocumentCreateHandler)
		r.Get("/", complianceDocumentReadHandler)
		r.Get("/{id}", complianceDocumentReadOneHandler)
		r.Put("/{id}", complianceDocumentUpdateHandler)
		r.Delete("/{id}", complianceDocumentDeleteHandler)
	})

	r.Route("/equipment-categories", func(r chi.Router) {
		r.Post("/", equipmentCategoryCreateHandler)
		r.Get("/", equipmentCategoryReadHandler)
		r.Get("/{id}", equipmentCategoryReadOneHandler)
		r.Put("/{id}", equipmentCategoryUpdateHandler)
		r.Delete("/{id}", equipmentCategoryDeleteHandler)
	})

	r.Route("/equipment-docs", func(r chi.Router) {
		r.Post("/", equipmentDocCreateHandler)
		r.Get("/", equipmentDocReadHandler)
		r.Get("/{id}", equipmentDocReadOneHandler)
		r.Put("/{id}", equipmentDocUpdateHandler)
		r.Delete("/{id}", equipmentDocDeleteHandler)
	})

	r.Route("/equipment", func(r chi.Router) {
		r.Post("/", equipmentCreateHandler)
		r.Get("/", equipmentReadHandler)
		r.Get("/{id}", equipmentReadOneHandler)
		r.Put("/{id}", equipmentUpdateHandler)
		r.Delete("/{id}", equipmentDeleteHandler)
	})

	r.Route("/inventory", func(r chi.Router) {
		r.Post("/", inventoryCreateHandler)
		r.Get("/", inventoryReadHandler)
		r.Get("/{id}", inventoryReadOneHandler)
		r.Put("/{id}", inventoryUpdateHandler)
		r.Delete("/{id}", inventoryDeleteHandler)
	})

	r.Route("/maintenance-history", func(r chi.Router) {
		r.Post("/", maintenanceHistoryCreateHandler)
		r.Get("/", maintenanceHistoryReadHandler)
		r.Get("/{id}", maintenanceHistoryReadOneHandler)
		r.Put("/{id}", maintenanceHistoryUpdateHandler)
		r.Delete("/{id}", maintenanceHistoryDeleteHandler)
	})

	r.Route("/maintenance-parts-usage", func(r chi.Router) {
		r.Post("/", maintenancePartsUsageCreateHandler)
		r.Get("/", maintenancePartsUsageReadHandler)
		r.Get("/{id}", maintenancePartsUsageReadOneHandler)
		r.Put("/{id}", maintenancePartsUsageUpdateHandler)
		r.Delete("/{id}", maintenancePartsUsageDeleteHandler)
	})

	r.Route("/maintenance-schedule", func(r chi.Router) {
		r.Post("/", maintenanceScheduleCreateHandler)
		r.Get("/", maintenanceScheduleReadHandler)
		r.Get("/{id}", maintenanceScheduleReadOneHandler)
		r.Put("/{id}", maintenanceScheduleUpdateHandler)
		r.Delete("/{id}", maintenanceScheduleDeleteHandler)
	})

	r.Route("/maintenance-types", func(r chi.Router) {
		r.Post("/", maintenanceTypeCreateHandler)
		r.Get("/", maintenanceTypeReadHandler)
		r.Get("/{id}", maintenanceTypeReadOneHandler)
		r.Put("/{id}", maintenanceTypeUpdateHandler)
		r.Delete("/{id}", maintenanceTypeDeleteHandler)
	})

	r.Route("/notifications", func(r chi.Router) {
		r.Post("/", notificationCreateHandler)
		r.Get("/", notificationReadHandler)
		r.Get("/{id}", notificationReadOneHandler)
		r.Put("/{id}", notificationUpdateHandler)
		r.Delete("/{id}", notificationDeleteHandler)
	})

	r.Route("/purchase-orders", func(r chi.Router) {
		r.Post("/", purchaseOrderCreateHandler)
		r.Get("/", purchaseOrderReadHandler)
		r.Get("/{id}", purchaseOrderReadOneHandler)
		r.Put("/{id}", purchaseOrderUpdateHandler)
		r.Delete("/{id}", purchaseOrderDeleteHandler)
	})

	r.Route("/roles", func(r chi.Router) {
		r.Post("/", roleCreateHandler)
		r.Get("/", roleReadHandler)
		r.Get("/{id}", roleReadOneHandler)
		r.Put("/{id}", roleUpdateHandler)
		r.Delete("/{id}", roleDeleteHandler)
	})

	r.Route("/service-providers", func(r chi.Router) {
		r.Post("/", serviceProviderCreateHandler)
		r.Get("/", serviceProviderReadHandler)
		r.Get("/{id}", serviceProviderReadOneHandler)
		r.Put("/{id}", serviceProviderUpdateHandler)
		r.Delete("/{id}", serviceProviderDeleteHandler)
	})

	r.Route("/suppliers", func(r chi.Router) {
		r.Post("/", supplierCreateHandler)
		r.Get("/", supplierReadHandler)
		r.Get("/{id}", supplierReadOneHandler)
		r.Put("/{id}", supplierUpdateHandler)
		r.Delete("/{id}", supplierDeleteHandler)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userCreateHandler)
		r.Get("/", userReadHandler)
		r.Get("/{id}", userReadOneHandler)
		r.Put("/{id}", userUpdateHandler)
		r.Delete("/{id}", userDeleteHandler)
	})

	err := http.ListenAndServe(":8181", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
