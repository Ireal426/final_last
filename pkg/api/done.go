package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		now := time.Now()
		next, errNext := NextDate(now, task.Date, task.Repeat)
		if errNext != nil {
			sendError(w, errNext.Error())
			return
		}
		err = db.UpdateTaskDate(id, next)
	}

	if err != nil {
		sendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`))
}