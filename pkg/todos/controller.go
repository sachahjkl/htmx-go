package todos

import (
	"sachahjkl/htmx_go/pkg/common/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := app.Group("/todos")
	routes.Get("/", h.GetTodos)
	routes.Post("/", h.AddTodo)
	routes.Put("/:id/toggle", h.ToggleTodo)
	routes.Delete("/:id", h.DeleteTodo)
}

type AddTodoRequestBody struct {
	Title string `json:"title" form:"todo-title"`
}

func (h handler) GetTodos(c *fiber.Ctx) error {
	var todos []models.Todo

	// get all the todos
	result := h.DB.Model(&models.Todo{}).Find(&todos)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return c.Render("todos/index", fiber.Map{
		"Todos": todos,
	}, "layouts/main")
}

func (h handler) AddTodo(c *fiber.Ctx) error {

	body := AddTodoRequestBody{}

	// parse body, attach to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var todo models.Todo

	todo.Title = body.Title
	todo.Done = false

	// insert new db entry
	if result := h.DB.Create(&todo); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	var todos []models.Todo

	// get all the todos
	result := h.DB.Model(&models.Todo{}).Find(&todos)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return c.Render("todos/list", fiber.Map{
		"Todos": todos,
	})
}

func (h handler) ToggleTodo(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var todo models.Todo
	result := h.DB.First(&todo, id)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	// toggle the todo
	todo.Done = !todo.Done

	result = h.DB.Save(&todo)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	var todos []models.Todo

	// get all the todos
	result = h.DB.Model(&models.Todo{}).Find(&todos)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return c.Render("todos/list", fiber.Map{
		"Todos": todos,
	})

}

func (h handler) DeleteTodo(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var todo models.Todo
	result := h.DB.First(&todo, id)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	result = h.DB.Delete(&todo)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	var todos []models.Todo

	// get all the todos
	result = h.DB.Model(&models.Todo{}).Find(&todos)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return c.Render("todos/list", fiber.Map{
		"Todos": todos,
	})

}
