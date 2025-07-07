package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/nuuid"
)

type PropertyStatus string
type PropertyType string
type PropertyAreaUnit string

const (
	// PropertyStatusNotInUse indicates a property that is not actively in use
	PropertyStatusNotInUse PropertyStatus = "not_in_use"
	// PropertyStatusInUse indicates a property that is actively in use
	PropertyStatusInUse PropertyStatus = "in_use"
	// PropertyStatusRetired indicates a property that is rented out
	PropertyStatusRented PropertyStatus = "rented"
	// PropertyStatusSold indicates a property that has been sold
	PropertyStatusSold PropertyStatus = "sold"
)

const (
	// PropertyTypeLand indicates a property of type Land
	PropertyTypeLand PropertyType = "land"
	// PropertyTypeHouse indicates a property of type House
	PropertyTypeHouse PropertyType = "house"
	// PropertyTypeApartment indicates a property of type Apartment
	PropertyTypeApartment PropertyType = "apartment"
)

const (
	// PropertyAreaUnitSQFT indicates a property with areas measured in square feet
	PropertyAreaUnitSQFT PropertyAreaUnit = "sqft"
	// PropertyAreaUnitSQM indicates a property with areas measured in square meters
	PropertyAreaUnitSQM PropertyAreaUnit = "sqm"
)

const (
	// PropertyColumnID represents the corresponding column in Property table
	PropertyColumnID filter.Field = "properties.entity_id"
	// PropertyColumnName represents the corresponding column in Property table
	PropertyColumnName filter.Field = "properties.name"
	// PropertyColumnAddress represents the corresponding column in Property table
	PropertyColumnAddress filter.Field = "properties.address"
	// PropertyColumnTotalArea represents the corresponding column in Property table
	PropertyColumnTotalArea filter.Field = "properties.total_area"
	// PropertyColumnBuildingArea represents the corresponding column in Property table
	PropertyColumnBuildingArea filter.Field = "properties.building_area"
	// PropertyColumnAreaUnit represents the corresponding column in Property table
	PropertyColumnAreaUnit filter.Field = "properties.area_unit"
	// PropertyColumnType represents the corresponding column in Property table
	PropertyColumnType filter.Field = "properties.type"
	// PropertyColumnTitleHolder represents the corresponding column in Property table
	PropertyColumnTitleHolder filter.Field = "properties.title_holder"
	// PropertyColumnTaxIdentifier represents the corresponding column in Property table
	PropertyColumnTaxIdentifier filter.Field = "properties.tax_identifier"
	// PropertyColumnPurchaseDate represents the corresponding column in Property table
	PropertyColumnPurchaseDate filter.Field = "properties.purchase_date"
	// PropertyColumnInitialValue represents the corresponding column in Property table
	PropertyColumnInitialValue filter.Field = "properties.initial_value"
	// PropertyColumnInitialValueDate represents the corresponding column in Property table
	PropertyColumnInitialValueDate filter.Field = "properties.initial_value_date"
	// PropertyColumnCurrentValue represents the corresponding column in Property table
	PropertyColumnCurrentValue filter.Field = "properties.current_value"
	// PropertyColumnCurrentvalueDate represents the corresponding column in Property table
	PropertyColumnCurrentvalueDate filter.Field = "properties.current_value_date"
	// PropertyColumnAnnualAppreciationPercent represents the corresponding column in Property table
	PropertyColumnAnnualAppreciationPercent filter.Field = "properties.annual_appreciation_percent"
	// PropertyColumnStatus represents the corresponding column in Property table
	PropertyColumnStatus filter.Field = "properties.status"
	// PropertyColumnCreated represents the corresponding column in Property table
	PropertyColumnCreated filter.Field = "properties.created"
	// PropertyColumnCreatedBy represents the corresponding column in Property table
	PropertyColumnCreatedBy filter.Field = "properties.created_by"
	// PropertyColumnUpdated represents the corresponding column in Property table
	PropertyColumnUpdated filter.Field = "properties.updated"
	// PropertyColumnUpdatedBy represents the corresponding column in Property table
	PropertyColumnUpdatedBy filter.Field = "properties.updated_by"
	// PropertyColumnDeleted represents the corresponding column in Property table
	PropertyColumnDeleted filter.Field = "properties.deleted"
	// PropertyColumnDeletedBy represents the corresponding column in Property table
	PropertyColumnDeletedBy filter.Field = "properties.deleted_by"
)

