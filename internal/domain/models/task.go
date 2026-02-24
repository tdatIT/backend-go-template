package models

import "time"

// TaskGroup represents a collection of tasks in a todo list.
type TaskGroup struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint64    `json:"user_id" gorm:"index"`
	Icon        string    `json:"icon,omitempty" gorm:"size:100"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy   uint64    `json:"created_by" gorm:"index"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	//Relationships
	User  *User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Tasks []*Task `json:"tasks,omitempty" gorm:"foreignKey:GroupID;references:ID"`
}

func (TaskGroup) TableName() string {
	return "task_groups"
}

// Task represents a single todo item.
type Task struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"size:200;not null"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	Status      string     `json:"status" gorm:"size:20;not null;default:pending"`
	Priority    int        `json:"priority" gorm:"not null;default:0"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Order       int        `json:"order" gorm:"not null;default:0"`
	GroupID     uint64     `json:"group_id" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy   uint64     `json:"created_by" gorm:"index"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	//Relationships
	Group *TaskGroup `json:"group,omitempty" gorm:"foreignKey:GroupID;references:ID"`
}

func (Task) TableName() string {
	return "tasks"
}
