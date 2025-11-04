package auth

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

const (
	WildcardToken        = "*"
	PartDividerToken     = ":"
	SubPartDividerToken  = ","
	DefaultCaseSensitive = true
)

type AuthorizeManager struct {
}

func (a *AuthorizeManager) GetAuthorizeInfo(uid string) (*AuthorizeInfo, error) {
	// TODO 加载用户的权限信息
	return nil, nil
}

func (a *AuthorizeManager) getAuthorizeCacheKey(uid string) string {
	// TODO 获取缓存权限信息的key
	return ""
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
	cachedPermissions := mapset.NewSet(info.Permissions...)
	for _, perm := range permissions {
		if !cachedPermissions.ContainsOne(perm) {
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
	cachedPermissions := mapset.NewSet(info.Permissions...)
	for _, perm := range permissions {
		if cachedPermissions.ContainsOne(perm) {
			return true
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
