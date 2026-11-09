package main

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func repairRoute(app *fiber.App, db *gorm.DB) {
	// ดึงรายการแจ้งซ่อมทั้งหมด
	app.Get("/api/tickets", func(c fiber.Ctx) error {
		tickets, err := getAllRepairTickets(db)
		if err != nil {
			log.Printf("%v", err)
			return c.Status(500).JSON(fiber.Map{"message": "Internal server error"})
		}
		return c.Status(200).JSON(tickets)
	})

	// สร้างรายการแจ้งซ่อม
	app.Post("/api/tickets", func(c fiber.Ctx) error {
		ticket := new(RepairTicket)
		if err := c.Bind().Body(ticket); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
		}

		ticket.Status = "Pending"
		if err := createRepairTicket(db, ticket); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		}

		return c.Status(201).JSON(fiber.Map{"message": "Ticket created successfully", "ticket": ticket})
	})

	// อัปเดตสถานะ / มอบหมายช่าง
	app.Patch("/api/tickets/:id/progress", func(c fiber.Ctx) error {
		id := c.Params("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid request parameter"})
		}

		_, err = getRepairTicketByID(db, idInt)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"message": "Ticket not found"})
		}

		type UpdateInput struct {
			Status       string `json:"status"`
			TechnicianID *uint  `json:"technician_id"`
		}

		var input UpdateInput
		if err := c.Bind().Body(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
		}

		updates := map[string]interface{}{}
		if input.Status != "" {
			updates["status"] = input.Status
		}
		if input.TechnicianID != nil {
			updates["technician_id"] = input.TechnicianID
		}

		if err := updateRepairTicket(db, idInt, updates); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Internal server error"})
		}

		updatedTicket, _ := getRepairTicketByID(db, idInt)
		return c.Status(200).JSON(fiber.Map{"message": "Updated successfully", "ticket": updatedTicket})
	})

	// ลบใบแจ้งซ่อม
	app.Delete("/api/tickets/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "Invalid request parameter"})
		}

		if err := deleteRepairTicket(db, idInt); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Internal server error"})
		}

		return c.Status(200).JSON(fiber.Map{"message": "Ticket deleted successfully"})
	})
}
