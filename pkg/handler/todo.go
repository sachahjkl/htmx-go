package handler

import (
	"encoding/json"
	"sachahjkl/htmx_go/pkg/common"
	"sachahjkl/htmx_go/pkg/config"
	"sachahjkl/htmx_go/pkg/middleware"
	"sachahjkl/htmx_go/pkg/model"
	"time"

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
	routes := app.Group("/todos", middleware.AuthenticatedOnly)

	routes.Get("/", h.GetTodos).Name("todos")
	routes.Post("/", h.AddTodo)
	routes.Put("/:id/toggle", h.ToggleTodo)
	routes.Delete("/:id", h.DeleteTodo)
	routes.Get("/export", h.ExportTodosJson)
	routes.Get("/settings", h.Settings)
}

type AddTodoRequestBody struct {
	Title string `json:"title" form:"todo-title"`
}

func (h handler) GetTodos(c *fiber.Ctx) error {
	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	// get all the todos
	todos, err := model.AllTodos(h.DB, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Called from HTMX
	if c.Get("HX-Request") != "" {
		return c.Render("todos/index", fiber.Map{
			"Todos":    todos,
			"Username": user.Username,
			"Config":   h.Config,
		})

	}

	// Called from the browser
	return c.Render("todos/index", fiber.Map{
		"Todos":    todos,
		"Username": user.Username,
		"Config":   h.Config,
	}, "layouts/main")
}

func (h handler) AddTodo(c *fiber.Ctx) error {
	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	body := AddTodoRequestBody{}

	// parse body, unmarshall to AddTodoRequestBody struct
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todo, err := model.AddTodo(h.DB, body.Title, false, user.ID)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Render("todos/list_item", todo)
}

func (h handler) ToggleTodo(c *fiber.Ctx) error {
	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	param := struct {
		ID uint `params:"id"`
	}{}

	err := c.ParamsParser(&param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todo, err := model.ToggleTodo(h.DB, param.ID, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Render("todos/list_item", todo)

}

func (h handler) DeleteTodo(c *fiber.Ctx) error {

	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	param := struct {
		ID uint `params:"id"`
	}{}

	err := c.ParamsParser(&param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = model.DeleteTodo(h.DB, param.ID, user.ID)

	if err != nil {
		fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(nil)

}

func (h handler) ExportTodosJson(c *fiber.Ctx) error {

	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	// get all the todos
	todos, err := model.AllTodos(h.DB, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// map todos to json
	jsonTodos := make([]model.TodoJSON, len(*todos))
	for i, todo := range *todos {
		jsonTodos[i] = model.TodoJSON{
			Title:   todo.Title,
			Done:    todo.Done,
			Created: todo.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// construct filename with username and date from date library
	filename := "export_" + user.Username + "_" + time.Now().Format("2006_01_02_15_04_05") + ".json"
	c.Attachment(filename)

	// set content type
	c.Type("application/json")

	json, err := json.MarshalIndent(jsonTodos, "", "  ")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(json)

}

func (h handler) Settings(c *fiber.Ctx) error {
	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	// Called from HTMX
	if c.Get("HX-Request") != "" {
		return c.Render("todos/settings", fiber.Map{
			"Username": user.Username,
			"Config":   h.Config,
		})
	}

	// Called from the browser
	return c.Render("todos/settings", fiber.Map{
		"Username": user.Username,
		"Config":   h.Config,
	}, "layouts/main")
}
