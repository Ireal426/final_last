package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		sendError(w, "Ошибка при получении списка задач")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(TasksResp{Tasks: tasks}); err != nil {
		sendError(w, "Ошибка кодирования JSON")
		return
	}
}