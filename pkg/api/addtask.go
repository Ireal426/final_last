package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "JSON deserialization error", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendError(w, "Task title not specified", http.StatusBadRequest)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	
	if task.Date == "" {
		task.Date = now.Format(TimeLayout)
	}

	t, err := time.Parse(TimeLayout, task.Date)
	if err != nil {
		sendError(w, "The date is in an invalid format", http.StatusBadRequest)
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(TimeLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = next
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		sendError(w, "Error adding to database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    
    resp := map[string]interface{}{"id": strconv.FormatInt(id, 10)}
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("error encoding addtask response: %v", err)
    }
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "No identifier specified", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "JSON deserialization error", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		sendError(w, "Task ID not specified", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		sendError(w, "Task title not specified", http.StatusBadRequest)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	t, err := time.Parse(TimeLayout, task.Date)
	if err != nil {
		sendError(w, "The date is in an invalid format", http.StatusBadRequest)
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(TimeLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = next
		}
	}

	err = db.UpdateTask(task)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`))
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "No identifier specified", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`))
}