package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TODO struct {
	Title string
	Done  bool
}

type TODOS []*TODO

func (todos *TODOS) AddItem(title string) {
	for _, todo := range *todos {
		if todo.Title == title {
			return
		}
	}
	*todos = append(*todos, &TODO{Title: title})
}

func (todos TODOS) MarkDone(title string) {
	for _, todo := range todos {
		if todo.Title == title {
			todo.Done = true
			return
		}
	}
}

func (todos *TODOS) RemoveItem(title string) {
	pos := -1
	for i, todo := range *todos {
		if todo.Title == title {
			pos = i
			break
		}
	}
	if pos == -1 {
		return
	}
	*todos = append((*todos)[:pos], (*todos)[pos+1:]...)
}

func (todos TODOS) GetTODOList() []*TODO {
	return todos
}

func main() {
	todos := TODOS{}
	r := gin.Default()
	r.GET("todos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, todos.GetTODOList())
	})
	r.POST("todo", func(ctx *gin.Context) {
		var todo TODO
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			log.Println(err)
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		todos.AddItem(todo.Title)
		ctx.String(http.StatusCreated, "The new TODO has been added")
	})
	r.PATCH("todo", func(ctx *gin.Context) {
		var todo TODO
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		todos.MarkDone(todo.Title)
		ctx.String(http.StatusOK, "Your TODO has been marked as done")
	})
	r.DELETE("todo", func(ctx *gin.Context) {
		var todo TODO
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		todos.RemoveItem(todo.Title)
		ctx.String(http.StatusOK, "Your TODO is deleted successfully")
	})
	r.Run()
}
