package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Employee struct {
	gorm.Model
	Firstname    string
	Lastname     string
	Birthdate    time.Time
	DepartmentID uint
}

type Department struct {
	gorm.Model
	Name      string
	Location  string
	Employees []Employee
}

type DatabaseConfig struct {
	Username, Password, Host, Port, DatabaseName string
}

func (conf DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.DatabaseName,
	)
}

func main() {
	dbconfig := DatabaseConfig{
		Username:     "root",
		Password:     "1234",
		Host:         "localhost",
		Port:         "3306",
		DatabaseName: "sample",
	}
	db, err := gorm.Open(mysql.Open(dbconfig.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Department{}, &Employee{})

	// add some data
	department := Department{
		Name:     "Science",
		Location: "Bangkok",
	}
	db.Create(&department)
	employee1 := Employee{
		Firstname:    "Roger",
		Lastname:     "Stevens",
		Birthdate:    time.Date(1996, 10, 3, 0, 0, 0, 0, time.FixedZone("UTC+7", 7*60*60)),
		DepartmentID: department.ID,
	}
	employee2 := Employee{
		Firstname:    "Piers",
		Lastname:     "Lambert",
		Birthdate:    time.Date(1994, 3, 22, 0, 0, 0, 0, time.FixedZone("UTC+7", 7*60*60)),
		DepartmentID: department.ID,
	}
	db.Create(&employee1)
	db.Create(&employee2)

	// Query data
	type Result struct {
		Firstname string
		Name      string
		Location  string
	}
	var results []Result
	db.
		Model(&Employee{}).
		Select("employees.firstname, departments.name, departments.location").
		Joins("left join departments on departments.id = employees.department_id").
		Scan(&results)
	fmt.Println(results)
}
