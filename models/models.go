package models
 
import (
        "time"
        "gorm.io/gorm"
)
 
type Task struct {
        gorm.Model
        ID        uint
        TaskName  string `validate:"omitempty,ascii"`
        Assignee  string
        CreatedAt time.Time
        IsDone    bool `gorm:"default:false" json:"isDone"`
        UserID    uint
}
type TaskResponse struct {
        ID        uint
        TaskName  string `validate:"omitempty,ascii"`
        Assignee  string
        CreatedAt time.Time
        IsDone    bool `gorm:"default:false" json:"isDone"`
        UserID    uint
}
 
type User struct {
        gorm.Model
        ID        uint   `json:"id"`
        Username  string `json:"username" validate:"omitempty,min=5,max=16,alphanum"`
        Password  string `json:"password" validate:"omitempty,min=8,max=20,alphanum"`
        Tasks     []Task `json:"tasks"`
}
type UserResponse struct {
        ID       uint   `json:"id"`
        Username string `json:"username"`
        Tasks    []Task `json:"tasks" gorm:"foreignKey:TaskName"`
}
