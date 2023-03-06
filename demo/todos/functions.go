package todos

import (
	"encoding/json"
	"io/ioutil"
)

// ToDoPageData data structure
type ToDoPageData struct {
	PageTitle string
	ToDos     []ToDo
}

// ToDo data structure
type ToDo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Note  string `json:"note"`
	Done  bool   `json:"done"`
}

func (t ToDo) ToString() string {

	bytes, _ := json.Marshal(t)
	return string(bytes)

}

func GetToDos() []ToDo {

	todos := make([]ToDo, 3)

	raw, _ := ioutil.ReadFile("./todos.json")

	json.Unmarshal(raw, &todos)

	return todos

}