const (
	// PropertyValueColumnID represents the corresponding column in the Property Value table
	PropertyValueColumnID filter.Field = "property_values.entity_id"
	// PropertyValueColumnPropertyID represents the corresponding column in the Property Value table
	PropertyValueColumnPropertyID filter.Field = "property_values.property_entity_id"
	// PropertyValueColumnDate represents the corresponding column in the Property Value table
	PropertyValueColumnDate filter.Field = "property_values.date"
	// PropertyValueColumnValue represents the corresponding column in the Property Value table
	PropertyValueColumnValue filter.Field = "property_values.value"
	// PropertyValueColumnCreated represents the corresponding column in the Property Value table
	PropertyValueColumnCreated filter.Field = "property_values.created"
	// PropertyValueColumnCreatedBy represents the corresponding column in the Property Value table
	PropertyValueColumnCreatedBy filter.Field = "property_values.created_by"
	// PropertyValueColumnUpdated represents the corresponding column in the Property Value table
	PropertyValueColumnUpdated filter.Field = "property_values.updated"
	// PropertyValueColumnUpdatedBy represents the corresponding column in the Property Value table
	PropertyValueColumnUpdatedBy filter.Field = "property_values.updated_by"
	// PropertyValueColumnDeleted represents the corresponding column in the Property Value table
	PropertyValueColumnDeleted filter.Field = "property_values.deleted"
	// PropertyValueColumnDeletedBy represents the corresponding column in the Property Value table
	PropertyValueColumnDeletedBy filter.Field = "property_values.deleted_by"
)

// Property represents a Property object
type Property struct {
	ID                        uuid.UUID        `db:"entity_id" validate:"min=36,max=36"`
	Name                      string           `db:"name" validate:"max=255"`
	Address                   string           `db:"address" validate:"max=255"`
	TotalArea                 float64          `db:"total_area" validate:"max=255"`
	BuildingArea              float64          `db:"building_area" validate:"min=0"`
	AreaUnit                  PropertyAreaUnit `db:"area_unit"`
	Type                      PropertyType     `db:"type"`
	TitleHolder               string           `db:"title_holder" validate:"max=255"`
	TaxIdentifier             string           `db:"tax_identifier" validate:"max=255"`
	PurchaseDate              time.Time        `db:"purchase_date"`
	InitialValue              float64          `db:"initial_value" validate:"min=0"`
	InitialValueDate          time.Time        `db:"initial_value_date"`
	CurrentValue              float64          `db:"current_value" validate:"min=0"`
	CurrentValueDate          time.Time        `db:"current_value_date"`
	AnnualAppreciationPercent float64          `db:"annual_appreciation_percent"`
	Status                    PropertyStatus   `db:"status"`
	Created                   time.Time        `db:"created"`
	CreatedBy                 uuid.UUID        `db:"created_by" validate:"min=36,max=36"`
	Updated                   null.Time        `db:"updated"`
	UpdatedBy                 nuuid.NUUID      `db:"updated_by" validate:"min=36,max=36"`
	Deleted                   null.Time        `db:"deleted"`
	DeletedBy                 nuuid.NUUID      `db:"deleted_by" validate:"min=36,max=36"`
	Values                    []PropertyValue  `db:"-"`
}

// NewPropertyFromInput creates a new Property from its input object
func NewPropertyFromInput(input PropertyInput, userID uuid.UUID) (p Property) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	p = Property{
		ID:                        newUUID,
		Name:                      input.Name,
		Address:                   input.Address,
		TotalArea:                 input.TotalArea,
		BuildingArea:              input.BuildingArea,
		AreaUnit:                  input.AreaUnit,
		Type:                      input.Type,
		TitleHolder:               input.TitleHolder,
		TaxIdentifier:             input.TaxIdentifier,
		PurchaseDate:              input.PurchaseDate.Time(),
		InitialValue:              input.InitialValue,
		InitialValueDate:          input.InitialValueDate.Time(),
		CurrentValue:              input.CurrentValue,
		CurrentValueDate:          input.CurrentValueDate.Time(),
		AnnualAppreciationPercent: input.AnnualAppreciationPercent,
		Status:                    input.Status,
		Created:                   now,
		CreatedBy:                 userID,
	}

	values := make([]PropertyValue, 0)

	// add initial value as first value
	values = append(values, NewPropertyValueFromInput(PropertyValueInput{
		Date:  input.InitialValueDate,
		Value: input.InitialValue,
	}, p.ID, userID))

	// add current value as second value if necessary
	if input.CurrentValue != input.InitialValue || !input.CurrentValueDate.Time().Equal(input.InitialValueDate.Time()) {
		values = append(values, NewPropertyValueFromInput(PropertyValueInput{
			Date:  input.CurrentValueDate,
			Value: input.CurrentValue,
		}, p.ID, userID))
	}

	p.Values = values

	// TODO: validate?

	return
}

