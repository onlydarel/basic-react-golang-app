package api

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/onlydarel/basic-react-golang-app/internal/models"
	"log"
)

var db *sql.DB

func SetDatabase(database *sql.DB) {
	db = database
}

func GetTodos(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, status, body FROM todos")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query todos"})
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Status, &todo.Body); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan todos"})
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

// Create a Todo
func AddTodos(c *fiber.Ctx) error {

	todo := &models.Todo{}

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	_, err := db.Exec("insert into todos (status, body) values ($1, $2)", false, todo.Body)
	if err != nil {
		log.Fatal("Error inserting todo", err)
	}

	return c.Status(201).JSON(todo)

}

func UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println(id)
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo id is required"})
	}

	// Fetch the current status of the todo with the specified id
	var todo models.Todo
	err := db.QueryRow("SELECT id, status FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query todos"})
	}

	// Toggle the completion status
	newStatus := !todo.Status
	_, err = db.Exec("UPDATE todos SET status = $1 WHERE id = $2", newStatus, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
	}

	// Update the todo struct and respond with it
	todo.Status = newStatus
	return c.Status(201).JSON(todo)
}

func DeleteTodo(c *fiber.Ctx) error {
	delete_id := c.Params("id")
	if delete_id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo id is required"})
	}

	var todo models.Todo
	err := db.QueryRow("SELECT id FROM todos WHERE id = $1", delete_id).Scan(&todo.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Todo not found, enter the correct id!"})
	}

	_, err = db.Exec("delete from todos where id = $1", delete_id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
	}

	return c.Status(200).JSON(fiber.Map{"msg": "successfully deleted todo"})
}
