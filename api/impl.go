package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
	"github.com/joremysh/fliqt/internal/repository"
	"github.com/joremysh/fliqt/internal/service"
)

var _ ServerInterface = (*HRSystem)(nil)

type HRSystem struct {
	gdb             *gorm.DB
	employeeService service.EmployeeService
	dayOffService   service.DayOffService
}

func NewHRSystem(gdb *gorm.DB) *HRSystem {
	employeeRepo := repository.NewEmployeeRepo(gdb)
	dayOffRepo := repository.NewDayOffRepo(gdb)

	return &HRSystem{
		gdb:             gdb,
		employeeService: service.NewEmployeeService(employeeRepo),
		dayOffService:   service.NewDayOffService(dayOffRepo, employeeRepo),
	}
}

func (s *HRSystem) GetLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, Pong{
		StartTime: time.Now().Format(time.RFC3339),
	})
}

func ConvertToEmployeeResponse(employee *model.Employee) *Employee {
	return &Employee{
		Address:     employee.Address,
		Email:       openapi_types.Email(employee.Email),
		Id:          int64(employee.ID),
		Name:        employee.Name,
		OnboardDate: openapi_types.Date{Time: employee.OnboardDate},
		PhoneNumber: employee.PhoneNumber,
		Salary:      employee.Salary,
		Department:  EmployeeDepartment(employee.Department),
	}
}

func ConvertToDayOffResponse(record *model.DayOffRecord) *DayOffRecord {
	return &DayOffRecord{
		DayOffType: DayOffRecordDayOffType(record.DayOffType),
		EmployeeID: int64(record.EmployeeID),
		EndTime:    record.EndTime,
		Reason:     record.Reason,
		StartTime:  record.StartTime,
	}
}

func (s *HRSystem) ListEmployees(c *gin.Context, params ListEmployeesParams) {
	result, err := s.employeeService.ListEmployees(c.Request.Context(), parseListParams(params))
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &ListEmployeesResponse{
		Data:       make([]Employee, len(result.Data)),
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalCount: result.TotalCount,
	}
	for i, employee := range result.Data {
		converted := ConvertToEmployeeResponse(&employee)
		resp.Data[i] = *converted
	}
	c.JSON(http.StatusOK, resp)
}

func parseListParams(params ListEmployeesParams) *model.ListParams {
	listParams := &model.ListParams{}
	if params.PageSize != nil {
		listParams.PageSize = *params.PageSize
	}
	if params.Page != nil {
		listParams.Page = *params.Page
	}
	if params.SortBy != nil {
		listParams.SortBy = string(*params.SortBy)
	}
	if params.SortOrder != nil {
		listParams.SortOrder = string(*params.SortOrder)
	}
	if params.Filters != nil {
		listParams.Filters = *params.Filters
	}
	return listParams
}

func sendErrorResponse(c *gin.Context, code int, errMsg string) {
	c.JSON(code, Error{
		Code:    int32(code),
		Message: errMsg,
	})
}

func (s *HRSystem) AddEmployee(c *gin.Context) {
	var newEmployee NewEmployee
	err := c.Bind(&newEmployee)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for NewEmployee")
		return
	}

	created, err := s.employeeService.CreateEmployee(c.Request.Context(), &model.Employee{
		Name:        newEmployee.Name,
		Email:       string(newEmployee.Email),
		PhoneNumber: newEmployee.PhoneNumber,
		Department:  string(newEmployee.Department),
		Address:     newEmployee.Address,
		Salary:      newEmployee.Salary,
		OnboardDate: newEmployee.OnboardDate.Time,
	})
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, ConvertToEmployeeResponse(created))
}

func (s *HRSystem) UpdateEmployee(c *gin.Context, id int64) {
	var newEmployee NewEmployee
	err := c.Bind(&newEmployee)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for NewEmployee")
		return
	}

	req := &model.Employee{
		Name:        newEmployee.Name,
		Email:       string(newEmployee.Email),
		PhoneNumber: newEmployee.PhoneNumber,
		Department:  string(newEmployee.Department),
		Address:     newEmployee.Address,
		Salary:      newEmployee.Salary,
		OnboardDate: newEmployee.OnboardDate.Time,
	}
	req.ID = uint(id)

	updated, err := s.employeeService.UpdateEmployee(c.Request.Context(), req)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ConvertToEmployeeResponse(updated))
}

func (s *HRSystem) DeleteEmployee(c *gin.Context, id int64) {
	if err := s.employeeService.DeleteEmployee(c.Request.Context(), uint(id)); err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, id)
}

func (s *HRSystem) FindEmployeeByID(c *gin.Context, id int64) {
	employee, err := s.employeeService.GetEmployee(c.Request.Context(), uint(id))
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ConvertToEmployeeResponse(employee))
}

func (s *HRSystem) SubmitDayOff(c *gin.Context, id int64) {
	var dayOffRecord DayOffRecord
	err := c.Bind(&dayOffRecord)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for DayOffRecord")
		return
	}

	created, err := s.dayOffService.SubmitDayOff(c.Request.Context(), &model.DayOffRecord{
		EmployeeID: uint(id),
		DayOffType: string(dayOffRecord.DayOffType),
		Reason:     dayOffRecord.Reason,
		StartTime:  dayOffRecord.StartTime,
		EndTime:    dayOffRecord.EndTime,
	})
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, ConvertToDayOffResponse(created))
}

func (s *HRSystem) ListDayOffs(c *gin.Context, id int64, params ListDayOffsParams) {
	listParams := &model.ListParams{}
	if params.PageSize != nil {
		listParams.PageSize = *params.PageSize
	}
	if params.Page != nil {
		listParams.Page = *params.Page
	}
	if params.SortBy != nil {
		listParams.SortBy = string(*params.SortBy)
	}
	if params.SortOrder != nil {
		listParams.SortOrder = string(*params.SortOrder)
	}
	if params.Filters != nil {
		listParams.Filters = *params.Filters
	}

	// result, err := s.employeeService.ListEmployees(c.Request.Context(), parseListParams(params))
	// if err != nil {
	// 	sendErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	//
	// resp := &ListEmployeesResponse{
	// 	Data:       make([]Employee, len(result.Data)),
	// 	Page:       result.Page,
	// 	PageSize:   result.PageSize,
	// 	TotalCount: result.TotalCount,
	// }
	// for i, employee := range result.Data {
	// 	converted := ConvertToEmployeeResponse(&employee)
	// 	resp.Data[i] = *converted
	// }
	// c.JSON(http.StatusOK, resp)

	// result,err:=s.dayOffService.ListDayOffs(c.Request.Context(), uint(id),listParams)
	// if err != nil {
	// 	sendErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// TODO implement me
	panic("implement me")
}

func (s *HRSystem) CancelDayOff(c *gin.Context, id int64) {
	// TODO implement me
	panic("implement me")
}