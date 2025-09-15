package allmaster

import (
	"mjsubackend/model/master/department"
	"mjsubackend/model/master/departmentform"
	"mjsubackend/model/master/form"
	"mjsubackend/model/master/role"
	"mjsubackend/model/master/roleform"
	"mjsubackend/model/master/userrole"
)

type MasterData struct {
	Department     []department.Department         `json:"department"`
	DepartmentForm []departmentform.DepartmentForm `json:"department_form"`
	Form           []form.Form                     `json:"form"`
	Role           []role.Role                     `json:"role"`
	RoleForm       []roleform.RoleForm             `json:"role_form"`
	UserRole       []userrole.UserRole             `json:"user_role"`
}
