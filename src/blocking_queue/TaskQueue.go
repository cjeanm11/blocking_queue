package blocking_queue

import (
	"time"
)

// TODO add priority

type Task struct {
	Typ       byte        // Typ is the type identifier of the task.
	StartTime time.Time   // StartTime is the timestamp when the task started.
	EndTime   time.Time   // EndTime is the timestamp when the task ended.
	ID        int         // ID is an identifier for the task. (TODO: Provide a description)
	Data      interface{} // Data is the payload or additional data associated with the task.
	Status    string      // Status is the current status of the task. (TODO: Provide a description)
	Result    interface{} // Result stores the output or result of the task.
	Error     error       // Error holds information about any errors that occurred during task execution. (TODO: Handle error reporting)
	Owner     string      // Owner is the owner or user responsible for the task. (TODO: Add func or method to set owner)
	CreatedAt time.Time   // CreatedAt is the timestamp when the task was created. (TODO: improve)
	UpdatedAt time.Time   // UpdatedAt is the timestamp when the task was last updated. (TODO: improve)
}

type TaskQueue []*Task
