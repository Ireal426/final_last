package api

import (
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		now := time.Now()
		next, errNext := NextDate(now, task.Date, task.Repeat)
		if errNext != nil {
			sendError(w, errNext.Error(), http.StatusBadRequest)
			return
		}
		err = db.UpdateTaskDate(id, next)
	}

	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		log.Printf("error writing done response: %v", err)
	}
}