package main

import "gorm.io/gorm"

func createRepairTicket(db *gorm.DB, ticket *RepairTicket) error {
	return db.Create(ticket).Error
}

func getAllRepairTickets(db *gorm.DB) ([]RepairTicket, error) {
	var tickets []RepairTicket
	result := db.Preload("User").Preload("Technician").Find(&tickets)
	return tickets, result.Error
}

func getRepairTicketByID(db *gorm.DB, id int) (*RepairTicket, error) {
	var ticket RepairTicket
	result := db.Preload("User").Preload("Technician").First(&ticket, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ticket, nil
}

func updateRepairTicket(db *gorm.DB, id int, updates map[string]interface{}) error {
	return db.Model(&RepairTicket{}).Where("id = ?", id).Updates(updates).Error
}

func deleteRepairTicket(db *gorm.DB, id int) error {
	return db.Delete(&RepairTicket{}, id).Error
}
