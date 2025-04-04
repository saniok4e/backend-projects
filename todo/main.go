package main

import (
	"bufio"
	"fmt"
	"os"
	//"strings"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
}

var Tasks []Task
var NextID int = 1

func AddTask(title string) {
	task := Task{
		ID:        NextID,
		Title:     title,
		Completed: false,
	}
	Tasks = append(Tasks, task)
	fmt.Printf("✅ Задача[%d] \"%s\" успешно добавлена\n", task.ID, task.Title)
	NextID++
}

func CompleteTask(id int) {
	for i, task := range Tasks {
		if task.ID == id {
			Tasks[i].Completed = true
			fmt.Printf("✅ Задача [%d] помечена как выполненная\n", Tasks[i].ID)
			return
		}
	}
	fmt.Printf("❌ Задача с номером [%d] не найдена\n", id)
}

func ShowTasks(filter string) {
	switch filter {
	case "all":
		fmt.Println("\n--- Задачи ---")
		for _, val := range Tasks {
			status := "❌"
			if val.Completed {
				status = "✅"
			}
			fmt.Printf("%s Задача[%d] - %s\n", status, val.ID, val.Title)
		}
	case "done":
		fmt.Println("\n--- Выполненные задачи ---")
		for _, val := range Tasks {
			if val.Completed {
				fmt.Printf("✅ Задача[%d] - %s\n", val.ID, val.Title)
			}
		}
	case "pending":
		fmt.Println("\n--- Задачи в процессе ---")
		for _, val := range Tasks {
			if !val.Completed {
				fmt.Printf("❌ Задача[%d] - %s\n", val.ID, val.Title)
			}
		}
	default:
		fmt.Println("\n❌ Неизвестный фильтр")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n------ Моя тудушка ------")
		fmt.Println("Выбрать опцию:")
		fmt.Println("1. Добавить задачу")
		fmt.Println("2. Пометить как выполненную")
		fmt.Println("3. Показать задачи")
		fmt.Println("4. Выйти")
		fmt.Print("Выбор: ")

		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "1":
			fmt.Print("Введите задачу для выполнения (или напишите 'exit' для отмены): ")
			scanner.Scan()
			task := scanner.Text()
			if task != "exit" {
				AddTask(task)
			}
		case "2":
			fmt.Print("Введите айди задачи, которую пометить как выполненную: ")
			scanner.Scan()
			var id int
			fmt.Sscanf(scanner.Text(), "%d", &id)
			CompleteTask(id)
		case "3":
			fmt.Print("Введите фильтр ('pending', 'done', 'all'): ")
			scanner.Scan()
			filter := scanner.Text()
			ShowTasks(filter)
		case "4":
			fmt.Println("👋 Выходим...")
			return
		default:
			fmt.Println("❌ Неправильный ввод")
		}
	}
}
