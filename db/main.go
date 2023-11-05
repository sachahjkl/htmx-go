package db

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var idCount = 0

func get_id() int {

	currID := idCount
	idCount += 1

	return currID
}

var Todos = []Todo{
	{
		Id:    get_id(),
		Title: "Clean",
		Done:  true,
	},
	{
		Id:    get_id(),
		Title: "Cook",
		Done:  false,
	},
}

func AddTodo(title string, done bool) {
	Todos = append(Todos, Todo{
		Id:    get_id(),
		Title: title,
		Done:  done,
	})
}

func ToggleTodo(id int) {
	for i := range Todos {
		if Todos[i].Id == id {
			Todos[i].Done = !Todos[i].Done
			return
		}
	}
}
