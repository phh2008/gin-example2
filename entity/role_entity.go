package entity

type RoleEntity struct {
	BaseEntity
	RoleCode string // 角色编码
	RoleName string // 角色名称
}

func (RoleEntity) TableName() string {
	return "sys_role"
}
