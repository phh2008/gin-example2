package entity

type UserRoleEntity struct {
	Id     int64
	UserId int64
	RoleId int64
}

func (UserRoleEntity) TableName() string {
	return "sys_user_role"
}
