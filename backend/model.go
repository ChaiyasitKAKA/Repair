package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Role     string `json:"role" gorm:"default:user"` // "user", "technician", "admin"
}

type RepairTicket struct {
	gorm.Model
	Title        string `json:"title" gorm:"not null"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	Status       string `json:"status" gorm:"default:Pending"` // Pending, In Progress, Completed, Cancelled
	UserID       uint   `json:"user_id"`
	User         User   `json:"user" gorm:"foreignKey:UserID"`
	TechnicianID *uint  `json:"technician_id"`
	Technician   *User  `json:"technician" gorm:"foreignKey:TechnicianID"`
}