// AttachValues attaches Property Values to a Property
func (p *Property) AttachValues(values []PropertyValue, clearBeforeAttach bool) {
	if clearBeforeAttach {
		p.Values = []PropertyValue{}
	}

	for _, value := range values {
		if value.PropertyID == p.ID {
			p.Values = append(p.Values, value)
		}
	}
}

// Update performs an update on a Property
func (p *Property) Update(input PropertyInput, userID uuid.UUID) error {
	if p.Deleted.Valid || p.DeletedBy.Valid {
		return failure.OperationNotPermitted("update", "Property", "already deleted")
	}

	now := time.Now()

	p.Name = input.Name
	p.Address = input.Address
	p.TotalArea = input.TotalArea
	p.BuildingArea = input.BuildingArea
	p.AreaUnit = input.AreaUnit
	p.Type = input.Type
	p.TitleHolder = input.TitleHolder
	p.TaxIdentifier = input.TaxIdentifier
	p.PurchaseDate = input.PurchaseDate.Time()
	p.InitialValue = input.InitialValue
	p.InitialValueDate = input.InitialValueDate.Time()
	p.AnnualAppreciationPercent = input.AnnualAppreciationPercent
	p.Status = input.Status
	p.Updated = null.TimeFrom(now)
	p.UpdatedBy = nuuid.From(userID)

	return nil
}

// Delete performs a delete on a Property
func (p *Property) Delete(userID uuid.UUID) error {
	if p.Deleted.Valid || p.DeletedBy.Valid {
		return failure.OperationNotPermitted("delete", "Property", "already deleted")
	}

	now := time.Now()

	p.Deleted = null.TimeFrom(now)
	p.DeletedBy = nuuid.From(userID)

	deletedValues := make([]PropertyValue, 0)
	for _, value := range p.Values {
		err := value.Delete(userID)
		if err != nil {
			return err
		}

		deletedValues = append(deletedValues, value)
	}

	p.Values = deletedValues

	return nil
}

// SetCurrentValue sets a new current value and current value date on a Property
func (p *Property) SetCurrentValue(input PropertyValueInput, userID uuid.UUID) error {
	now := time.Now()

	p.CurrentValue = input.Value
	p.CurrentValueDate = input.Date.Time()
	p.Updated = null.TimeFrom(now)
	p.UpdatedBy = nuuid.From(userID)

	// TODO: Validate ?

	return nil
}

// ToOutput converts a Property to its JSON-compatible object representation
func (p *Property) ToOutput() PropertyOutput {
	o := PropertyOutput{
		ID:                        p.ID,
		Name:                      p.Name,
		Address:                   p.Address,
		TotalArea:                 p.TotalArea,
		BuildingArea:              p.BuildingArea,
		AreaUnit:                  p.AreaUnit,
		Type:                      p.Type,
		TitleHolder:               p.TitleHolder,
		TaxIdentifier:             p.TaxIdentifier,
		PurchaseDate:              cachetime.CacheTime(p.PurchaseDate),
		InitialValue:              p.InitialValue,
		InitialValueDate:          cachetime.CacheTime(p.InitialValueDate),
		CurrentValue:              p.CurrentValue,
		CurrentValueDate:          cachetime.CacheTime(p.CurrentValueDate),
		AnnualAppreciationPercent: p.AnnualAppreciationPercent,
		Status:                    p.Status,
		Created:                   cachetime.CacheTime(p.Created),
		CreatedBy:                 p.CreatedBy,
		Updated:                   cachetime.NCacheTime(p.Updated),
		UpdatedBy:                 p.UpdatedBy,
		Deleted:                   cachetime.NCacheTime(p.Deleted),
		DeletedBy:                 p.DeletedBy,
	}

	vvOutput := make([]PropertyValueOutput, 0)
	for _, pv := range p.Values {
		vvOutput = append(vvOutput, pv.ToOutput())
	}

	o.Values = vvOutput

	return o
}

