package filter

import (
	"fmt"
)

// Field represents an SQL field
type Field string

// Operator represents an SQL operator
type Operator string

const (
	// OperatorEqual represents an SQL operator of the same name
	OperatorEqual Operator = "eq"
	// OperatorNot represents an SQL operator of the same name
	OperatorNot Operator = "not"
	// OperatorNotEqual represents an SQL operator of the same name
	OperatorNotEqual Operator = "noteq"
	// OperatorLessThan represents an SQL operator of the same name
	OperatorLessThan Operator = "lt"
	// OperatorLessThanEqual represents an SQL operator of the same name
	OperatorLessThanEqual Operator = "lte"
	// OperatorGreaterThan represents an SQL operator of the same name
	OperatorGreaterThan Operator = "gt"
	// OperatorGreaterThanEqual represents an SQL operator of the same name
	OperatorGreaterThanEqual Operator = "gte"
	// OperatorAnd represents an SQL operator of the same name
	OperatorAnd Operator = "and"
	// OperatorOr represents an SQL operator of the same name
	OperatorOr Operator = "or"
	// OperatorLike represents an SQL operator of the same name
	OperatorLike Operator = "like"
	// OperatorIn represents an SQL operator of the same name
	OperatorIn Operator = "in"
)

// OperandMap is the map of operands to its query string equivalent
var OperandMap = map[Operator]string{
	OperatorEqual:            " = ",
	OperatorNot:              " ! ",
	OperatorNotEqual:         " != ",
	OperatorLessThan:         " < ",
	OperatorLessThanEqual:    " <= ",
	OperatorGreaterThan:      " > ",
	OperatorGreaterThanEqual: " >= ",
	OperatorAnd:              " AND ",
	OperatorOr:               " OR ",
	OperatorLike:             " LIKE ",
	OperatorIn:               " IN ",
}

// QueryPart represents part of a query
type QueryPart interface {
	ToString() string
}

// Clause represents a simple clause
type Clause struct {
	Operand1 interface{}
	Operand2 interface{}
	Operator Operator
}

// GetArgs gets the arguments required for this clause
func (c *Clause) GetArgs(args []interface{}) []interface{} {
	switch c.Operand1.(type) {
	case Field, *Field:
	case Clause:
		clause1 := c.Operand1.(Clause)
		args = clause1.GetArgs(args)
	case *Clause:
		clause1 := c.Operand1.(*Clause)
		args = clause1.GetArgs(args)
	default:
		args = append(args, c.Operand1)
	}

	switch c.Operand2.(type) {
	case Field, *Field:
	case Clause:
		clause2 := c.Operand2.(Clause)
		args = clause2.GetArgs(args)
	case *Clause:
		clause1 := c.Operand2.(*Clause)
		args = clause1.GetArgs(args)
	default:
		args = append(args, c.Operand2)
	}

	return args
}

func (c *Clause) handlePointers() error {
	switch c.Operand1.(type) {
	case *Field:
		if c.Operand1 == nil {
			return fmt.Errorf("operand1 is nil, expected non-nil *Field")
		}
		operand1 := c.Operand1.(*Field)
		c.Operand1 = *operand1
	case *Clause:
		if c.Operand1 == nil {
			return fmt.Errorf("operand1 is nil, expected non-nil *Clause")
		}
		operand1 := c.Operand1.(*Clause)
		c.Operand1 = *operand1
	default:
	}

	switch c.Operand2.(type) {
	case *Field:
		if c.Operand2 == nil {
			return fmt.Errorf("operand2 is nil, expected non-nil *Field")
		}
		operand2 := c.Operand2.(*Field)
		c.Operand2 = *operand2
	case *Clause:
		if c.Operand2 == nil {
			return fmt.Errorf("operand2 is nil, expected non-nil *Clause")
		}
		operand2 := c.Operand2.(*Clause)
		c.Operand2 = *operand2
	default:
	}

	return nil
}

