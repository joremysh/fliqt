package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"

	"github.com/joremysh/fliqt/api"
	"github.com/joremysh/fliqt/internal/model"
	"github.com/joremysh/fliqt/internal/repository"
	"github.com/joremysh/fliqt/internal/service"
	"github.com/joremysh/fliqt/pkg/cache"
)

var _ api.ServerInterface = (*HRSystem)(nil)
var StartUp string

type HRSystem struct {
	gdb             *gorm.DB
	employeeService service.EmployeeService
	dayOffService   service.DayOffService
}

func NewHRSystem(gdb *gorm.DB, redisClient *cache.RedisClient) *HRSystem {
	employeeRepo := repository.NewEmployeeRepo(gdb)
	dayOffRepo := repository.NewDayOffRepo(gdb)

	return &HRSystem{
		gdb:             gdb,
		employeeService: service.NewEmployeeService(employeeRepo, redisClient),
		dayOffService:   service.NewDayOffService(dayOffRepo, employeeRepo),
	}
}

func (s *HRSystem) GetLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, api.Pong{
		StartTime: StartUp,
	})
}

func ConvertToEmployeeResponse(employee *model.Employee) *api.Employee {
	return &api.Employee{
		Address:     employee.Address,
		Email:       openapitypes.Email(employee.Email),
		Id:          int64(employee.ID),
		Name:        employee.Name,
		OnboardDate: openapitypes.Date{Time: employee.OnboardDate},
		PhoneNumber: employee.PhoneNumber,
		Salary:      employee.Salary,
		Department:  api.EmployeeDepartment(employee.Department),
	}
}

func ConvertToDayOffResponse(record *model.DayOffRecord) *api.DayOffRecord {
	return &api.DayOffRecord{
		DayOffType: api.DayOffRecordDayOffType(record.DayOffType),
		EmployeeID: int64(record.EmployeeID),
		EndTime:    record.EndTime,
		Reason:     record.Reason,
		StartTime:  record.StartTime,
	}
}

func (s *HRSystem) ListEmployees(c *gin.Context, params api.ListEmployeesParams) {
	result, err := s.employeeService.ListEmployees(c.Request.Context(), parseListParams(params))
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &api.ListEmployeesResponse{
		Data:       make([]api.Employee, len(result.Data)),
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

func parseListParams(params api.ListEmployeesParams) *model.ListParams {
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
	c.JSON(code, api.Error{
		Code:    code,
		Message: errMsg,
	})
}

func (s *HRSystem) AddEmployee(c *gin.Context) {
	var newEmployee api.NewEmployee
	err := c.Bind(&newEmployee)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for Employee")
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
	var newEmployee api.NewEmployee
	err := c.Bind(&newEmployee)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for Employee")
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

	c.JSON(http.StatusNoContent, id)
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
	var dayOffRecord api.DayOffRecord
	err := c.Bind(&dayOffRecord)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for DayOff Record")
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

func (s *HRSystem) ListDayOffs(c *gin.Context, id int64, params api.ListDayOffsParams) {
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

	result, err := s.dayOffService.ListDayOffs(c.Request.Context(), uint(id), listParams)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &api.ListDayOffsResponse{
		Data:       make([]api.DayOffRecord, len(result.Data)),
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalCount: result.TotalCount,
	}
	for i, datum := range result.Data {
		converted := ConvertToDayOffResponse(&datum)
		resp.Data[i] = *converted
	}
	c.JSON(http.StatusOK, resp)
}

func (s *HRSystem) CancelDayOff(c *gin.Context, id int64) {
	var request api.CancelDayOffJSONBody
	err := c.Bind(&request)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid format for Cancel Day Off")
		return
	}

	if err = s.dayOffService.CancelDayOff(c.Request.Context(), uint(id), request.CancellationReason); err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, id)
}
