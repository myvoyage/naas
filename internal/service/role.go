package service

import (
	"context"
	"errors"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/pkg/logger"
)

type role struct {
}

func (r *role) Recursive() []*model.Role {
	var (
		rootRoles []*model.Role
		err       error
	)
	ctx := store.NewDBContext()
	rootRoles, err = dao.Role.SelectByRoot(ctx)
	if err != nil {
		logger.Errorf("dao.Role.SelectByRoot: %s", err)
	}
	r.recursive(ctx, rootRoles)
	return rootRoles
}

func (r *role) recursive(ctx context.Context, roles []*model.Role) {
	if len(roles) == 0 {
		return
	}
	var (
		childRoles []*model.Role
		err        error
	)
	for _, role := range roles {
		childRoles, err = dao.Role.SelectByParentCode(ctx, role.Code)
		if err != nil {
			logger.Errorf("dao.Role.SelectByRoot: %s", err)
		}
		r.recursive(ctx, childRoles)
		role.ChildRoles = childRoles
	}
}

// GetAllRoleByUserID ...
func (r *role) GetAllRoleByUserID(userID uint64) (roles []*model.UserRole, err error) {
	ctx := store.NewDBContext()
	roles, err = dao.UserRole.SelectAllByUserID(ctx, userID)
	return
}

// AddResourceWebRoute 添加web路由资源角色
func (r *role) AddResourceWebRoute(roleCode string, resourceWebRouteID uint64) (err error) {
	ctx := store.NewDBContext()
	var exist bool
	exist, err = dao.RoleResourceWebRoute.ExistByRoleCodeAndResourceWebRouteID(ctx, roleCode, resourceWebRouteID)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("Web Routing conditions exist")
		return
	}
	err = dao.RoleResourceWebRoute.Insert(ctx, &model.RoleResourceWebRoute{
		RoleCode:           roleCode,
		ResourceWebRouteID: resourceWebRouteID,
	})
	if err != nil {
		return
	}
	var resourceWebRoute *model.ResourceWebRoute
	resourceWebRoute, err = dao.ResourceWebRoute.Select(ctx, resourceWebRouteID)
	if err != nil {
		return
	}
	sub, dom, obj, act := formatPolicy(roleCode, resourceWebRoute)
	_, err = casbin.Enforcer.AddPolicy(sub, dom, obj, act)
	if err != nil {
		logger.Errorf("casbin.Enforcer.AddPolicy: %s", err)
	}
	return
}