// ToQueryString returns the string representation of the clause
func (c *Clause) ToQueryString() (string, error) {
	err := c.handlePointers()
	if err != nil {
		return "", err
	}

	switch c.Operand1.(type) {
	case Field:
		switch c.Operand2.(type) {
		case Field:
			return c.toStringFieldVsField()
		case Clause:
			return c.toStringFieldVsClause()
		default:
			return c.toStringFieldVsValue()
		}
	case Clause:
		switch c.Operand2.(type) {
		case Field:
			return c.toStringClauseVsField()
		case Clause:
			return c.toStringClauseVsClause()
		default:
			return "", fmt.Errorf("unsupported filter param combination: clause vs value")
		}
	default:
		switch c.Operand2.(type) {
		case Field:
			return "", fmt.Errorf("unsupported filter param combination: value vs field")
		case Clause:
			return "", fmt.Errorf("unsupported filter param combination: value vs clause")
		default:
			return "", fmt.Errorf("unsupported filter param combination: value vs value")
		}
	}
}

func (c *Clause) toStringFieldVsField() (string, error) {
	field1, ok := c.Operand1.(Field)
	if !ok {
		return "", fmt.Errorf("failed processing field [%#v]", field1)
	}
	field2, ok := c.Operand2.(Field)
	if !ok {
		return "", fmt.Errorf("failed processing field [%#v]", field2)
	}
	return fmt.Sprintf("(%s %s %s)", field1, OperandMap[c.Operator], field2), nil
}

func (c *Clause) toStringFieldVsClause() (string, error) {
	field, ok := c.Operand1.(Field)
	if !ok {
		return "", fmt.Errorf("failed processing field [%#v]", field)
	}
	clause, ok := c.Operand2.(Clause)
	if !ok {
		return "", fmt.Errorf("failed processing clause [%#v]", clause)
	}
	clauseStr, err := clause.ToQueryString()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s %s %s)", field, OperandMap[c.Operator], clauseStr), nil
}

func (c *Clause) toStringFieldVsValue() (string, error) {
	field, ok := c.Operand1.(Field)
	if !ok {
		return "", fmt.Errorf("failed processing field [%#v]", field)
	}
	if c.Operator == OperatorIn {
		return fmt.Sprintf("(%s %s (?))", field, OperandMap[c.Operator]), nil
	}
	return fmt.Sprintf("(%s %s ?)", field, OperandMap[c.Operator]), nil
}

func (c *Clause) toStringClauseVsField() (string, error) {
	clause, ok := c.Operand1.(Clause)
	if !ok {
		return "", fmt.Errorf("failed processing clause [%#v]", clause)
	}
	clauseStr, err := clause.ToQueryString()
	if err != nil {
		return "", err
	}
	field, ok := c.Operand2.(Field)
	if !ok {
		return "", fmt.Errorf("failed processing field [%#v]", field)
	}
	return fmt.Sprintf("(%s %s %s)", clauseStr, OperandMap[c.Operator], field), nil
}

func (c *Clause) toStringClauseVsClause() (string, error) {
	clause1, ok := c.Operand1.(Clause)
	if !ok {
		return "", fmt.Errorf("failed processing clause [%#v]", clause1)
	}
	clause1Str, err := clause1.ToQueryString()
	if err != nil {
		return "", err
	}
	clause2, ok := c.Operand2.(Clause)
	if !ok {
		return "", fmt.Errorf("failed processing clause [%#v]", clause2)
	}
	clause2Str, err := clause2.ToQueryString()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("(%s %s %s)", clause1Str, OperandMap[c.Operator], clause2Str), nil
}

// Pagination represents the pagination part of an SQL query
type Pagination struct {
	Page     int
	PageSize int
}

// GetOffset returns the offset required for the current page
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetArgs gets the arguments required for this pagination
func (p *Pagination) GetArgs(args []interface{}) []interface{} {
	args = append(args, p.PageSize, p.GetOffset())
	return args
}

