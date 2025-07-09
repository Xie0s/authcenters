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

// TagRepository 标签数据访问接口
type TagRepository interface {
	// Create 创建标签
	Create(tag *models.Tag) error

	// GetByID 通过ID获取标签
	GetByID(id string) (*models.Tag, error)

	// GetByName 通过名称获取标签
	GetByName(name string) (*models.Tag, error)

	// List 获取标签列表
	List(page, pageSize int) ([]*models.Tag, int64, error)

	// Update 更新标签
	Update(id string, data *models.Tag) error

	// Delete 删除标签
	Delete(id string) error

	// Search 搜索标签
	Search(keyword string, page, pageSize int) ([]*models.Tag, int64, error)

	// GetByCategory 按类别获取标签
	GetByCategory(categoryID string) ([]*models.Tag, error)

	// GetPopularTags 获取热门标签
	GetPopularTags(limit int) ([]*models.Tag, error)

	// UpdateUsageCount 更新使用次数
	UpdateUsageCount(id string, increment int64) error

	// GetByNames 通过名称列表获取标签
	GetByNames(names []string) ([]*models.Tag, error)
}

// tagRepository 标签仓储实现
type tagRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewTagRepository 创建标签仓储
func NewTagRepository(db *mongo.Database) TagRepository {
	return &tagRepository{
		db:         db,
		collection: db.Collection("tags"),
	}
}

// Create 创建标签
func (r *tagRepository) Create(tag *models.Tag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置创建时间
	now := time.Now()
	tag.CreatedAt = now
	tag.LastUsedAt = now

	// 默认使用次数为0
	if tag.UsageCount == 0 {
		tag.UsageCount = 0
	}

	// 插入标签
	result, err := r.collection.InsertOne(ctx, tag)
	if err != nil {
		return err
	}

	// 设置生成的ID
	tag.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID 通过ID获取标签
func (r *tagRepository) GetByID(id string) (*models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid tag ID format")
	}

	var tag models.Tag
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return &tag, nil
}

// GetByName 通过名称获取标签
func (r *tagRepository) GetByName(name string) (*models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tag models.Tag
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}

	return &tag, nil
}

// List 获取标签列表
func (r *tagRepository) List(page, pageSize int) ([]*models.Tag, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "usage_count", Value: -1}, {Key: "created_at", Value: -1}}) // 按使用次数和创建时间排序

	// 查询标签
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tags []*models.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// Update 更新标签
func (r *tagRepository) Update(id string, data *models.Tag) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid tag ID format")
	}

	// 创建更新文档
	updateDoc := bson.M{"$set": data}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("tag not found")
	}

	return nil
}

// Delete 删除标签
func (r *tagRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid tag ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("tag not found")
	}

	return nil
}

// Search 搜索标签
func (r *tagRepository) Search(keyword string, page, pageSize int) ([]*models.Tag, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 构建搜索条件
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	// 设置查询选项
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{Key: "usage_count", Value: -1}}) // 按使用次数排序

	// 查询标签
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var tags []*models.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// GetByCategory 按类别获取标签
func (r *tagRepository) GetByCategory(categoryID string) ([]*models.Tag, error) {
	// 这里需要通过知识文档来关联标签和分类
	// 暂时返回空结果，实际实现需要根据业务逻辑调整
	return []*models.Tag{}, nil
}

// GetPopularTags 获取热门标签
func (r *tagRepository) GetPopularTags(limit int) ([]*models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "usage_count", Value: -1}}) // 按使用次数倒序

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []*models.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// UpdateUsageCount 更新使用次数
func (r *tagRepository) UpdateUsageCount(id string, increment int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid tag ID format")
	}

	update := bson.M{
		"$inc": bson.M{"usage_count": increment},
		"$set": bson.M{"last_used_at": time.Now()},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("tag not found")
	}

	return nil
}

// GetByNames 通过名称列表获取标签
func (r *tagRepository) GetByNames(names []string) ([]*models.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"name": bson.M{"$in": names}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []*models.Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}
