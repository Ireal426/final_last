package api

import (
	"encoding/json"
	"net/http"
	"log"

	"go_final_project/pkg/db"
)

const TasksLimit = 50

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(TasksLimit)
	if err != nil {
		sendError(w, "Error getting task list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(TasksResp{Tasks: tasks}); err != nil {
        log.Printf("JSON encoding error: %v", err)
        return
    }
}