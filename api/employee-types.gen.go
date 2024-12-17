// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for DayOffRecordDayOffType.
const (
	Bereavement   DayOffRecordDayOffType = "bereavement"
	PTO           DayOffRecordDayOffType = "PTO"
	ParentalLeave DayOffRecordDayOffType = "parental leave"
	SickLeave     DayOffRecordDayOffType = "sick leave"
)

// Defines values for EmployeeDepartment.
const (
	EmployeeDepartmentDesign         EmployeeDepartment = "Design"
	EmployeeDepartmentEngineering    EmployeeDepartment = "Engineering"
	EmployeeDepartmentFinancial      EmployeeDepartment = "Financial"
	EmployeeDepartmentGeneralAffairs EmployeeDepartment = "General affairs"
	EmployeeDepartmentSales          EmployeeDepartment = "Sales"
)

// Defines values for NewEmployeeDepartment.
const (
	NewEmployeeDepartmentDesign         NewEmployeeDepartment = "Design"
	NewEmployeeDepartmentEngineering    NewEmployeeDepartment = "Engineering"
	NewEmployeeDepartmentFinancial      NewEmployeeDepartment = "Financial"
	NewEmployeeDepartmentGeneralAffairs NewEmployeeDepartment = "General affairs"
	NewEmployeeDepartmentSales          NewEmployeeDepartment = "Sales"
)

// Defines values for ListEmployeesParamsSortBy.
const (
	Department  ListEmployeesParamsSortBy = "department"
	Email       ListEmployeesParamsSortBy = "email"
	Name        ListEmployeesParamsSortBy = "name"
	OnboardDate ListEmployeesParamsSortBy = "onboardDate"
)

// Defines values for ListEmployeesParamsSortOrder.
const (
	ListEmployeesParamsSortOrderAsc  ListEmployeesParamsSortOrder = "asc"
	ListEmployeesParamsSortOrderDesc ListEmployeesParamsSortOrder = "desc"
)

// Defines values for ListDayOffsParamsSortBy.
const (
	DayOffType ListDayOffsParamsSortBy = "dayOffType"
	StartTime  ListDayOffsParamsSortBy = "startTime"
)

// Defines values for ListDayOffsParamsSortOrder.
const (
	ListDayOffsParamsSortOrderAsc  ListDayOffsParamsSortOrder = "asc"
	ListDayOffsParamsSortOrderDesc ListDayOffsParamsSortOrder = "desc"
)

// DayOffRecord defines model for DayOffRecord.
type DayOffRecord struct {
	DayOffType DayOffRecordDayOffType `json:"dayOffType"`

	// EmployeeID Unique id of the employee
	EmployeeID int64     `json:"employeeID"`
	EndTime    time.Time `json:"endTime"`
	Reason     string    `json:"reason"`
	StartTime  time.Time `json:"startTime"`
}

// DayOffRecordDayOffType defines model for DayOffRecord.DayOffType.
type DayOffRecordDayOffType string

// Employee defines model for Employee.
type Employee struct {
	Address    string              `json:"address"`
	Department EmployeeDepartment  `json:"department"`
	Email      openapi_types.Email `json:"email"`

	// Id Unique id of the employee
	Id int64 `json:"id"`

	// Name Name of the employee
	Name        string             `json:"name"`
	OnboardDate openapi_types.Date `json:"onboardDate"`
	PhoneNumber string             `json:"phoneNumber"`
	Salary      int                `json:"salary"`
}

// EmployeeDepartment defines model for Employee.Department.
type EmployeeDepartment string

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// ListEmployeesResponse defines model for ListEmployeesResponse.
type ListEmployeesResponse struct {
	Data []Employee `json:"data"`

	// Page Current page number
	Page int `json:"page"`

	// PageSize Number of items per page
	PageSize int `json:"pageSize"`

	// TotalCount Total number of records
	TotalCount int64 `json:"totalCount"`
}

// NewEmployee defines model for NewEmployee.
type NewEmployee struct {
	Address    string                `json:"address"`
	Department NewEmployeeDepartment `json:"department"`
	Email      openapi_types.Email   `json:"email"`

	// Name Name of the employee
	Name        string             `json:"name"`
	OnboardDate openapi_types.Date `json:"onboardDate"`
	PhoneNumber string             `json:"phoneNumber"`
	Salary      int                `json:"salary"`
}

// NewEmployeeDepartment defines model for NewEmployee.Department.
type NewEmployeeDepartment string

// Pong defines model for Pong.
type Pong struct {
	StartTime string `json:"startTime"`
}

// ListEmployeesParams defines parameters for ListEmployees.
type ListEmployeesParams struct {
	Page      *int                          `form:"page,omitempty" json:"page,omitempty"`
	PageSize  *int                          `form:"pageSize,omitempty" json:"pageSize,omitempty"`
	SortBy    *ListEmployeesParamsSortBy    `form:"sortBy,omitempty" json:"sortBy,omitempty"`
	SortOrder *ListEmployeesParamsSortOrder `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`

	// Filters Key-value pairs for filtering records (e.g., filters[department]=Engineering&filters[name]=John)
	Filters *map[string]string `json:"filters,omitempty"`
}

// ListEmployeesParamsSortBy defines parameters for ListEmployees.
type ListEmployeesParamsSortBy string

// ListEmployeesParamsSortOrder defines parameters for ListEmployees.
type ListEmployeesParamsSortOrder string

// CancelDayOffJSONBody defines parameters for CancelDayOff.
type CancelDayOffJSONBody struct {
	CancellationReason string `json:"cancellationReason"`
}

// ListDayOffsParams defines parameters for ListDayOffs.
type ListDayOffsParams struct {
	Page          *int                        `form:"page,omitempty" json:"page,omitempty"`
	PageSize      *int                        `form:"pageSize,omitempty" json:"pageSize,omitempty"`
	SortBy        *ListDayOffsParamsSortBy    `form:"sortBy,omitempty" json:"sortBy,omitempty"`
	SortOrder     *ListDayOffsParamsSortOrder `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
	StartTimeFrom *openapi_types.Date         `form:"startTimeFrom,omitempty" json:"startTimeFrom,omitempty"`
	StartTimeTo   *openapi_types.Date         `form:"startTimeTo,omitempty" json:"startTimeTo,omitempty"`

	// Filters Key-value pairs for filtering records (e.g., filters[dayOffType]=PTO)
	Filters *map[string]string `json:"filters,omitempty"`
}

// ListDayOffsParamsSortBy defines parameters for ListDayOffs.
type ListDayOffsParamsSortBy string

// ListDayOffsParamsSortOrder defines parameters for ListDayOffs.
type ListDayOffsParamsSortOrder string

// AddEmployeeJSONRequestBody defines body for AddEmployee for application/json ContentType.
type AddEmployeeJSONRequestBody = NewEmployee

// CancelDayOffJSONRequestBody defines body for CancelDayOff for application/json ContentType.
type CancelDayOffJSONRequestBody CancelDayOffJSONBody

// UpdateEmployeeJSONRequestBody defines body for UpdateEmployee for application/json ContentType.
type UpdateEmployeeJSONRequestBody = NewEmployee

// SubmitDayOffJSONRequestBody defines body for SubmitDayOff for application/json ContentType.
type SubmitDayOffJSONRequestBody = DayOffRecord
