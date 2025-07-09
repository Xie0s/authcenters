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

// CategoryRepository 类别数据访问接口
type CategoryRepository interface {
	// Create 创建类别
	Create(category *models.Category) error

	// GetByID 通过ID获取类别
	GetByID(id string) (*models.Category, error)

	// GetByName 通过名称获取类别
	GetByName(name string) (*models.Category, error)

	// List 获取类别列表
	List(page, pageSize int) ([]*models.Category, int64, error)

	// Update 更新类别
	Update(id string, data *models.Category) error

	// Delete 删除类别
	Delete(id string) error

	// GetByParent 获取子类别
	GetByParent(parentID string) ([]*models.Category, error)

	// GetRootCategories 获取根类别
	GetRootCategories() ([]*models.Category, error)

	// GetCategoryTree 获取类别树
	GetCategoryTree() ([]*models.Category, error)

	// UpdateDocumentCount 更新文档数量
	UpdateDocumentCount(id string, count int64) error
}

// categoryRepository 类别仓储实现
type categoryRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewCategoryRepository 创建类别仓储
func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	return &categoryRepository{
		db:         db,
		collection: db.Collection("categories"),
	}
}

// Create 创建类别
func (r *categoryRepository) Create(category *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	// 默认状态为激活
	if category.Status == "" {
		category.Status = "active"
	}

	// 如果有父类别，需要更新父类别的children字段
	if category.ParentID != nil {
		// 插入类别
		result, err := r.collection.InsertOne(ctx, category)
		if err != nil {
			return err
		}
		category.ID = result.InsertedID.(primitive.ObjectID)

		// 更新父类别的children字段
		_, err = r.collection.UpdateOne(
			ctx,
			bson.M{"_id": *category.ParentID},
			bson.M{"$addToSet": bson.M{"children": category.ID}},
		)
		if err != nil {
			// 如果更新父类别失败，删除刚创建的类别
			r.collection.DeleteOne(ctx, bson.M{"_id": category.ID})
			return err
		}
	} else {
		// 插入根类别
		result, err := r.collection.InsertOne(ctx, category)
		if err != nil {
			return err
		}
		category.ID = result.InsertedID.(primitive.ObjectID)
	}

	return nil
}

// GetByID 通过ID获取类别
func (r *categoryRepository) GetByID(id string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category ID format")
	}

	var category models.Category
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

// GetByName 通过名称获取类别
func (r *categoryRepository) GetByName(name string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var category models.Category
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

// List 获取类别列表
func (r *categoryRepository) List(page, pageSize int) ([]*models.Category, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "level", Value: 1}, {Key: "sort_order", Value: 1}}) // 按级别和排序顺序

	// 查询类别
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// Update 更新类别
func (r *categoryRepository) Update(id string, data *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID format")
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
		return errors.New("category not found")
	}

	return nil
}

// Delete 删除类别
func (r *categoryRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID format")
	}

	// 检查是否有子类别
	childCount, err := r.collection.CountDocuments(ctx, bson.M{"parent_id": objectID})
	if err != nil {
		return err
	}

	if childCount > 0 {
		return errors.New("cannot delete category with children")
	}

	// 获取类别信息以便更新父类别
	var category models.Category
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("category not found")
		}
		return err
	}

	// 删除类别
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("category not found")
	}

	// 如果有父类别，从父类别的children中移除
	if category.ParentID != nil {
		_, err = r.collection.UpdateOne(
			ctx,
			bson.M{"_id": *category.ParentID},
			bson.M{"$pull": bson.M{"children": objectID}},
		)
		// 这里不返回错误，因为类别已经删除成功
	}

	return nil
}

// GetByParent 获取子类别
func (r *categoryRepository) GetByParent(parentID string) ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, errors.New("invalid parent ID format")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"parent_id": objectID}, options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetRootCategories 获取根类别
func (r *categoryRepository) GetRootCategories() ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"parent_id": nil}, options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryTree 获取类别树
func (r *categoryRepository) GetCategoryTree() ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "level", Value: 1}, {Key: "sort_order", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// UpdateDocumentCount 更新文档数量
func (r *categoryRepository) UpdateDocumentCount(id string, count int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"document_count": count,
			"updated_at":     time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("category not found")
	}

	return nil
}
