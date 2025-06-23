package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/nuuid"
)

type VehicleStatus string
type VehicleType string

const (
	// VehicleStatusInUse indicates a vehicle that is actively in use
	VehicleStatusInUse VehicleStatus = "in_use"
	// VehicleStatusRetired indicates a vehicle that is no longer in use
	VehicleStatusRetired VehicleStatus = "retired"
	// VehicleStatusSold indicates a vehicle that has been sold
	VehicleStatusSold VehicleStatus = "sold"
)

const (
	// VehicleTypeCar indicates a vehicle of type Car
	VehicleTypeCar VehicleType = "car"
	// VehicleTypeTruck indicates a vehicle of type Truck
	VehicleTypeTruck VehicleType = "truck"
	// VehicleTypeBicycle indicates a vehicle of type Bicycle
	VehicleTypeBicycle VehicleType = "bicycle"
	// VehicleTypeOther indicates a vehicle of type Other
	VehicleTypeOther VehicleType = "other"
)

const (
	// VehicleColumnID represents the corresponding column in Vehicle table
	VehicleColumnID filter.Field = "vehicles.entity_id"
	// VehicleColumnName represents the corresponding column in Vehicle table
	VehicleColumnName filter.Field = "vehicles.name"
	// VehicleColumnMake represents the corresponding column in Vehicle table
	VehicleColumnMake filter.Field = "vehicles.make"
	// VehicleColumnModel represents the corresponding column in Vehicle table
	VehicleColumnModel filter.Field = "vehicles.model"
	// VehicleColumnYear represents the corresponding column in Vehicle table
	VehicleColumnYear filter.Field = "vehciles.year"
	// VehicleColumnType represents the corresponding column in Vehicle table
	VehicleColumnType filter.Field = "vehicles.type"
	// VehicleColumnTitleHolder represents the corresponding column in Vehicle table
	VehicleColumnTitleHolder filter.Field = "vehicles.title_holder"
	// VehicleColumnLicensePlateNumber represents the corresponding column in Vehicle table
	VehicleColumnLicensePlateNumber filter.Field = "vehicles.license_plate_number"
	// VehicleColumnPurchaseDate represents the corresponding column in Vehicle table
	VehicleColumnPurchaseDate filter.Field = "vehicles.purchase_date"
	// VehicleColumnInitialValue represents the corresponding column in Vehicle table
	VehicleColumnInitialValue filter.Field = "vehicles.initial_value"
	// VehicleColumnInitialValueDate represents the corresponding column in Vehicle table
	VehicleColumnInitialValueDate filter.Field = "vehicles.initial_value_date"
	// VehicleColumnCurrentValue represents the corresponding column in Vehicle table
	VehicleColumnCurrentValue filter.Field = "vehicles.current_value"
	// VehicleColumnCurrentvalueDate represents the corresponding column in Vehicle table
	VehicleColumnCurrentvalueDate filter.Field = "vehicles.current_value_date"
	// VehicleColumnAnnualDepreciationPercent represents the corresponding column in Vehicle table
	VehicleColumnAnnualDepreciationPercent filter.Field = "vehicles.annual_depreciation_percent"
	// VehicleColumnStatus represents the corresponding column in Vehicle table
	VehicleColumnStatus filter.Field = "vehicles.status"
	// VehicleColumnCreated represents the corresponding column in Vehicle table
	VehicleColumnCreated filter.Field = "vehicles.created"
	// VehicleColumnCreatedBy represents the corresponding column in Vehicle table
	VehicleColumnCreatedBy filter.Field = "vehicles.created_by"
	// VehicleColumnUpdated represents the corresponding column in Vehicle table
	VehicleColumnUpdated filter.Field = "vehicles.updated"
	// VehicleColumnUpdatedBy represents the corresponding column in Vehicle table
	VehicleColumnUpdatedBy filter.Field = "vehicles.updated_by"
	// VehicleColumnDeleted represents the corresponding column in Vehicle table
	VehicleColumnDeleted filter.Field = "vehicles.deleted"
	// VehicleColumnDeletedBy represents the corresponding column in Vehicle table
	VehicleColumnDeletedBy filter.Field = "vehicles.deleted_by"
)

const (
	// VehicleValueColumnID represents the corresponding column in the Vehicle Value table
	VehicleValueColumnID filter.Field = "vehicle_values.entity_id"
	// VehicleValueColumnDate represents the corresponding column in the Vehicle Value table
	VehicleValueColumnDate filter.Field = "vehicle_values.date"
	// VehicleValueColumnValue represents the corresponding column in the Vehicle Value table
	VehicleValueColumnValue filter.Field = "vehicle_values.value"
	// VehicleValueColumnCreated represents the corresponding column in the Vehicle Value table
	VehicleValueColumnCreated filter.Field = "vehicle_values.created"
	// VehicleValueColumnCreatedBy represents the corresponding column in the Vehicle Value table
	VehicleValueColumnCreatedBy filter.Field = "vehicle_values.created_by"
	// VehicleValueColumnUpdated represents the corresponding column in the Vehicle Value table
	VehicleValueColumnUpdated filter.Field = "vehicle_values.updated"
	// VehicleValueColumnUpdatedBy represents the corresponding column in the Vehicle Value table
	VehicleValueColumnUpdatedBy filter.Field = "vehicle_values.updated_by"
	// VehicleValueColumnDeleted represents the corresponding column in the Vehicle Value table
	VehicleValueColumnDeleted filter.Field = "vehicle_values.deleted"
	// VehicleValueColumnDeletedBy represents the corresponding column in the Vehicle Value table
	VehicleValueColumnDeletedBy filter.Field = "vehicle_values.deleted_by"
)

