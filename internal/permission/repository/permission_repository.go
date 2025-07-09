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

// PermissionRepository 权限数据访问接口
type PermissionRepository interface {
	// Create 创建权限
	Create(permission *models.Permission) error

	// GetByID 通过ID获取权限
	GetByID(id string) (*models.Permission, error)

	// GetByName 通过名称获取权限
	GetByName(name string) (*models.Permission, error)

	// List 获取权限列表
	List(page, pageSize int) ([]*models.Permission, int64, error)

	// Update 更新权限
	Update(id string, data *models.Permission) error

	// Delete 删除权限
	Delete(id string) error

	// GetByCategory 按类别获取权限
	GetByCategory(category string) ([]*models.Permission, error)

	// GetByResource 按资源获取权限
	GetByResource(resource string) ([]*models.Permission, error)

	// GetByAction 按操作获取权限
	GetByAction(action string) ([]*models.Permission, error)

	// GetByResourceAndAction 按资源和操作获取权限
	GetByResourceAndAction(resource, action string) (*models.Permission, error)
}

// permissionRepository 权限仓储实现
type permissionRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewPermissionRepository 创建权限仓储
func NewPermissionRepository(db *mongo.Database) PermissionRepository {
	return &permissionRepository{
		db:         db,
		collection: db.Collection("permissions"),
	}
}

// Create 创建权限
func (r *permissionRepository) Create(permission *models.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间
	permission.CreatedAt = time.Now()

	// 插入权限
	result, err := r.collection.InsertOne(ctx, permission)
	if err != nil {
		return err
	}

	// 设置生成的ID
	permission.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 通过ID获取权限
func (r *permissionRepository) GetByID(id string) (*models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid permission ID format")
	}

	var permission models.Permission
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

// GetByName 通过名称获取权限
func (r *permissionRepository) GetByName(name string) (*models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var permission models.Permission
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

// List 获取权限列表
func (r *permissionRepository) List(page, pageSize int) ([]*models.Permission, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "category", Value: 1}, {Key: "resource", Value: 1}, {Key: "action", Value: 1}}) // 按分类、资源、操作排序

	// 查询权限
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// Update 更新权限
func (r *permissionRepository) Update(id string, data *models.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid permission ID format")
	}

	// 创建更新文档
	updateDoc := bson.M{"$set": data}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("permission not found")
	}

	return nil
}

// Delete 删除权限
func (r *permissionRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid permission ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("permission not found")
	}

	return nil
}

// GetByCategory 按类别获取权限
func (r *permissionRepository) GetByCategory(category string) ([]*models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetByResource 按资源获取权限
func (r *permissionRepository) GetByResource(resource string) ([]*models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"resource": resource})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetByAction 按操作获取权限
func (r *permissionRepository) GetByAction(action string) ([]*models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"action": action})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetByResourceAndAction 按资源和操作获取权限
func (r *permissionRepository) GetByResourceAndAction(resource, action string) (*models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var permission models.Permission
	err := r.collection.FindOne(ctx, bson.M{"resource": resource, "action": action}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}