// PropertyInput represents an input struct for Property entity
type PropertyInput struct {
	ID                        uuid.UUID           `json:"id"`
	Name                      string              `json:"name"`
	Address                   string              `json:"address"`
	TotalArea                 float64             `json:"totalArea"`
	BuildingArea              float64             `json:"buildingArea"`
	AreaUnit                  PropertyAreaUnit    `json:"areaUnit"`
	Type                      PropertyType        `json:"type"`
	TitleHolder               string              `json:"titleHolder"`
	TaxIdentifier             string              `json:"taxIdentifier"`
	PurchaseDate              cachetime.CacheTime `json:"purchaseDate"`
	InitialValue              float64             `json:"initialValue"`
	InitialValueDate          cachetime.CacheTime `json:"initialValueDate"`
	CurrentValue              float64             `json:"currentValue"`
	CurrentValueDate          cachetime.CacheTime `json:"currentValueDate"`
	AnnualAppreciationPercent float64             `json:"annualDepreciationPercent"`
	Status                    PropertyStatus      `json:"status"`
}

// PropertyOutput is the JSON-compatible object representation of Property
type PropertyOutput struct {
	ID                        uuid.UUID             `json:"id"`
	Name                      string                `json:"name"`
	Address                   string                `json:"address"`
	TotalArea                 float64               `json:"totalArea"`
	BuildingArea              float64               `json:"buildingArea"`
	AreaUnit                  PropertyAreaUnit      `json:"areaUnit"`
	Type                      PropertyType          `json:"type"`
	TitleHolder               string                `json:"titleHolder"`
	TaxIdentifier             string                `json:"taxIdentifier"`
	PurchaseDate              cachetime.CacheTime   `json:"purchaseDate"`
	InitialValue              float64               `json:"initialValue"`
	InitialValueDate          cachetime.CacheTime   `json:"initialValueDate"`
	CurrentValue              float64               `json:"currentValue"`
	CurrentValueDate          cachetime.CacheTime   `json:"currentValueDate"`
	AnnualAppreciationPercent float64               `json:"annualAppreciationPercent"`
	Status                    PropertyStatus        `json:"status"`
	Created                   cachetime.CacheTime   `json:"created"`
	CreatedBy                 uuid.UUID             `json:"createdBy"`
	Updated                   cachetime.NCacheTime  `json:"updated,omitempty"`
	UpdatedBy                 nuuid.NUUID           `json:"updatedBy,omitempty"`
	Deleted                   cachetime.NCacheTime  `json:"deleted,omitempty"`
	DeletedBy                 nuuid.NUUID           `json:"deletedBy,omitempty"`
	Values                    []PropertyValueOutput `json:"values"`
}

// PropertyValue represents a snapshot of a Property's value at a given time
type PropertyValue struct {
	ID         uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	PropertyID uuid.UUID   `db:"property_entity_id" validate:"min=36,max=36"`
	Date       time.Time   `db:"date"`
	Value      float64     `db:"value" validate:"min=0"`
	Created    time.Time   `db:"created"`
	CreatedBy  uuid.UUID   `db:"created_by" validate:"min=36,max=36"`
	Updated    null.Time   `db:"updated"`
	UpdatedBy  nuuid.NUUID `db:"updated_by" validate:"min=36,max=36"`
	Deleted    null.Time   `db:"deleted"`
	DeletedBy  nuuid.NUUID `db:"deleted_by" validate:"min=36,max=36"`
}

func NewPropertyValueFromInput(input PropertyValueInput, propertyID uuid.UUID, userID uuid.UUID) (pv PropertyValue) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	pv = PropertyValue{
		ID:         newUUID,
		PropertyID: propertyID,
		Date:       input.Date.Time(),
		Value:      input.Value,
		Created:    now,
		CreatedBy:  userID,
	}

	// TODO: validate:

	return
}

// Update performs an update on a Property Value
func (pv *PropertyValue) Update(input PropertyValueInput, userID uuid.UUID) error {
	now := time.Now()

	pv.Date = input.Date.Time()
	pv.Value = input.Value
	pv.Updated = null.TimeFrom(now)
	pv.UpdatedBy = nuuid.From(userID)

	// TODO: Validate ?

	return nil
}

// Delete performs a delete on a Property Value
func (pv *PropertyValue) Delete(userID uuid.UUID) error {
	if pv.Deleted.Valid || pv.DeletedBy.Valid {
		return failure.OperationNotPermitted("delete", "Property Value", "already deleted")
	}

	now := time.Now()

	pv.Deleted = null.TimeFrom(now)
	pv.DeletedBy = nuuid.From(userID)

	return nil
}

