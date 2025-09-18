package routing

import (
	"mjsubackend/handler"
	"mjsubackend/helper"
	barangkeluar "mjsubackend/model/barang-keluar"
	"mjsubackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func BarangKeluarRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	barangkeluarRepository := barangkeluar.NewRepository(db)
	barangkeluarService := barangkeluar.NewService(barangkeluarRepository)

	barangkeluarHandler := handler.NewBarangKeluarHandler(userService, barangkeluarService, validate)

	barangkeluarRouting := app.Group("/barangkeluar") // /api

	barangkeluarRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	barangkeluarRouting.Post("/create", barangkeluarHandler.CreateBarangKeluar)
	barangkeluarRouting.Get("/list", barangkeluarHandler.GetBarangKeluar)
	barangkeluarRouting.Get("/list/pagination", barangkeluarHandler.GetListBarangKeluar)
	barangkeluarRouting.Get("/detail/:id", barangkeluarHandler.GetBarangKeluarById)
	barangkeluarRouting.Put("/update/:id", barangkeluarHandler.UpdateBarangKeluar)
	barangkeluarRouting.Delete("/delete/:id", barangkeluarHandler.DeleteBarangKeluar)

}
