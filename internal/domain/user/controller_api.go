package user

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"golang-starter-kit/pkg/logger"
	"golang-starter-kit/pkg/response"
	"golang-starter-kit/pkg/validator"
)

type APIController struct {
	service *Service
}

func NewAPIController(service *Service) *APIController {
	return &APIController{service: service}
}

func (ctrl *APIController) List(c fiber.Ctx) error {
	ctx := c.Context()
	items, err := ctrl.service.List(ctx)
	if err != nil {
		return err
	}

	traceID := logger.CorrelationIDFromContext(ctx)
	meta := response.NewMeta(traceID)
	meta.Total = int64(len(items))

	return c.JSON(response.Success("User list retrieved", items, meta))
}

func (ctrl *APIController) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.Context()

	item, err := ctrl.service.GetByID(ctx, id)
	if err != nil {
		return err
	}

	traceID := logger.CorrelationIDFromContext(ctx)
	return c.JSON(response.Success("User detail retrieved", item, response.NewMeta(traceID)))
}

func (ctrl *APIController) Create(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error("Invalid request body", nil))
	}

	if valErrors := validator.ValidateStruct(req); len(valErrors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(response.ValidationFailed(valErrors))
	}

	ctx := c.Context()
	item, err := ctrl.service.Create(ctx, req)
	if err != nil {
		return err
	}

	traceID := logger.CorrelationIDFromContext(ctx)
	return c.Status(http.StatusCreated).JSON(response.Success("User created successfully", item, response.NewMeta(traceID)))
}

func (ctrl *APIController) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Error("Invalid request body", nil))
	}

	if valErrors := validator.ValidateStruct(req); len(valErrors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(response.ValidationFailed(valErrors))
	}

	ctx := c.Context()
	item, err := ctrl.service.Update(ctx, id, req)
	if err != nil {
		return err
	}

	traceID := logger.CorrelationIDFromContext(ctx)
	return c.JSON(response.Success("User updated successfully", item, response.NewMeta(traceID)))
}

func (ctrl *APIController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.Context()

	if err := ctrl.service.Delete(ctx, id); err != nil {
		return err
	}

	traceID := logger.CorrelationIDFromContext(ctx)
	return c.JSON(response.Success("User deleted successfully", nil, response.NewMeta(traceID)))
}
