package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// RoleResourceWebRouter ...
type RoleResourceWebRouter interface {
	Insert(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (roleResourceWebRoute *model.RoleResourceWebRoute, err error)
	SelectAll(ctx context.Context) (roleResourceWebRoutes []*model.RoleResourceWebRoute, err error)
	Update(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error)
	ExistByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode string, resourceWebRouteID uint64) (exist bool, err error)
}

type roleResourceWebRoute struct {
}

func (*roleResourceWebRoute) Insert(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(roleResourceWebRoute).Error
	return
}
func (*roleResourceWebRoute) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.RoleResourceWebRoute{}, id).Error
	return
}
func (*roleResourceWebRoute) Select(ctx context.Context, id uint64) (roleResourceWebRoute *model.RoleResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.RoleResourceWebRoute
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	roleResourceWebRoute = &dbResult
	return
}
func (*roleResourceWebRoute) Update(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(roleResourceWebRoute).Update(roleResourceWebRoute).Error
	return
}

func (*roleResourceWebRoute) SelectAll(ctx context.Context) (roleResourceWebRoutes []*model.RoleResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.RoleResourceWebRoute{}).Find(&roleResourceWebRoutes).Error
	return
}

func (*roleResourceWebRoute) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var count uint
	err = gdb.Model(&model.RoleResourceWebRoute{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// ExistByRoleCodeAndResourceWebRouteID 判断根据RoleCode和资源web路由ID
func (r *roleResourceWebRoute) ExistByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode string, resourceWebRouteID uint64) (exist bool, err error) {
	return r.exist(ctx, "role_code = ? and resource_web_route_id = ?", roleCode, resourceWebRouteID)
}
