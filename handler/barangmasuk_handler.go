package handler

import (
	barangmasuk "mjsubackend/model/barang-masuk"
	"mjsubackend/model/user"
	"mjsubackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type barangmasukHandler struct {
	userService        user.Service
	barangmasukService barangmasuk.Service
	v                  *validator.Validate
}

func NewBarangMasukHandler(userService user.Service, barangmasuksService barangmasuk.Service, v *validator.Validate) *barangmasukHandler {
	return &barangmasukHandler{
		userService,
		barangmasuksService,
		v,
	}
}

// Master Data

func (h *barangmasukHandler) CreateBarangMasuk(c *fiber.Ctx) error {
	barangmasuksInput := new(barangmasuk.RegisterBarangMasukInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(barangmasuksInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*barangmasuksInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdBarangMasuk, createdBarangMasukErr := h.barangmasukService.CreateBarangMasuk(*barangmasuksInput)

	if createdBarangMasukErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdBarangMasukErr.Error(),
		})
	}

	return c.Status(201).JSON(createdBarangMasuk)
}

func (h *barangmasukHandler) GetBarangMasuk(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	brands, brandsErr := h.barangmasukService.FindBarangMasuk()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *barangmasukHandler) GetBarangMasukById(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	heavyEquipments, heavyEquipmentsErr := h.barangmasukService.FindBarangMasukById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *barangmasukHandler) GetListBarangMasuk(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	page := c.Query("page")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if page == "" {
		pageNumber = 1
	}

	var filterBarangMasuk barangmasuk.SortFilterBarangMasuk

	filterBarangMasuk.AssetID = c.Query("asset_id")
	filterBarangMasuk.Tanggal = c.Query("tanggal")
	filterBarangMasuk.JumlahMasuk = c.Query("jumlah_masuk")
	filterBarangMasuk.Field = c.Query("field")
	filterBarangMasuk.Sort = c.Query("sort")

	listBarangMasuk, listBarangMasukErr := h.barangmasukService.GetListBarangMasuk(pageNumber, filterBarangMasuk)

	if listBarangMasukErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listBarangMasukErr.Error(),
		})
	}

	return c.Status(200).JSON(listBarangMasuk)
}

func (h *barangmasukHandler) UpdateBarangMasuk(c *fiber.Ctx) error {
	//Get User
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}
	// Check User Login
	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	//Get Input
	inputUpdateBarangMasuk := new(barangmasuk.RegisterBarangMasukInput)
	if err := c.BodyParser(inputUpdateBarangMasuk); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//Get ID
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	updateBarangMasuk, err := h.barangmasukService.UpdateBarangMasuk(*inputUpdateBarangMasuk, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateBarangMasuk)
}

func (h *barangmasukHandler) DeleteBarangMasuk(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID parameter",
			"error":   err.Error(),
		})
	}

	barangmasuk, err := h.barangmasukService.FindBarangMasukById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find barangmasuk",
			"error":   err.Error(),
		})
	}

	// Optional: Use barangmasuk.ID for extra safety
	if _, err := h.barangmasukService.DeleteBarangMasuk(barangmasuk.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete barangmasuk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete barangmasuk",
	})
}
