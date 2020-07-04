package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	gorsk "github.com/kerti/balances/backend"
	"github.com/kerti/balances/backend/pkg/utl/secure"
	"github.com/satori/uuid"
)

func main() {
	var host = os.Getenv("DB_HOST")
	var port = os.Getenv("DB_PORT")
	var user = os.Getenv("DB_USER")
	var pass = os.Getenv("DB_PASS")
	var name = os.Getenv("DB_NAME")
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", user, pass, host, port, name)

	db, err := gorm.Open("mysql", dsn)
	checkErr(err)
	defer db.Close()

	db.AutoMigrate(
		&gorsk.Company{},
		&gorsk.Location{},
		&gorsk.Role{},
		&gorsk.User{})

	var companyID = uuid.NewV4()
	var locationID = uuid.NewV4()
	var userID = uuid.NewV4()
	var superUserRoleID = uuid.NewV4()
	sec := secure.New(1, nil)

	db.Create(&gorsk.Company{ID: companyID, Name: "admin_company", Active: true})

	db.Create(&gorsk.Location{
		ID:        locationID,
		Name:      "admin_location",
		Active:    true,
		Address:   "admin_address",
		CompanyID: companyID})

	db.Create(&gorsk.Role{ID: superUserRoleID, AccessLevel: 100, Name: "SUPER_ADMIN"})
	db.Create(&gorsk.Role{ID: uuid.NewV4(), AccessLevel: 110, Name: "ADMIN"})
	db.Create(&gorsk.Role{ID: uuid.NewV4(), AccessLevel: 120, Name: "COMPANY_ADMIN"})
	db.Create(&gorsk.Role{ID: uuid.NewV4(), AccessLevel: 130, Name: "LOCATION_ADMIN"})
	db.Create(&gorsk.Role{ID: uuid.NewV4(), AccessLevel: 140, Name: "USER"})

	db.Create(&gorsk.User{
		ID:         userID,
		FirstName:  "John",
		LastName:   "Doe",
		Username:   "johndoe",
		Password:   sec.Hash("admin"),
		Email:      "johndoe@mail.com",
		Mobile:     "+6280989999",
		Phone:      "+62274147",
		Address:    "1234 Alpha Stret",
		Active:     true,
		RoleID:     superUserRoleID,
		CompanyID:  companyID,
		LocationID: locationID,
	})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
