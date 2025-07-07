package service

import (
	"math"

	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// PropertyImpl is the service provider implementation
type PropertyImpl struct {
	Repository repository.Property `inject:"propertyRepository"`
}

// Startup performs startup functions
func (s *PropertyImpl) Startup() {
	logger.Trace("Property Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *PropertyImpl) Shutdown() {
	logger.Trace("Property Service shutting down...")
}

// Create creates a new Property
func (s *PropertyImpl) Create(input model.PropertyInput, userID uuid.UUID) (*model.Property, error) {
	property := model.NewPropertyFromInput(input, userID)
	err := s.Repository.Create(property)
	if err != nil {
		return nil, err
	}
	return &property, err
}

// GetByID fetches a Property by its ID
func (s *PropertyImpl) GetByID(id uuid.UUID, withValues bool, valueStartDate, valueEndDate cachetime.NCacheTime, pageSize *int) (*model.Property, error) {
	properties, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(properties) != 1 {
		return nil, failure.EntityNotFound("get by ID", "Property")
	}

	property := properties[0]

	if withValues {
		filter := model.PropertyValueFilterInput{
			PropertyIDs: &[]uuid.UUID{id},
		}

		if valueStartDate.Valid {
			filter.StartDate = valueStartDate
		}

		if valueEndDate.Valid {
			filter.EndDate = valueEndDate
		}

		if pageSize != nil {
			filter.PageSize = pageSize
		}

		values, _, err := s.Repository.ResolveValuesByFilter(filter.ToFilter())
		if err != nil {
			return nil, err
		}

		property.AttachValues(values, true)
	}

	return &property, nil
}

// GetByFilter fetches a set of Properties by its filter
func (s *PropertyImpl) GetByFilter(input model.PropertyFilterInput) ([]model.Property, model.PageInfoOutput, error) {
	return s.Repository.ResolveByFilter(input.ToFilter())
}

// Update updates an existing Property
func (s *PropertyImpl) Update(input model.PropertyInput, userID uuid.UUID) (*model.Property, error) {
	properties, err := s.Repository.ResolveByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return nil, err
	}

	if len(properties) != 1 {
		return nil, failure.EntityNotFound("update", "Property")
	}

	property := properties[0]

	err = property.Update(input, userID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.Update(property)
	if err != nil {
		return nil, err
	}

	return &property, err
}

// Delete deletes an existing Property. The method will find all the property's values
// and delete all of them also.
func (s *PropertyImpl) Delete(id uuid.UUID, userID uuid.UUID) (*model.Property, error) {
	properties, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(properties) != 1 {
		return nil, failure.EntityNotFound("delete", "Property")
	}

	property := properties[0]

	// pre-validate to save one database call
	if !property.Deleted.Valid && !property.DeletedBy.Valid {
		filter := model.PropertyValueFilterInput{}
		filter.PropertyIDs = &[]uuid.UUID{property.ID}

		page := 1
		pageSize := math.MaxInt

		filter.Page = &page
		filter.PageSize = &pageSize

		values, _, err := s.Repository.ResolveValuesByFilter(filter.ToFilter())
		if err != nil {
			return nil, err
		}

		property.AttachValues(values, true)
	}

	err = property.Delete(userID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.Update(property)
	if err != nil {
		return nil, err
	}

	return &property, err
}

// CreateValue creates a new Property Value
func (s *PropertyImpl) CreateValue(input model.PropertyValueInput, userID uuid.UUID) (*model.PropertyValue, error) {
	properties, err := s.Repository.ResolveByIDs([]uuid.UUID{input.PropertyID})
	if err != nil {
		return nil, err
	}

	if len(properties) != 1 {
		return nil, failure.EntityNotFound("create value", "Property Value")
	}

	property := properties[0]

	if property.Deleted.Valid || property.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("add value", "Property", "the Property is already deleted")
	}

	if property.Status == model.PropertyStatusSold {
		return nil, failure.OperationNotPermitted("add value", "Property", "the Property has been sold")
	}

	lastValues, err := s.Repository.ResolveLastValuesByPropertyID(property.ID, 1)
	if err != nil {
		return nil, err
	}

	if len(lastValues) != 1 {
		return nil, failure.EntityNotFound("create value", "Property Last Value")
	}

	lastValue := lastValues[0]
	isNewerValue := lastValue.Date.Before(input.Date.Time())
	var propertyToUpdate *model.Property

	if isNewerValue {
		property.SetCurrentValue(input, userID)
		propertyToUpdate = &property
	}

	propertyValue := model.NewPropertyValueFromInput(input, property.ID, userID)
	err = s.Repository.CreateValue(propertyValue, propertyToUpdate)
	if err != nil {
		return nil, err
	}

	return &propertyValue, nil
}

// GetValueByID fetches a Property Value by its ID
func (s *PropertyImpl) GetValueByID(id uuid.UUID) (*model.PropertyValue, error) {
	values, err := s.Repository.ResolveValuesByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(values) != 1 {
		return nil, failure.EntityNotFound("get by ID", "Property Value")
	}

	return &values[0], nil
}

// GetValuesByFilter fetches a set of Property Values by its filter
func (s *PropertyImpl) GetValuesByFilter(input model.PropertyValueFilterInput) ([]model.PropertyValue, model.PageInfoOutput, error) {
	return s.Repository.ResolveValuesByFilter(input.ToFilter())
}

// UpdateValue updates an existing Property Value
func (s *PropertyImpl) UpdateValue(input model.PropertyValueInput, userID uuid.UUID) (*model.PropertyValue, error) {
	properties, err := s.Repository.ResolveByIDs([]uuid.UUID{input.PropertyID})
	if err != nil {
		return nil, err
	}

	if len(properties) != 1 {
		return nil, failure.EntityNotFound("update", "Property Value")
	}

	property := properties[0]

	if property.Deleted.Valid || property.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("update", "Property Value", "the Property is already deleted")
	}

	if property.Status == model.PropertyStatusSold {
		return nil, failure.OperationNotPermitted("update", "Property Value", "the Property has been sold")
	}

	propertyValues, err := s.Repository.ResolveValuesByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return nil, err
	}

	if len(propertyValues) != 1 {
		return nil, failure.EntityNotFound("update", "Property Value")
	}

	propertyValue := propertyValues[0]

	if propertyValue.Deleted.Valid || propertyValue.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("update", "Property Value", "the Property Value is already deleted")
	}

	err = propertyValue.Update(input, userID)
	if err != nil {
		return nil, err
	}

	currentValues, err := s.Repository.ResolveLastValuesByPropertyID(property.ID, 1)
	if err != nil {
		return nil, err
	}

	if len(currentValues) != 1 {
		return nil, failure.EntityNotFound("update", "Property Value")
	}

	currentValue := currentValues[0]
	isNewerOrCurrentValue := currentValue.Date.Before(input.Date.Time()) || input.ID == currentValue.ID
	var propertyToUpdate *model.Property

	if isNewerOrCurrentValue {
		property.SetCurrentValue(input, userID)
		propertyToUpdate = &property
	}

	err = s.Repository.UpdateValue(propertyValue, propertyToUpdate)
	if err != nil {
		return nil, err
	}

	return &propertyValue, nil
}

