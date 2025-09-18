package routing

import (
	"mjsubackend/handler"
	"mjsubackend/helper"
	"mjsubackend/model/asset"
	"mjsubackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func AssetRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	assetRepository := asset.NewRepository(db)
	assetService := asset.NewService(assetRepository)

	assetHandler := handler.NewAssetHandler(userService, assetService, validate)

	assetRouting := app.Group("/asset") // /api

	assetRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	assetRouting.Post("/create", assetHandler.CreateAsset)
	assetRouting.Get("/list", assetHandler.GetAsset)
	assetRouting.Get("/list/pagination", assetHandler.GetListAsset)
	assetRouting.Get("/detail/:id", assetHandler.GetAssetById)
	assetRouting.Get("/detail/ByName/:assetname", assetHandler.GetAssetByName)
	assetRouting.Put("/update/:id", assetHandler.UpdateAsset)
	assetRouting.Delete("/delete/:id", assetHandler.DeleteAsset)
	assetRouting.Get("/export", assetHandler.GetExportReportAsset)
	assetRouting.Get("/list/rangkuman", assetHandler.GetListReportAsset)

}