// GetPageCount calculates the page count based on pagination settings and total item count
func (p *Pagination) GetPageCount(itemCount int) (pageCount int) {
	pageCount = itemCount / p.PageSize
	pagedItems := p.PageSize * pageCount
	if pagedItems < itemCount {
		pageCount++
	}
	return
}

// ToQueryString returns the string representation of the pagination
func (p *Pagination) ToQueryString() string {
	return fmt.Sprintf("LIMIT ? OFFSET ?")
}

// Filter represents a generic SQL filter
type Filter struct {
	TableName      string
	DeletedColumn  string
	Clause         *Clause
	IncludeDeleted bool
	Pagination     Pagination
}

// GetArgs gets the arguments required for this filter
func (f *Filter) GetArgs(withPagination bool) []interface{} {
	args := make([]interface{}, 0)

	if f.Clause != nil {
		args = f.Clause.GetArgs(args)
	}

	if withPagination {
		args = f.Pagination.GetArgs(args)
	}

	return args
}

// AddClause adds a clause to this Filter
func (f *Filter) AddClause(clause Clause, operator Operator) {
	if f.Clause == nil {
		f.Clause = &clause
	} else {
		newClause := Clause{
			Operand1: f.Clause,
			Operand2: clause,
			Operator: operator,
		}
		f.Clause = &newClause
	}
}

// ToQueryString converts the Filter to a string query
func (f *Filter) ToQueryString() (string, error) {
	var err error
	clauseStr := ""

	if f.Clause != nil {
		clauseStr, err = f.Clause.ToQueryString()
		if err != nil {
			return "", err
		}
	}

	if !f.IncludeDeleted {
		if len(f.DeletedColumn) == 0 {
			f.DeletedColumn = "deleted"
		}
		if len(clauseStr) > 0 {
			clauseStr = fmt.Sprintf("(%s) AND %s.%s IS NULL ", clauseStr, f.TableName, f.DeletedColumn)
		} else {
			clauseStr = fmt.Sprintf(" %s.%s IS NULL ", f.TableName, f.DeletedColumn)
		}
	}
	return " WHERE " + clauseStr, nil
}

// BaseFilterInput is the base type for all filter inputs
type BaseFilterInput struct {
	Keyword        *string `json:"keyword,omitempty"`
	IncludeDeleted *bool   `json:"includeDeleted,omitempty"`
	Page           *int    `json:"page,omitempty"`
	PageSize       *int    `json:"pageSize,omitempty"`
}

// GetKeywordFilter produces the filter object from a list of searchable fields
func (f *BaseFilterInput) GetKeywordFilter(fields []Field, exact bool) (clause *Clause) {
	if f.Keyword == nil {
		return nil
	}

	keyword := *f.Keyword
	operator := OperatorEqual

	if !exact {
		keyword = "%" + keyword + "%"
		operator = OperatorLike
	}

	for idx, field := range fields {
		newClause := &Clause{
			Operand1: field,
			Operand2: keyword,
			Operator: operator,
		}
		if idx == 0 {
			clause = newClause
		} else {
			clause = &Clause{
				Operand1: clause,
				Operand2: newClause,
				Operator: OperatorOr,
			}
		}
	}
	return
}

// GetIncludeDeleted returns the boolean parameter which indicates whether the filter should include deleted items
func (f *BaseFilterInput) GetIncludeDeleted() bool {
	if f.IncludeDeleted == nil {
		includeDeleted := false
		f.IncludeDeleted = &includeDeleted
	}

	return *f.IncludeDeleted
}

// GetPagination returns the pagination object from a filter input
func (f *BaseFilterInput) GetPagination() Pagination {
	page := 1
	if f.Page != nil {
		page = *f.Page
	}

	pageSize := 10
	if f.PageSize != nil {
		pageSize = *f.PageSize
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
