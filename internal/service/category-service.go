package service

import (
	"github.com/Fajar3108/mafi-course-be/database"
	"github.com/Fajar3108/mafi-course-be/internal/model"
	"github.com/Fajar3108/mafi-course-be/internal/request"
	"github.com/Fajar3108/mafi-course-be/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		db: database.DB(),
	}
}

func (cs *CategoryService) GetAll(page, limit int) (categories *[]model.Category, err error) {
	offset := (page - 1) * limit

	result := cs.db.Order("name ASC").Offset(offset).Limit(limit).Find(&categories)

	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return categories, nil
}

func (cs *CategoryService) GetBySlug(slug string) (category *model.Category, err error) {
	result := cs.db.First(&category, "slug = ?", slug)

	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	if category.ID == "" {
		return nil, fiber.NewError(fiber.StatusNotFound, "Category not found")
	}

	return category, nil
}

func (cs *CategoryService) Store(categoryRequest *request.CategoryRequest) (category *model.Category, err error) {
	id, err := helpers.GenerateUUID()

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	category = &model.Category{
		ID:   id,
		Name: categoryRequest.Name,
		Slug: helpers.Slug(categoryRequest.Name),
	}

	if result := cs.db.Create(category); result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return category, nil
}

func (cs *CategoryService) Update(categoryRequest *request.CategoryRequest, slug string) (category *model.Category, err error) {
	category, err = cs.GetBySlug(slug)

	if err != nil {
		return nil, err
	}

	category.Name = categoryRequest.Name
	category.Slug = helpers.Slug(categoryRequest.Name)

	if result := cs.db.Save(category); result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return category, nil
}

func (cs *CategoryService) Destroy(slug string) error {
	category, err := cs.GetBySlug(slug)

	if err != nil {
		return err
	}

	if result := cs.db.Delete(category); result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return nil
}
