package service

import (
	"context"
	"strconv"
	"time"

	"com.example/example/entity"
	"com.example/example/model"
	"com.example/example/model/result"
	"com.example/example/pkg/exception"
	"com.example/example/pkg/logger"
	"com.example/example/pkg/xjwt"
	"com.example/example/repository"
	"github.com/casbin/casbin/v2"
	"github.com/cristalhq/jwt/v5"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepository *repository.UserRepository
	jwt            *xjwt.JwtHelper
	enforcer       *casbin.Enforcer
}

// NewUserService 创建服务
func NewUserService(
	userRepository *repository.UserRepository,
	jwt *xjwt.JwtHelper,
	enforcer *casbin.Enforcer,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		jwt:            jwt,
		enforcer:       enforcer,
	}
}

// ListPage 用户列表
func (a *UserService) ListPage(ctx context.Context, req model.UserListReq) *result.Result[model.PageData[model.UserModel]] {
	data := a.userRepository.ListPage(ctx, req)
	return result.Ok[model.PageData[model.UserModel]](data)
}

// CreateByEmail 根据邮箱创建用户
func (a *UserService) CreateByEmail(ctx context.Context, email model.UserEmailRegister) *result.Result[model.UserModel] {
	user := a.userRepository.GetByEmail(ctx, email.Email)
	if user.Id > 0 {
		return result.FailMsg[model.UserModel]("email 已存在")
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(email.Password), 1)
	if err != nil {
		logger.Errorf("生成密码出错：%s", err.Error())
		return result.Error[model.UserModel](err)
	}

	user = entity.UserEntity{
		Email:    email.Email,
		RealName: email.Email,
		UserName: email.Email,
		Password: string(pwd),
		Status:   1,
		RoleCode: "",
	}
	user, err = a.userRepository.Add(ctx, user)
	if err != nil {
		logger.Errorf("创建用户出错：%s", err.Error())
		return result.FailMsg[model.UserModel]("创建用户出错")
	}
	var userModel model.UserModel
	copier.Copy(&userModel, &user)
	return result.Ok[model.UserModel](userModel)
}

// LoginByEmail 邮箱登录
func (a *UserService) LoginByEmail(ctx context.Context, loginModel model.UserLoginModel) *result.Result[string] {
	user := a.userRepository.GetByEmail(ctx, loginModel.Email)
	if user.Id == 0 {
		return result.FailMsg[string]("用户或密码错误")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginModel.Password))
	if err != nil {
		return result.FailMsg[string]("用户或密码错误")
	}
	// 生成token
	userClaims := xjwt.UserClaims{}
	userClaims.ID = strconv.FormatInt(user.Id, 10)
	userClaims.Role = user.RoleCode
	userClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
	token, err := a.jwt.CreateToken(userClaims)
	if err != nil {
		logger.Errorf("生成token错误：%s", err.Error())
		return result.Error[string](exception.SysError)
	}
	return result.Ok[string](token.String())
}

// AssignRole 给用户分配角色
func (a *UserService) AssignRole(ctx context.Context, userRole model.AssignRoleModel) *result.Result[any] {
	err := a.userRepository.SetRole(ctx, userRole.UserId, userRole.RoleCode)
	if err != nil {
		logger.Errorf("db update error: %s", err.Error())
		return result.FailMsg[any]("分配角色出错")
	}
	// 更新casbin中的用户与角色关系
	uid := strconv.FormatInt(userRole.UserId, 10)
	_, _ = a.enforcer.DeleteRolesForUser(uid)
	// 角色为空，表示清除此用户的角色,无需添加
	if userRole.RoleCode != "" {
		_, _ = a.enforcer.AddGroupingPolicy(uid, userRole.RoleCode)
	}
	return result.Success[any]()
}

// DeleteById 根据ID删除
func (a *UserService) DeleteById(ctx context.Context, id int64) *result.Result[any] {
	err := a.userRepository.DeleteById(ctx, id)
	if err != nil {
		logger.Errorf("delete error: %s", err.Error())
		return result.FailMsg[any]("刪除出错")
	}
	// 清除 casbin 中用户信息
	_, err = a.enforcer.DeleteRolesForUser(strconv.FormatInt(id, 10))
	if err != nil {
		logger.Errorf("Enforcer.DeleteRolesForUser error: %s", err)
	}
	return result.Success[any]()
}
