package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-final-project/internal/db"
	"go-final-project/internal/helpers"
	"go-final-project/internal/schedule"
	"net/http"
	"strings"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func writeDBError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrTaskNotFound) {
		helpers.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.FormValue("search"))
	tasks, err := db.Tasks(50, search)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helpers.WriteJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeDBError(w, err)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeDBError(w, err)
			return
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]string{})
		return
	}

	next, err := schedule.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeDBError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]string{})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeDBError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, task)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	task.Title = strings.TrimSpace(task.Title)
	task.Date = strings.TrimSpace(task.Date)
	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.Title == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	if err := checkDate(&task); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprint(id)})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	task.ID = strings.TrimSpace(task.ID)
	task.Date = strings.TrimSpace(task.Date)
	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.ID == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if task.Title == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	if err := checkDate(&task); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeDBError(w, err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeDBError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]string{})
}

func checkDate(task *db.Task) error {
	now := schedule.StartOfDay(time.Now())

	if task.Date == "" {
		task.Date = now.Format(schedule.DateLayout)
	}

	t, err := time.Parse(schedule.DateLayout, task.Date)
	if err != nil {
		return err
	}
	t = schedule.StartOfDay(t)

	next := ""
	if task.Repeat != "" {
		next, err = schedule.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	// Если task.Date < today:
	// - без repeat -> today
	// - с repeat -> next
	if schedule.AfterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(schedule.DateLayout)
		} else {
			task.Date = next
		}
	}
	return nil
}
