package example

import (
	"bytes"
	"net/http"

	"github.com/gofiber/fiber/v3"

	exampleviews "golang-starter-kit/resources/views/web/example"
)

type WebController struct {
	service *Service
}

func NewWebController(service *Service) *WebController {
	return &WebController{service: service}
}

func (ctrl *WebController) Index(c fiber.Ctx) error {
	ctx := c.Context()
	items, err := ctrl.service.List(ctx)
	if err != nil {
		return err
	}

	viewItems := make([]exampleviews.CrudItem, len(items))
	for i, item := range items {
		viewItems[i] = exampleviews.CrudItem{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
		}
	}

	buf := new(bytes.Buffer)
	component := exampleviews.Index(viewItems)
	if err := component.Render(ctx, buf); err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(buf.String())
}

func (ctrl *WebController) Create(c fiber.Ctx) error {
	ctx := c.Context()
	buf := new(bytes.Buffer)
	component := exampleviews.Create()
	if err := component.Render(ctx, buf); err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(buf.String())
}

func (ctrl *WebController) Store(c fiber.Ctx) error {
	var req CreateCrudRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid form input")
	}

	ctx := c.Context()
	if _, err := ctrl.service.Create(ctx, req); err != nil {
		return err
	}

	return c.Redirect().To("/example/cruds")
}