// DeleteValue deletes an existing Property Value
func (s *PropertyImpl) DeleteValue(id uuid.UUID, userID uuid.UUID) (*model.PropertyValue, error) {
	propertyValues, err := s.Repository.ResolveValuesByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(propertyValues) != 1 {
		return nil, failure.EntityNotFound("delete", "Property Value")
	}

	propertyValue := propertyValues[0]

	if propertyValue.Deleted.Valid || propertyValue.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("delete", "Property Value", "the Property Value is already deleted")
	}

	properties, err := s.Repository.ResolveByIDs([]uuid.UUID{propertyValue.PropertyID})
	if err != nil {
		return nil, err
	}

	if len(properties) != 1 {
		return nil, failure.EntityNotFound("delete", "Property")
	}

	property := properties[0]

	if property.Deleted.Valid || property.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("delete", "Property Value", "the Property is already deleted")
	}

	if property.Status == model.PropertyStatusSold {
		return nil, failure.OperationNotPermitted("delete", "Property Value", "the Property has been sold")
	}

	propertyValue.Delete(userID)

	currentValues, err := s.Repository.ResolveLastValuesByPropertyID(property.ID, 2)
	if err != nil {
		return nil, err
	}

	if len(currentValues) < 1 {
		return nil, failure.EntityNotFound("delete", "Property Current Value")
	}

	if len(currentValues) < 2 {
		return nil, failure.OperationNotPermitted("delete", "Property Value", "cannot delete the only Property Value belonging to a Property")
	}

	currentValueDeleted := propertyValue.ID.String() == currentValues[0].ID.String()
	var propertyToUpdate *model.Property

	if currentValueDeleted {
		newCurrentValueInput := model.PropertyValueInput{
			ID:         currentValues[1].ID,
			PropertyID: currentValues[1].PropertyID,
			Value:      currentValues[1].Value,
			Date:       cachetime.CacheTime(currentValues[1].Date),
		}
		property.SetCurrentValue(newCurrentValueInput, userID)
		propertyToUpdate = &property
	}

	err = s.Repository.UpdateValue(propertyValue, propertyToUpdate)
	if err != nil {
		return nil, err
	}

	return &propertyValue, nil
}
