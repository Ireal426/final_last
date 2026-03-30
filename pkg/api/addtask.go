package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "Ошибка десериализации JSON")
		return
	}

	if task.Title == "" {
		sendError(w, "Не указан заголовок задачи")
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	
	if task.Date == "" {
		task.Date = now.Format(TimeLayout)
	}

	t, err := time.Parse(TimeLayout, task.Date)
	if err != nil {
		sendError(w, "Дата представлена в неверном формате")
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(TimeLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			task.Date = next
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		sendError(w, "Ошибка при добавлении в базу данных")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "Задача не найдена")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "Ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		sendError(w, "Не указан идентификатор задачи")
		return
	}
	if task.Title == "" {
		sendError(w, "Не указан заголовок задачи")
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	t, err := time.Parse(TimeLayout, task.Date)
	if err != nil {
		sendError(w, "Дата представлена в неверном формате")
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(TimeLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			task.Date = next
		}
	}

	err = db.UpdateTask(task)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`))
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`))
}