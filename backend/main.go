package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbname   = "repair_db"
)

var SecretKey = "my_super_secret_key_change_this"

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("การเชื่อมต่อฐานข้อมูลล้มเหลว: " + err.Error())
	}

	fmt.Println("การเชื่อมต่อฐานข้อมูลสำเร็จ... ^_^")

	// AutoMigrate ตารางทั้งหมด
	db.AutoMigrate(&User{}, &RepairTicket{})
	fmt.Println("Migration สำเร็จ... ^_^")

	app := fiber.New()
	app.Use(cors.New())

	// เรียกใช้งาน Route ที่แยกไฟล์ไว้
	repairRoute(app, db)
	// (สามารถใส่ Auth Routes เช่น Register / Login ตรงนี้หรือแยกไฟล์เพิ่มได้ตามสะดวกครับ)

	log.Fatal(app.Listen(":8080"))
}
