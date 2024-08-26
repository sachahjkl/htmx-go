package handler

import (
	"encoding/json"
	"sachahjkl/htmx_go/internal/common"
	"sachahjkl/htmx_go/internal/config"
	"sachahjkl/htmx_go/internal/middleware"
	"sachahjkl/htmx_go/internal/model"
	"strconv"
	"strings"
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
	routes.Put("/:id<int>/toggle", h.ToggleTodo)
	routes.Get("/export-json", h.ExportTodosJson)
	routes.Get("/export-csv", h.ExportTodosCsv)
	routes.Delete("/all", h.RemoveAllTodos)
	routes.Delete("/all-finished", h.RemoveAllFinishedTodos)
	routes.Delete("/:id<int>", h.DeleteTodo)
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

	bind := fiber.Map{
		"Todos":    todos,
		"Username": user.Username,
		"Config":   h.Config,
	}

	layout := ""

	// Called from HTMX, no layout
	if c.Get("HX-Request") == "" {
		layout = "layouts/main"
	}

	return c.Render("todos/index", bind, layout)
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

func (h handler) ExportTodosCsv(c *fiber.Ctx) error {

	user := c.Locals(common.USER_KEY).(*model.User)

	if user == nil {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	// get all the todos
	todos, err := model.AllTodos(h.DB, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	const separator = ";"

	// map todos to csv
	csvTodos := make([]string, len(*todos))
	columnNames := []string{"Done", "Title", "Created"}
	for i, todo := range *todos {
		csvTodos[i] = strconv.FormatBool(todo.Done) + separator + todo.Title + separator + todo.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// construct filename with username and date from date library
	filename := "export_" + user.Username + "_" + time.Now().Format("2006_01_02_15_04_05") + ".csv"
	c.Attachment(filename)

	// set content type
	c.Type("text/csv")

	return c.SendString(strings.Join(columnNames, separator) + "\n" + strings.Join(csvTodos, "\n"))
}

func (h handler) RemoveAllTodos (c *fiber.Ctx) error { 
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h handler) RemoveAllFinishedTodos (c *fiber.Ctx) error { 
	return c.SendStatus(fiber.StatusNotImplemented)
}
