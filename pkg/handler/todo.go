package handler

import (
	"sachahjkl/htmx_go/pkg/config"
	"sachahjkl/htmx_go/pkg/middleware"
	"sachahjkl/htmx_go/pkg/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterTodoRoutes(app *fiber.App, db *gorm.DB, c *config.Config) {
	h := &handler{
		DB:     db,
		Config: c,
	}

	middleware := middleware.New(db)

	// Protect all "todo" routes (meanning innacessible to unlogged users)
	routes := app.Group("/todos", middleware.Protected)

	routes.Get("/", h.GetTodos).Name("todos")
	routes.Post("/", h.AddTodo)
	routes.Put("/:id/toggle", h.ToggleTodo)
	routes.Delete("/:id", h.DeleteTodo)
}

type AddTodoRequestBody struct {
	Title string `json:"title" form:"todo-title"`
}

func (h handler) GetTodos(c *fiber.Ctx) error {
	var todos []model.Todo

	// get all the todos
	result := h.DB.Model(&model.Todo{}).Find(&todos)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return c.Render("todos/index", fiber.Map{
		"Todos": todos,
	}, "layouts/main")
}

func (h handler) AddTodo(c *fiber.Ctx) error {

	body := AddTodoRequestBody{}

	// parse body, unmarshall to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todo, err := model.AddTodo(h.DB, body.Title, false, 0)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Render("todos/list_item", todo)
}

func (h handler) ToggleTodo(c *fiber.Ctx) error {

	param := struct {
		ID uint `params:"id"`
	}{}

	err := c.ParamsParser(&param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todo, err := model.ToggleTodo(h.DB, param.ID, 0)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Render("todos/list_item", todo)

}

func (h handler) DeleteTodo(c *fiber.Ctx) error {

	param := struct {
		ID uint `params:"id"`
	}{}

	err := c.ParamsParser(&param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = model.DeleteTodo(h.DB, param.ID, 0)

	if err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(nil)

}
