package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"com.example/example/pkg/cache"
	"com.example/example/repository"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/spf13/cast"
)

const (
	WildcardToken        = "*"
	PartDividerToken     = ":"
	SubPartDividerToken  = ","
	DefaultCaseSensitive = true
)

type AuthorizeManager struct {
	cache          cache.Storage
	userRoleRepo   *repository.UserRoleRepository
	permissionRepo *repository.PermissionRepository
}

func NewAuthorizeManager(cache cache.Storage,
	userRoleRepo *repository.UserRoleRepository,
	permissionRepo *repository.PermissionRepository) *AuthorizeManager {
	return &AuthorizeManager{
		cache:          cache,
		userRoleRepo:   userRoleRepo,
		permissionRepo: permissionRepo,
	}
}

func (a *AuthorizeManager) GetAuthorizeInfo(uid string) (*AuthorizeInfo, error) {
	// 加载用户的权限信息
	var info AuthorizeInfo
	val, _ := a.cache.Get(a.getAuthorizeCacheKey(uid))
	if val != nil {
		// 存在缓存
		err := json.Unmarshal([]byte(val.(string)), &info)
		if err != nil {
			slog.Error("invalid data in cache", "error", err)
			return &info, err
		}
	} else {
		// 数据库中查询
		res, err := a.loadAuthorizeInfo(uid)
		if err != nil {
			slog.Error("load authorizeInfo error", "error", err)
			return &info, err
		}
		if res == nil {
			return &info, err
		}
		info = *res
		// 写入缓存
		data, err := json.Marshal(info)
		if err == nil {
			err := a.cache.Set(a.getAuthorizeCacheKey(uid), string(data), 10*time.Minute)
			if err != nil {
				slog.Error("cache set value error", "error", err)
			}
		}
	}
	return &info, nil
}

func (a *AuthorizeManager) getAuthorizeCacheKey(uid string) string {
	// 获取缓存权限信息的key
	return fmt.Sprintf("authz:uid:%s", uid)
}

// HasPermission 是否拥有权限
func (a *AuthorizeManager) HasPermission(uid string, permission string) bool {
	return a.HasAllPermission(uid, []string{permission})
}

// HasAllPermission 是否拥有所有权限
func (a *AuthorizeManager) HasAllPermission(uid string, permissions []string) bool {
	info, err := a.GetAuthorizeInfo(uid)
	if err != nil || info == nil {
		return false
	}
	for _, perm := range permissions {
		var allow bool
		for _, p := range info.Permissions {
			if a.matchPermission(p, perm) {
				allow = true
				break
			}
		}
		if !allow {
			return false
		}
	}
	return true
}

// HasAnyPermission 是否拥有任一权限
func (a *AuthorizeManager) HasAnyPermission(uid string, permissions []string) bool {
	info, err := a.GetAuthorizeInfo(uid)
	if err != nil || info == nil {
		return false
	}
	for _, perm := range permissions {
		for _, p := range info.Permissions {
			if a.matchPermission(p, perm) {
				return true
			}
		}
	}
	return false
}

// HasRole 是否拥有角色
func (a *AuthorizeManager) HasRole(uid string, role string) bool {
	return a.HasAllRole(uid, []string{role})
}

// HasAllRole 是否拥有所有角色
func (a *AuthorizeManager) HasAllRole(uid string, roles []string) bool {
	info, err := a.GetAuthorizeInfo(uid)
	if err != nil || info == nil {
		return false
	}
	cachedRoles := mapset.NewSet(info.Roles...)
	for _, role := range roles {
		if !cachedRoles.ContainsOne(role) {
			return false
		}
	}
	return true
}

// HasAnyRole 是否拥有任一角色
func (a *AuthorizeManager) HasAnyRole(uid string, roles []string) bool {
	info, err := a.GetAuthorizeInfo(uid)
	if err != nil || info == nil {
		return false
	}
	cachedRoles := mapset.NewSet(info.Roles...)
	for _, role := range roles {
		if cachedRoles.ContainsOne(role) {
			return true
		}
	}
	return false
}

// matchPermission Matches permission with wildcards support | 权限匹配（支持通配符）
func (a *AuthorizeManager) matchPermission(pattern, permission string) bool {
	// Exact match or wildcard | 精确匹配或通配符
	if pattern == WildcardToken || pattern == permission {
		return true
	}
	// 支持通配符，例如 user:* 匹配 user:add, user:delete等
	wildcardSuffix := PartDividerToken + WildcardToken
	if strings.HasSuffix(pattern, wildcardSuffix) {
		prefix := strings.TrimSuffix(pattern, WildcardToken)
		return strings.HasPrefix(permission, prefix)
	}
	// 支持 user:*:view 这样的模式
	if strings.Contains(pattern, WildcardToken) {
		parts := strings.Split(pattern, PartDividerToken)
		permParts := strings.Split(permission, PartDividerToken)
		if len(parts) != len(permParts) {
			return false
		}
		for i, part := range parts {
			if part != WildcardToken && part != permParts[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (a *AuthorizeManager) ClearCache(uid ...string) error {
	var keys []string
	for _, id := range uid {
		keys = append(keys, a.getAuthorizeCacheKey(id))
	}
	_ = a.cache.Delete(keys...)
	return nil
}

func (a *AuthorizeManager) loadAuthorizeInfo(uid string) (*AuthorizeInfo, error) {
	// 加载用户的权限信息
	userId := cast.ToInt64(uid)
	roles, err := a.userRoleRepo.FindRoleCodeByUserId(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	permissions, err := a.permissionRepo.FindByUserId(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	perms := mapset.NewSet[string]()
	for _, p := range permissions {
		if p.Action != "" {
			perms.Add(p.Action)
		}
	}
	info := &AuthorizeInfo{
		Uid:         uid,
		Permissions: perms.ToSlice(),
		Roles:       roles,
	}
	return info, nil
}
