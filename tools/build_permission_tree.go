package tools

import "github.com/sanyuanya/dongle/entity"

func BuildPermissionTree(permissions []*entity.PermissionMenu) []*entity.PermissionMenu {
	permissionMap := make(map[string]*entity.PermissionMenu)
	roots := make([]*entity.PermissionMenu, 0)

	// 将权限列表转换为映射
	for _, perm := range permissions {
		permissionMap[perm.SnowflakeId] = perm
	}

	// 构建父子关系
	for _, perm := range permissions {
		if perm.ParentId == "" {
			roots = append(roots, perm)
		} else {
			parent, exists := permissionMap[perm.ParentId]
			if exists {
				parent.Children = append(parent.Children, perm)
			}
		}
	}

	return roots
}
