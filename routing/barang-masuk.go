package routing

import (
	"mjsubackend/handler"
	"mjsubackend/helper"
	barangmasuk "mjsubackend/model/barang-masuk"
	"mjsubackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func BarangMasukRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	barangmasukRepository := barangmasuk.NewRepository(db)
	barangmasukService := barangmasuk.NewService(barangmasukRepository)

	barangmasukHandler := handler.NewBarangMasukHandler(userService, barangmasukService, validate)

	barangmasukRouting := app.Group("/barangmasuk") // /api

	barangmasukRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	barangmasukRouting.Post("/create", barangmasukHandler.CreateBarangMasuk)
	barangmasukRouting.Get("/list", barangmasukHandler.GetBarangMasuk)
	barangmasukRouting.Get("/list/pagination", barangmasukHandler.GetListBarangMasuk)
	barangmasukRouting.Get("/detail/:id", barangmasukHandler.GetBarangMasukById)
	barangmasukRouting.Put("/update/:id", barangmasukHandler.UpdateBarangMasuk)
	barangmasukRouting.Delete("/delete/:id", barangmasukHandler.DeleteBarangMasuk)

}
