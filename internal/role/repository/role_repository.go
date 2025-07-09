package repository

import (
	"authcenter/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleRepository 角色数据访问接口
type RoleRepository interface {
	// Create 创建角色
	Create(role *models.Role) error

	// GetByID 通过ID获取角色
	GetByID(id string) (*models.Role, error)

	// GetByName 通过名称获取角色
	GetByName(name string) (*models.Role, error)

	// List 获取角色列表
	List(page, pageSize int) ([]*models.Role, int64, error)

	// Update 更新角色
	Update(id string, data *models.Role) error

	// Delete 删除角色
	Delete(id string) error

	// AssignPermission 为角色分配权限
	AssignPermission(roleID, permissionID string) error

	// RemovePermission 移除角色权限
	RemovePermission(roleID, permissionID string) error

	// GetRolePermissions 获取角色权限
	GetRolePermissions(roleID string) ([]models.RolePermission, error)

	// GetRoleUsers 获取角色下的用户
	GetRoleUsers(roleID string) ([]*models.User, error)
}

// roleRepository 角色仓储实现
type roleRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *mongo.Database) RoleRepository {
	return &roleRepository{
		db:         db,
		collection: db.Collection("roles"),
	}
}

// Create 创建角色
func (r *roleRepository) Create(role *models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	// 默认状态为激活
	if role.Status == "" {
		role.Status = "active"
	}

	// 插入角色
	result, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return err
	}

	// 设置生成的ID
	role.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 通过ID获取角色
func (r *roleRepository) GetByID(id string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid role ID format")
	}

	var role models.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// GetByName 通过名称获取角色
func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// List 获取角色列表
func (r *roleRepository) List(page, pageSize int) ([]*models.Role, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "level", Value: -1}, {Key: "created_at", Value: -1}}) // 按级别和创建时间排序

	// 查询角色
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var roles []*models.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Update 更新角色
func (r *roleRepository) Update(id string, data *models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid role ID format")
	}

	// 设置更新时间
	data.UpdatedAt = time.Now()

	// 创建更新文档
	updateDoc := bson.M{"$set": data}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("role not found")
	}

	return nil
}

// Delete 删除角色
func (r *roleRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid role ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("role not found")
	}

	return nil
}

// AssignPermission 为角色分配权限
func (r *roleRepository) AssignPermission(roleID, permissionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	roleObjectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return errors.New("invalid role ID format")
	}

	permObjectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return errors.New("invalid permission ID format")
	}

	// 首先获取权限信息
	permCollection := r.db.Collection("permissions")
	var permission models.Permission
	err = permCollection.FindOne(ctx, bson.M{"_id": permObjectID}).Decode(&permission)
	if err != nil {
		return errors.New("permission not found")
	}

	// 创建角色权限对象
	rolePermission := models.RolePermission{
		PermissionID: permObjectID,
		Name:         permission.Name,
		Resource:     permission.Resource,
		Action:       permission.Action,
	}

	// 添加权限到角色
	update := bson.M{
		"$addToSet": bson.M{"permissions": rolePermission},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": roleObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("role not found")
	}

	return nil
}

// RemovePermission 移除角色权限
func (r *roleRepository) RemovePermission(roleID, permissionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	roleObjectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return errors.New("invalid role ID format")
	}

	permObjectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		return errors.New("invalid permission ID format")
	}

	// 移除权限
	update := bson.M{
		"$pull": bson.M{"permissions": bson.M{"permission_id": permObjectID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": roleObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("role not found")
	}

	return nil
}

// GetRolePermissions 获取角色权限
func (r *roleRepository) GetRolePermissions(roleID string) ([]models.RolePermission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return nil, errors.New("invalid role ID format")
	}

	var role models.Role
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return role.Permissions, nil
}

// GetRoleUsers 获取角色下的用户
func (r *roleRepository) GetRoleUsers(roleID string) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return nil, errors.New("invalid role ID format")
	}

	// 查询拥有该角色的用户
	userCollection := r.db.Collection("users")
	cursor, err := userCollection.Find(ctx, bson.M{"roles.role_id": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
