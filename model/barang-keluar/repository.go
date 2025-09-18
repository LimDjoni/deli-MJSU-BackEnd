package barangkeluar

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBarangKeluar(barangkeluars RegisterBarangKeluarInput) (BarangKeluar, error)
	FindBarangKeluar() ([]BarangKeluar, error)
	FindBarangKeluarById(id uint) (BarangKeluar, error)
	ListBarangKeluar(page int, sortFilter SortFilterBarangKeluar) (Pagination, error)
	UpdateBarangKeluar(inputBarangKeluar RegisterBarangKeluarInput, id int) (BarangKeluar, error)
	DeleteBarangKeluar(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateBarangKeluar(BarangKeluarInput RegisterBarangKeluarInput) (BarangKeluar, error) {
	var newBarangKeluar BarangKeluar

	newBarangKeluar.AssetID = BarangKeluarInput.AssetID
	newBarangKeluar.EmployeeID = BarangKeluarInput.EmployeeID
	newBarangKeluar.Tanggal = BarangKeluarInput.Tanggal
	newBarangKeluar.JumlahKeluar = BarangKeluarInput.JumlahKeluar

	err := r.db.Create(&newBarangKeluar).Error
	if err != nil {
		return newBarangKeluar, err
	}

	return newBarangKeluar, nil
}

func (r *repository) FindBarangKeluar() ([]BarangKeluar, error) {
	var barangkeluars []BarangKeluar

	queryFilter := "barang_keluars.id > 0"

	errFind := r.db.
		Preload("Asset").
		Preload("Employee").
		Where(queryFilter).
		Order("id ASC").Find(&barangkeluars).Error

	return barangkeluars, errFind
}

func (r *repository) FindBarangKeluarById(id uint) (BarangKeluar, error) {
	var barangkeluars BarangKeluar

	errFind := r.db.
		Preload("Asset").
		Preload("Employee").
		Where("id = ?", id).
		Order("id ASC").First(&barangkeluars).Error

	return barangkeluars, errFind
}

func (r *repository) ListBarangKeluar(page int, sortFilter SortFilterBarangKeluar) (Pagination, error) {
	var listBarangKeluar []BarangKeluar
	var pagination Pagination

	pagination.Limit = 10
	pagination.Page = page
	queryFilter := "barang_keluars.id > 0"
	querySort := "barang_keluars.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.EmployeeID != "" {
		queryFilter = queryFilter + " AND employee_id = " + sortFilter.EmployeeID
	}

	if sortFilter.Department != "" {
		queryFilter = queryFilter + " AND department_id = " + sortFilter.Department
	}

	if sortFilter.AssetID != "" {
		queryFilter = queryFilter + " AND assets.asset_type ILIKE '%" + sortFilter.AssetID + "%'"
	}

	if sortFilter.Tanggal != "" {
		queryFilter = queryFilter + " AND tanggal ILIKE '%" + sortFilter.Tanggal + "%'"
	}

	if sortFilter.JumlahKeluar != "" {
		queryFilter = queryFilter + " AND jumlah_keluar ILIKE '%" + sortFilter.JumlahKeluar + "%'"
	}

	errFind := r.db.
		Table("barang_keluars").
		Joins("LEFT JOIN assets ON assets.id = barang_keluars.asset_id").
		Joins("LEFT JOIN employees ON employees.id = barang_keluars.employee_id").
		Preload("Asset").
		Preload("Employee").
		Preload("Employee.Department").Where(queryFilter).Order(querySort).Scopes(paginateData(listBarangKeluar, &pagination, r.db, queryFilter)).Find(&listBarangKeluar).Error
	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listBarangKeluar

	return pagination, nil
}

func UintToPtr(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func (r *repository) UpdateBarangKeluar(inputBarangKeluar RegisterBarangKeluarInput, id int) (BarangKeluar, error) {

	var updatedBarangKeluar BarangKeluar
	errFind := r.db.Where("id = ?", id).First(&updatedBarangKeluar).Error

	if errFind != nil {
		return updatedBarangKeluar, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputBarangKeluar)

	if errorMarshal != nil {
		return updatedBarangKeluar, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedBarangKeluar, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedBarangKeluar).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedBarangKeluar, updateErr
	}

	return updatedBarangKeluar, nil
}

func (r *repository) DeleteBarangKeluar(id uint) (bool, error) {
	tx := r.db.Begin()
	var barangkeluars BarangKeluar

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&barangkeluars).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&barangkeluars).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