// Vehicle represents a Vehicle object
type Vehicle struct {
	ID                        uuid.UUID      `db:"entity_id" validate:"min=36,max=36"`
	Name                      string         `db:"name" validate:"max=255"`
	Make                      string         `db:"make" validate:"max=255"`
	Model                     string         `db:"model" validate:"max=255"`
	Year                      int            `db:"year" validate:"min=0"`
	Type                      VehicleType    `db:"type"`
	TitleHolder               string         `db:"title_holder" validate:"max=255"`
	LicensePlateNumber        string         `db:"license_plate_number" validate:"max=255"`
	PurchaseDate              time.Time      `db:"purhcase_date"`
	InitialValue              float64        `db:"initial_value" validate:"min=0"`
	InitialValueDate          time.Time      `db:"initial_value_date"`
	CurrentValue              float64        `db:"current_value" validate:"min=0"`
	CurrentValueDate          time.Time      `db:"current_value_date"`
	AnnualDepreciationPercent float64        `db:"annual_depreciation_percent"`
	Status                    VehicleStatus  `db:"status"`
	Created                   time.Time      `db:"created"`
	CreatedBy                 uuid.UUID      `db:"created_by" validate:"min=36,max=36"`
	Updated                   null.Time      `db:"updated"`
	UpdatedBy                 nuuid.NUUID    `db:"updated_by" validate:"min=36,max=36"`
	Deleted                   null.Time      `db:"deleted"`
	DeletedBy                 nuuid.NUUID    `db:"deleted_by" validate:"min=36,max=36"`
	Values                    []VehicleValue `db:"-"`
}

// VehicleInput represents an input struct for Vehicle entity
type VehicleInput struct {
	ID                        uuid.UUID           `json:"id"`
	Name                      string              `json:"name"`
	Make                      string              `json:"make"`
	Model                     string              `json:"model"`
	Year                      int                 `json:"year"`
	Type                      VehicleType         `json:"type"`
	TitleHolder               string              `json:"titleHolder"`
	LicensePlateNumber        string              `json:"licensePlateNumber"`
	PurchaseDate              cachetime.CacheTime `json:"purchaseDate"`
	InitialValue              float64             `json:"initialValue"`
	InitialValueDate          cachetime.CacheTime `json:"initialValueDate"`
	CurrentValue              float64             `json:"currentValue"`
	CurrentValueDate          cachetime.CacheTime `json:"currentValueDate"`
	AnnualDepreciationPercent float64             `json:"annualDepreciationPercent"`
	Status                    VehicleStatus       `json:"status"`
	Values                    []VehicleValueInput `json:"values"`
}

// VehicleOutput is the JSON-compatible object representation of Vehicle
type VehicleOutput struct {
	ID                        uuid.UUID            `json:"id"`
	Name                      string               `json:"name"`
	Make                      string               `json:"make"`
	Model                     string               `json:"model"`
	Year                      int                  `json:"year"`
	Type                      VehicleType          `json:"type"`
	TitleHolder               string               `json:"titleHolder"`
	LicensePlateNumber        string               `json:"licensePlateNumber"`
	PurchaseDate              cachetime.CacheTime  `json:"purchaseDate"`
	InitialValue              float64              `json:"initialValue"`
	InitialValueDate          cachetime.CacheTime  `json:"initialValueDate"`
	CurrentValue              float64              `json:"currentValue"`
	CurrentValueDate          cachetime.CacheTime  `json:"currentValueDate"`
	AnnualDepreciationPercent float64              `json:"annualDepreciationPercent"`
	Status                    VehicleStatus        `json:"status"`
	Created                   cachetime.CacheTime  `json:"created"`
	CreatedBy                 uuid.UUID            `json:"createdBy"`
	Updated                   cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy                 nuuid.NUUID          `json:"updatedBy,omitempty"`
	Deleted                   cachetime.NCacheTime `json:"deleted,omitempty"`
	DeletedBy                 nuuid.NUUID          `json:"deletedBy,omitempty"`
	Values                    []VehicleValueOutput `json:"values"`
}

// VehicleValue represents a snapshot of a Vehicle's value at a given time
type VehicleValue struct {
	ID        uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	VehicleID uuid.UUID   `db:"vehicle_entity_id" validate:"min=36,max=36"`
	Date      time.Time   `db:"date"`
	Value     float64     `db:"value" validate:"min=0"`
	Created   time.Time   `db:"created"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"min=36,max=36"`
	Updated   null.Time   `db:"updated"`
	UpdatedBy nuuid.NUUID `db:"updated_by" validate:"min=36,max=36"`
	Deleted   null.Time   `db:"deleted"`
	DeletedBy nuuid.NUUID `db:"deleted_by" validate:"min=36,max=36"`
}

// VehicleValueInput represents an input struct for Vehicle Value entity
type VehicleValueInput struct {
	ID        uuid.UUID           `json:"id"`
	VehicleID uuid.UUID           `json:"vehicleId"`
	Date      cachetime.CacheTime `json:"date"`
	Value     float64             `json:"value"`
}

// VehicleValueOutput is the JSON-compatible object representation of Vehicle Value
type VehicleValueOutput struct {
	ID        uuid.UUID            `json:"id"`
	VehicleID uuid.UUID            `json:"vehicleId"`
	Date      cachetime.CacheTime  `json:"date"`
	Value     float64              `json:"value"`
	Created   cachetime.CacheTime  `json:"created"`
	CreatedBy uuid.UUID            `json:"createdBy"`
	Updated   cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy nuuid.NUUID          `json:"updatedBy,omitempty"`
	Deleted   cachetime.NCacheTime `json:"deleted,omitempty"`
	DeletedBy nuuid.NUUID          `json:"deletedBy,omitempty"`
}

// VehicleFilterInput is the filter input object for Vehicles
type VehicleFilterInput struct {
	filter.BaseFilterInput
}
