package repository

import (
	"context"
	"errors"
	"time"

	"authcenter/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// Create 创建用户
	Create(user *models.User) error

	// GetByID 通过ID获取用户
	GetByID(id string) (*models.User, error)

	// GetByEmail 通过邮箱获取用户
	GetByEmail(email string) (*models.User, error)

	// GetByUsername 通过用户名获取用户
	GetByUsername(username string) (*models.User, error)

	// GetByPhone 通过手机号获取用户
	GetByPhone(phone string) (*models.User, error)

	// List 获取用户列表
	List(page, pageSize int) ([]*models.User, int64, error)

	// Update 更新用户
	Update(id string, data *models.User) error

	// Delete 删除用户
	Delete(id string) error

	// AssignRole 为用户分配角色
	AssignRole(userID, roleID string, grantedBy string) error

	// RemoveRole 移除用户角色
	RemoveRole(userID, roleID string) error

	// GetUserRoles 获取用户角色
	GetUserRoles(userID string) ([]models.UserRole, error)

	// GetUserPermissions 获取用户权限
	GetUserPermissions(userID string) ([]models.RolePermission, error)

	// UpdateLoginHistory 更新登录历史
	UpdateLoginHistory(userID string, ip string) error

	// CheckUserExists 检查用户是否存在
	CheckUserExists(username, email, phone string) (bool, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 默认状态为激活
	if user.Status == "" {
		user.Status = "active"
	}

	// 插入用户
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// 设置生成的ID
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 通过ID获取用户
func (r *userRepository) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail 通过邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername 通过用户名获取用户
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByPhone 通过手机号获取用户
func (r *userRepository) GetByPhone(phone string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(page, pageSize int) ([]*models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // 按创建时间倒序

	// 查询用户
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update 更新用户
func (r *userRepository) Update(id string, data *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
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
		return errors.New("user not found")
	}

	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// AssignRole 为用户分配角色
func (r *userRepository) AssignRole(userID, roleID string, grantedBy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	roleObjectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return errors.New("invalid role ID format")
	}

	grantedByObjectID, err := primitive.ObjectIDFromHex(grantedBy)
	if err != nil {
		return errors.New("invalid granted by ID format")
	}

	// 首先获取角色信息
	roleCollection := r.db.Collection("roles")
	var role models.Role
	err = roleCollection.FindOne(ctx, bson.M{"_id": roleObjectID}).Decode(&role)
	if err != nil {
		return errors.New("role not found")
	}

	// 创建用户角色
	userRole := models.UserRole{
		RoleID:    roleObjectID,
		RoleName:  role.Name,
		GrantedBy: grantedByObjectID,
		GrantedAt: time.Now(),
	}

	// 添加角色到用户
	update := bson.M{
		"$addToSet": bson.M{"roles": userRole},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// RemoveRole 移除用户角色
func (r *userRepository) RemoveRole(userID, roleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	roleObjectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return errors.New("invalid role ID format")
	}

	// 移除角色
	update := bson.M{
		"$pull": bson.M{"roles": bson.M{"role_id": roleObjectID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserRoles 获取用户角色
func (r *userRepository) GetUserRoles(userID string) ([]models.UserRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user.Roles, nil
}

// GetUserPermissions 获取用户权限
func (r *userRepository) GetUserPermissions(userID string) ([]models.RolePermission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// 使用聚合查询获取用户的所有权限
	pipeline := []bson.M{
		{"$match": bson.M{"_id": objectID}},
		{"$unwind": "$roles"},
		{"$lookup": bson.M{
			"from":         "roles",
			"localField":   "roles.role_id",
			"foreignField": "_id",
			"as":           "role_detail",
		}},
		{"$unwind": "$role_detail"},
		{"$unwind": "$role_detail.permissions"},
		{"$group": bson.M{
			"_id":        "$role_detail.permissions",
			"permission": bson.M{"$first": "$role_detail.permissions"},
		}},
		{"$replaceRoot": bson.M{"newRoot": "$permission"}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []models.RolePermission
	if err = cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// UpdateLoginHistory 更新登录历史
func (r *userRepository) UpdateLoginHistory(userID string, ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"login_history.last_login_at": time.Now(),
			"login_history.last_ip":       ip,
			"updated_at":                  time.Now(),
		},
		"$inc": bson.M{
			"login_history.login_count": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CheckUserExists 检查用户是否存在
func (r *userRepository) CheckUserExists(username, email, phone string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	filter := bson.M{"$or": []bson.M{}}

	if username != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"username": username})
	}
	if email != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"email": email})
	}
	if phone != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"phone": phone})
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
