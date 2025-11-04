package entity

type UserEntity struct {
	BaseEntity
	RealName string `json:"realName"`                // 姓名
	UserName string `json:"userName"`                // 用户名
	Email    string `json:"email"`                   // 邮箱
	Password string `json:"password"`                // 密码
	Status   int    `gorm:"default:1" json:"status"` //状态: 1-启用，2-禁用
	RoleCode string `json:"roleCode"`                // 角色编号
}

func (UserEntity) TableName() string {
	return "sys_user"
}
