package controller

import (
	"github.com/Fajar3108/mafi-course-be/internal/request"
	"github.com/Fajar3108/mafi-course-be/internal/service"
	"github.com/Fajar3108/mafi-course-be/pkg/helpers"
	"github.com/Fajar3108/mafi-course-be/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	service *service.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		service: service.NewCategoryService(),
	}
}

func (cc *CategoryController) Index(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	posts, err := cc.service.GetAll(page, limit)

	if err != nil {
		return err
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Categories retrieved successfully",
		posts,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (cc *CategoryController) Show(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	category, err := cc.service.GetBySlug(slug)

	if err != nil {
		return err
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Category retrieved successfully",
		category,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (cc *CategoryController) Store(ctx *fiber.Ctx) error {
	req := &request.CategoryRequest{}

	if err := validation.Validate(ctx, req); err != nil {
		return err
	}

	category, err := cc.service.Store(req)

	if err != nil {
		return err
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Category created successfully",
		category,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (cc *CategoryController) Update(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	req := &request.CategoryRequest{}

	if err := validation.Validate(ctx, req); err != nil {
		return err
	}

	category, err := cc.service.Update(req, slug)

	if err != nil {
		return err
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Category updated successfully",
		category,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (cc *CategoryController) Destroy(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	err := cc.service.Destroy(slug)

	if err != nil {
		return err
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Category deleted successfully",
		nil,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}
