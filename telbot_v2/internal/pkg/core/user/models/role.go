// This model contains list of accepted user roles

package models

const (
	roleAdminId    = 1
	roleUserId     = 2
	roleAadminName = "Admin"
	roleUserName   = "User"
)

type RoleId map[string]uint8
type RoleName map[uint8]string

// We use 2 symmetric maps in order to provide quick access to data by roles' Id and Name
var roleNames RoleName
var roleIds RoleId

func init() {
	roleNames = RoleName{}
	roleNames[roleAdminId] = roleAadminName
	roleNames[roleUserId] = roleUserName

	roleIds = RoleId{}
	roleIds[roleAadminName] = roleAdminId
	roleIds[roleUserName] = roleUserId
}

func GetRoleName(id uint8) string {
	return roleNames[id]
}

func GetRoleId(roleName string) uint8 {
	return roleIds[roleName]
}