// ToOutput converts a Property Value to its JSON-compatible object representation
func (pv *PropertyValue) ToOutput() PropertyValueOutput {
	return PropertyValueOutput{
		ID:         pv.ID,
		PropertyID: pv.PropertyID,
		Date:       cachetime.CacheTime(pv.Date),
		Value:      pv.Value,
		Created:    cachetime.CacheTime(pv.Created),
		CreatedBy:  pv.CreatedBy,
		Updated:    cachetime.NCacheTime(pv.Updated),
		UpdatedBy:  pv.UpdatedBy,
		Deleted:    cachetime.NCacheTime(pv.Deleted),
		DeletedBy:  pv.DeletedBy,
	}
}

// PropertyValueInput represents an input struct for Property Value entity
type PropertyValueInput struct {
	ID         uuid.UUID           `json:"id"`
	PropertyID uuid.UUID           `json:"propertyId"`
	Date       cachetime.CacheTime `json:"date"`
	Value      float64             `json:"value"`
}

// PropertyValueOutput is the JSON-compatible object representation of Property Value
type PropertyValueOutput struct {
	ID         uuid.UUID            `json:"id"`
	PropertyID uuid.UUID            `json:"propertyId"`
	Date       cachetime.CacheTime  `json:"date"`
	Value      float64              `json:"value"`
	Created    cachetime.CacheTime  `json:"created"`
	CreatedBy  uuid.UUID            `json:"createdBy"`
	Updated    cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy  nuuid.NUUID          `json:"updatedBy,omitempty"`
	Deleted    cachetime.NCacheTime `json:"deleted,omitempty"`
	DeletedBy  nuuid.NUUID          `json:"deletedBy,omitempty"`
}

// PropertyFilterInput is the filter input object for Propertys
type PropertyFilterInput struct {
	filter.BaseFilterInput
}

// ToFilter converts this entity-specific filter to a generic filter.Filter object
func (f *PropertyFilterInput) ToFilter() filter.Filter {
	keywordFields := []filter.Field{
		PropertyColumnName,
		PropertyColumnAddress,
		PropertyColumnType,
		PropertyColumnTitleHolder,
		PropertyColumnTaxIdentifier,
	}

	return filter.Filter{
		TableName:      "properties",
		Clause:         f.BaseFilterInput.GetKeywordFilter(keywordFields, false),
		IncludeDeleted: f.GetIncludeDeleted(),
		Pagination:     f.BaseFilterInput.GetPagination(),
	}
}

type PropertyValueFilterInput struct {
	filter.BaseFilterInput
	PropertyIDs *[]uuid.UUID         `json:"propertyIDs,omitempty"`
	StartDate   cachetime.NCacheTime `json:"startDate,omitempty"`
	EndDate     cachetime.NCacheTime `json:"endDate,omitempty"`
	ValueMin    *float64             `json:"valueMin,omitempty"`
	ValueMax    *float64             `json:"valueMax,omitempty"`
}

// ToFilter converts this entity-specific filter to a generic filter.Filter object
func (f *PropertyValueFilterInput) ToFilter() filter.Filter {
	theFilter := filter.Filter{
		TableName:      "property_values",
		IncludeDeleted: f.GetIncludeDeleted(),
		Pagination:     f.BaseFilterInput.GetPagination(),
	}

	if f.PropertyIDs != nil {
		if len(*f.PropertyIDs) > 0 {
			theFilter.AddClause(filter.Clause{
				Operand1: PropertyValueColumnPropertyID,
				Operand2: *f.PropertyIDs,
				Operator: filter.OperatorIn,
			}, filter.OperatorAnd)
		}
	}

	if f.StartDate.Valid {
		theFilter.AddClause(filter.Clause{
			Operand1: PropertyValueColumnDate,
			Operand2: f.StartDate.Time,
			Operator: filter.OperatorGreaterThanEqual,
		}, filter.OperatorAnd)
	}

	if f.EndDate.Valid {
		theFilter.AddClause(filter.Clause{
			Operand1: PropertyValueColumnDate,
			Operand2: f.EndDate.Time,
			Operator: filter.OperatorLessThanEqual,
		}, filter.OperatorAnd)
	}

	if f.ValueMin != nil {
		theFilter.AddClause(filter.Clause{
			Operand1: PropertyValueColumnValue,
			Operand2: *f.ValueMin,
			Operator: filter.OperatorGreaterThanEqual,
		}, filter.OperatorAnd)
	}

	if f.ValueMax != nil {
		theFilter.AddClause(filter.Clause{
			Operand1: PropertyValueColumnValue,
			Operand2: *f.ValueMax,
			Operator: filter.OperatorLessThanEqual,
		}, filter.OperatorAnd)
	}

	return theFilter
}
