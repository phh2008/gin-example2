package entity

type PermissionEntity struct {
	BaseEntity
	PermName string // 权限名称
	Url      string // URL路径
	Action   string // 权限动作：比如get、post、delete等
	PermType uint8  `gorm:"default:1"` // 权限类型：1-菜单、2-按钮
	ParentId int64  `gorm:"default:0"` // 父级ID：资源层级关系（0表示顶级）
}

func (PermissionEntity) TableName() string {
	return "sys_permission"
}
