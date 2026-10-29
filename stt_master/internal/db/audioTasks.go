package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

const (
	TABLE_AUDIO_TASKS            = "audio_tasks"
	AUDIO_TASKS_SATUS_PENDING    = "pending"
	AUDIO_TASKS_SATUS_PROCESSING = "processing"
	AUDIO_TASKS_SATUS_COMPLETED  = "completed"
	AUDIO_TASKS_SATUS_FAILED     = "failed"
)

var validStatuses = map[string]bool{
	AUDIO_TASKS_SATUS_PENDING:    true,
	AUDIO_TASKS_SATUS_PROCESSING: true,
	AUDIO_TASKS_SATUS_COMPLETED:  true,
	AUDIO_TASKS_SATUS_FAILED:     true,
}

type AudioTask struct {
	ID           int64          `json:"id"`
	FilePath     string         `json:"file_path"`
	Status       string         `json:"status"`
	ErrorMessage sql.NullString `json:"error_message"`
	UpdateAt     time.Time      `json:"updated_at"`
}

func (t AudioTask) MarshalJSON() ([]byte, error) {
	type Alias AudioTask
	var errStr *string

	if t.ErrorMessage.Valid {
		errStr = &t.ErrorMessage.String
	}

	return json.Marshal(&struct {
		Alias
		ErrorMessage *string `json:"error_message"`
	}{
		Alias:        Alias(t),
		ErrorMessage: errStr,
	})
}

func createAudioTasksTable(db *sql.DB) error {
	createAudiotracksTableSQL := fmt.Sprintf(`

	create table if not exists %s (
		id integer primary key autoincrement,
		file_path text not null unique,
		status text not null default '%s',
		error_message text,
		update_at datetime default current_timestamp
	);

	create trigger if not exists audio_tracks_update_at 
	after update on %s
	begin
		update %s set update_at = current_timestamp where id = old.id;
	end;
	`, TABLE_AUDIO_TASKS, AUDIO_TASKS_SATUS_PENDING, TABLE_AUDIO_TASKS, TABLE_AUDIO_TASKS)

	if _, err := db.Exec(createAudiotracksTableSQL); err != nil {
		fmt.Printf("Ошибка создания таблицы: %v\n", err)
		return err
	}

	return nil
}

func dropAudioTasksTable(db *sql.DB) error {
	dropSQL := fmt.Sprintf(`
		drop table if exists %s;

		drop trigger if exists audio_tracks_update_at;
	`, TABLE_AUDIO_TASKS)

	_, err := db.Exec(dropSQL)

	return err
}

func RecreateAudioTasksTable(db *sql.DB) error {
	err := dropAudioTasksTable(db)

	if err != nil {
		return err
	}

	return createAudioTasksTable(db)
}

func GetTaskById(db *sql.DB, taskId int64) (*AudioTask, error) {
	var task AudioTask

	query := fmt.Sprintf("select id, file_path, status, error_message, update_at from %s where id = ?", TABLE_AUDIO_TASKS)

	err := db.QueryRow(query, taskId).Scan(&task.ID, &task.FilePath, &task.Status, &task.ErrorMessage, &task.UpdateAt)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func GetAllTasks(db *sql.DB) ([]AudioTask, error) {
	query := fmt.Sprintf("select id, file_path, status, error_message, update_at from %s order by id asc", TABLE_AUDIO_TASKS)

	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("Error on 'GetAllTasks' %w", err)
	}

	defer rows.Close()

	tasks := make([]AudioTask, 0)

	for rows.Next() {
		var task AudioTask

		err := rows.Scan(&task.ID, &task.FilePath, &task.Status, &task.ErrorMessage, &task.UpdateAt)

		if err != nil {
			return nil, fmt.Errorf("Error on 'GetAllTasks' %w", err)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on 'GetAllTasks' %w", err)
	}

	return tasks, nil
}

func validateStatus(status string) error {
	if !validStatuses[status] {
		return fmt.Errorf("Invalid tasks status '%s'", status)
	}

	return nil
}

func AddTask(db *sql.DB, filepath string, status string) error {
	if status == "" {
		status = AUDIO_TASKS_SATUS_PENDING
	}

	if err := validateStatus(status); err != nil {
		return err
	}

	if filepath == "" {
		return fmt.Errorf("File path is empty")
	}

	query := fmt.Sprintf("insert into %s (file_path, status) values (?, ?)", TABLE_AUDIO_TASKS)

	_, err := db.Exec(query, filepath, status)

	if err != nil {
		return fmt.Errorf("Error on 'AddTask' %w", err)
	}

	return nil
}

func updateTask(db *sql.DB, taskId int64, status string, errorMessage string) error {
	if err := validateStatus(status); err != nil {
		return err
	}

	if status == AUDIO_TASKS_SATUS_FAILED && errorMessage == "" {
		return fmt.Errorf("ErrorMessage is empty when status '%s'", status)
	}

	query := fmt.Sprintf("update %s set status = ?, error_message = ? where id = ?", TABLE_AUDIO_TASKS)

	var errStr sql.NullString

	if status != AUDIO_TASKS_SATUS_FAILED {
		errorMessage = ""
	}

	if errorMessage != "" {
		errStr = sql.NullString{String: errorMessage, Valid: true}
	}

	_, err := db.Exec(query, status, errStr, taskId)

	if err != nil {
		return fmt.Errorf("Error on 'UpdateTable' %w", err)
	}

	return nil
}

func MarkPendingTaskById(db *sql.DB, taskId int64) error {
	return updateTask(db, taskId, AUDIO_TASKS_SATUS_PENDING, "")
}

func MarkProcessedTaskById(db *sql.DB, taskId int64) error {
	return updateTask(db, taskId, AUDIO_TASKS_SATUS_PROCESSING, "")
}

func MarkCompletedTaskById(db *sql.DB, taskId int64) error {
	return updateTask(db, taskId, AUDIO_TASKS_SATUS_COMPLETED, "")
}

func MarkFailedTaskById(db *sql.DB, taskId int64, errorMessage string) error {
	return updateTask(db, taskId, AUDIO_TASKS_SATUS_FAILED, errorMessage)
}

func ClaimNextTask(db *sql.DB) (*AudioTask, error) {
	var task AudioTask

	query := fmt.Sprintf(`
		update %s
		set status = '%s'
		where id = (
			select id from %s
			where status = '%s'
			order by id asc
			limit 1
		)
		returning id, file_path, status, error_message, update_at;
	`, TABLE_AUDIO_TASKS, AUDIO_TASKS_SATUS_PROCESSING, TABLE_AUDIO_TASKS, AUDIO_TASKS_SATUS_PENDING)

	err := db.QueryRow(query).Scan(&task.ID, &task.FilePath, &task.Status, &task.ErrorMessage, &task.UpdateAt)

	if err != nil {
		return nil, fmt.Errorf("Error on 'ClaimNextTask' %w", err)
	}

	return &task, nil
}
