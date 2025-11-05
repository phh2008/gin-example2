package service

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"com.example/example/entity"
	"com.example/example/manager"
	"com.example/example/model"
	"com.example/example/model/result"
	"com.example/example/pkg/exception"
	"com.example/example/pkg/xjwt"
	"com.example/example/repository"
	"github.com/cristalhq/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepository     *repository.UserRepository
	userRoleRepository *repository.UserRoleRepository
	authorizeManager   *manager.AuthorizeManager
	jwt                *xjwt.JwtHelper
}

// NewUserService 创建服务
func NewUserService(
	userRepository *repository.UserRepository,
	userRoleRepository *repository.UserRoleRepository,
	authorizeManager *manager.AuthorizeManager,
	jwt *xjwt.JwtHelper,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRoleRepository: userRoleRepository,
		authorizeManager:   authorizeManager,
		jwt:                jwt,
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
	pwd, err := bcrypt.GenerateFromPassword([]byte(email.Password), 10)
	if err != nil {
		slog.Error("生成密码出错", "error", err)
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
		slog.Error("创建用户出错", "error", err)
		return result.FailMsg[model.UserModel]("创建用户出错")
	}
	var userModel model.UserModel
	err = copier.Copy(&userModel, &user)
	if err != nil {
		slog.Error("copier.Copy拷贝出错", "error", err.Error())
		return result.Error[model.UserModel](err)
	}
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
		slog.Error("生成token错误", "error", err)
		return result.Error[string](exception.SysError)
	}
	return result.Ok[string](token.String())
}

// AssignRole 给用户分配角色
func (a *UserService) AssignRole(ctx context.Context, userRole model.AssignRoleModel) *result.Result[any] {
	err := a.userRepository.Transaction(ctx, func(tx context.Context) error {
		err := a.userRoleRepository.DeleteByUserId(tx, userRole.UserId)
		if err != nil {
			return err
		}
		err = a.userRoleRepository.CreateBatch(tx, []entity.UserRoleEntity{
			{
				UserId: userRole.UserId,
				RoleId: userRole.RoleId},
		})
		return err
	})
	if err != nil {
		slog.Error("分配角色出错", "error", err)
		return result.Error[any](err)
	}
	// 清空权限缓存
	_ = a.authorizeManager.ClearCache(cast.ToString(userRole.UserId))
	return result.Success[any]()
}

// DeleteById 根据ID删除
func (a *UserService) DeleteById(ctx context.Context, id int64) *result.Result[any] {
	err := a.userRepository.DeleteById(ctx, id)
	if err != nil {
		slog.Error("删除用户出错", "error", err)
		return result.FailMsg[any]("刪除出错")
	}
	// 清空分配的角色
	_ = a.userRoleRepository.DeleteByUserId(ctx, id)
	return result.Success[any]()
}
