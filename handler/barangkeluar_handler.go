package handler

import (
	barangkeluar "mjsubackend/model/barang-keluar"
	"mjsubackend/model/user"
	"mjsubackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type barangkeluarHandler struct {
	userService         user.Service
	barangkeluarService barangkeluar.Service
	v                   *validator.Validate
}

func NewBarangKeluarHandler(userService user.Service, barangkeluarsService barangkeluar.Service, v *validator.Validate) *barangkeluarHandler {
	return &barangkeluarHandler{
		userService,
		barangkeluarsService,
		v,
	}
}

// Master Data

func (h *barangkeluarHandler) CreateBarangKeluar(c *fiber.Ctx) error {
	barangkeluarsInput := new(barangkeluar.RegisterBarangKeluarInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(barangkeluarsInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*barangkeluarsInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdBarangKeluar, createdBarangKeluarErr := h.barangkeluarService.CreateBarangKeluar(*barangkeluarsInput)

	if createdBarangKeluarErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdBarangKeluarErr.Error(),
		})
	}

	return c.Status(201).JSON(createdBarangKeluar)
}

func (h *barangkeluarHandler) GetBarangKeluar(c *fiber.Ctx) error {
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

	brands, brandsErr := h.barangkeluarService.FindBarangKeluar()

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *barangkeluarHandler) GetBarangKeluarById(c *fiber.Ctx) error {
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

	heavyEquipments, heavyEquipmentsErr := h.barangkeluarService.FindBarangKeluarById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *barangkeluarHandler) GetListBarangKeluar(c *fiber.Ctx) error {
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

	var filterBarangKeluar barangkeluar.SortFilterBarangKeluar

	filterBarangKeluar.AssetID = c.Query("asset_id")
	filterBarangKeluar.EmployeeID = c.Query("employee_id")
	filterBarangKeluar.Department = c.Query("department")
	filterBarangKeluar.Tanggal = c.Query("tanggal")
	filterBarangKeluar.JumlahKeluar = c.Query("jumlah_keluar")
	filterBarangKeluar.Field = c.Query("field")
	filterBarangKeluar.Sort = c.Query("sort")

	listBarangKeluar, listBarangKeluarErr := h.barangkeluarService.GetListBarangKeluar(pageNumber, filterBarangKeluar)

	if listBarangKeluarErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listBarangKeluarErr.Error(),
		})
	}

	return c.Status(200).JSON(listBarangKeluar)
}

func (h *barangkeluarHandler) UpdateBarangKeluar(c *fiber.Ctx) error {
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
	inputUpdateBarangKeluar := new(barangkeluar.RegisterBarangKeluarInput)
	if err := c.BodyParser(inputUpdateBarangKeluar); err != nil {
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

	updateBarangKeluar, err := h.barangkeluarService.UpdateBarangKeluar(*inputUpdateBarangKeluar, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateBarangKeluar)
}

func (h *barangkeluarHandler) DeleteBarangKeluar(c *fiber.Ctx) error {
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

	barangkeluar, err := h.barangkeluarService.FindBarangKeluarById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find barangkeluar",
			"error":   err.Error(),
		})
	}

	// Optional: Use barangkeluar.ID for extra safety
	if _, err := h.barangkeluarService.DeleteBarangKeluar(barangkeluar.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete barangkeluar",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete barangkeluar",
	})
}
