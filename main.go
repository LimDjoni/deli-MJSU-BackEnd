package main

import (
	"mjsubackend/helper"
	"mjsubackend/model/master/apd"
	"mjsubackend/model/master/bpjskesehatan"
	"mjsubackend/model/master/bpjsketenagakerjaan"
	"mjsubackend/model/master/brand"
	"mjsubackend/model/master/department"
	"mjsubackend/model/master/departmentform"
	"mjsubackend/model/master/doh"
	"mjsubackend/model/master/form"
	"mjsubackend/model/master/history"
	"mjsubackend/model/master/jabatan"
	"mjsubackend/model/master/kartukeluarga"
	"mjsubackend/model/master/ktp"
	"mjsubackend/model/master/laporan"
	"mjsubackend/model/master/mcu"
	"mjsubackend/model/master/npwp"
	"mjsubackend/model/master/pendidikan"
	"mjsubackend/model/master/position"
	"mjsubackend/model/master/role"
	"mjsubackend/model/master/roleform"
	"mjsubackend/model/master/sertifikat"
	"mjsubackend/model/master/userrole"
	"mjsubackend/model/user"
	"mjsubackend/model/userposition"
	routing2 "mjsubackend/routing"

	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var port string
	port = "8080"

	dbUrl := helper.GetEnvWithKey("DATABASE_URL")
	dbUrlStg := helper.GetEnvWithKey("DATABASE_URL_STAGING")

	var dsn string

	if len(dbUrl) > 0 {
		dsn = dbUrl
	} else {
		dsn = dbUrlStg
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		createDB(dsn)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect database after creation: " + err.Error())
		}
	}

	if db != nil {
		// Auto Migrate All Table
		errMigrate := db.AutoMigrate(
			&position.Position{},
			&user.User{},
			&department.Department{},
			&role.Role{},
			&kartukeluarga.KartuKeluarga{},
			&ktp.KTP{},
			&pendidikan.Pendidikan{},
			&jabatan.Jabatan{},
			&sertifikat.Sertifikat{},
			&mcu.MCU{},
			&laporan.Laporan{},
			&apd.APD{},
			&npwp.NPWP{},
			&bpjskesehatan.BPJSKesehatan{},
			&bpjsketenagakerjaan.BPJSKetenagakerjaan{},
			&history.History{},
			&form.Form{},
			&userrole.UserRole{},
			&departmentform.DepartmentForm{},
			&roleform.RoleForm{},
			&brand.Brand{},
			&doh.DOH{},
			&userposition.UserPosition{},
		)
		fmt.Println(errMigrate)
	}

	var validate = validator.New()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, OPTIONS, PUT, DELETE",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization",
		MaxAge:           2592000,
	}))

	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	app.Listen(":" + port)

}

func createDB(dsn string) {
	// Base DSN use for if there is no database (only for creating new database)
	baseDsn := helper.GetEnvWithKey("BASE_DATABASE_URL_STAGING")
	db, err := gorm.Open(postgres.Open(baseDsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to base DB:", err)
		return
	} else {
		dbName := helper.GetEnvWithKey("DATABASE_NAME")
		dbExec := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		db = db.Exec(dbExec)

		if db.Error != nil {
			fmt.Println(db.Error)
			errAssumingExist := fmt.Sprintf("Unable to create DB %s, attempting to connect assuming it exists...", dbName)
			fmt.Println(errAssumingExist)
		}
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect db")
		return
	}
}

func Setup(db *gorm.DB, validate *validator.Validate, route fiber.Router) {
	routing2.UserRouting(db, route, validate)
	routing2.MasterRouting(db, route, validate)
	routing2.EmployeeRouting(db, route, validate)
}
