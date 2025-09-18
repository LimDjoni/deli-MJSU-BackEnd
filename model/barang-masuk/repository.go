package barangmasuk

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Repository interface {
	CreateBarangMasuk(barangmasuks RegisterBarangMasukInput) (BarangMasuk, error)
	FindBarangMasuk() ([]BarangMasuk, error)
	FindBarangMasukById(id uint) (BarangMasuk, error)
	ListBarangMasuk(page int, sortFilter SortFilterBarangMasuk) (Pagination, error)
	UpdateBarangMasuk(inputBarangMasuk RegisterBarangMasukInput, id int) (BarangMasuk, error)
	DeleteBarangMasuk(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateBarangMasuk(BarangMasukInput RegisterBarangMasukInput) (BarangMasuk, error) {
	var newBarangMasuk BarangMasuk

	newBarangMasuk.AssetID = BarangMasukInput.AssetID
	newBarangMasuk.Tanggal = BarangMasukInput.Tanggal
	newBarangMasuk.JumlahMasuk = BarangMasukInput.JumlahMasuk

	err := r.db.Create(&newBarangMasuk).Error
	if err != nil {
		return newBarangMasuk, err
	}

	return newBarangMasuk, nil
}

func (r *repository) FindBarangMasuk() ([]BarangMasuk, error) {
	var barangmasuks []BarangMasuk

	queryFilter := "barang_masuks.id > 0"

	errFind := r.db.
		Preload("Asset").
		Where(queryFilter).
		Order("id ASC").Find(&barangmasuks).Error

	return barangmasuks, errFind
}

func (r *repository) FindBarangMasukById(id uint) (BarangMasuk, error) {
	var barangmasuks BarangMasuk

	errFind := r.db.
		Preload("Asset").
		Where("id = ?", id).
		Order("id ASC").First(&barangmasuks).Error

	return barangmasuks, errFind
}

func (r *repository) ListBarangMasuk(page int, sortFilter SortFilterBarangMasuk) (Pagination, error) {
	var listBarangMasuk []BarangMasuk
	var pagination Pagination

	pagination.Limit = 10
	pagination.Page = page
	queryFilter := "barang_masuks.id > 0"
	querySort := "barang_masuks.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.AssetID != "" {
		queryFilter = queryFilter + " AND assets.asset_type ILIKE '%" + sortFilter.AssetID + "%'"
	}

	if sortFilter.Tanggal != "" {
		queryFilter = queryFilter + " AND tanggal ILIKE '%" + sortFilter.Tanggal + "%'"
	}

	if sortFilter.JumlahMasuk != "" {
		queryFilter = queryFilter + " AND jumlah_masuk ILIKE '%" + sortFilter.JumlahMasuk + "%'"
	}

	errFind := r.db.
		Table("barang_masuks").
		Joins("LEFT JOIN assets ON assets.id = barang_masuks.asset_id").
		Preload("Asset").Where(queryFilter).Order(querySort).Scopes(paginateData(listBarangMasuk, &pagination, r.db, queryFilter)).Find(&listBarangMasuk).Error
	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listBarangMasuk

	return pagination, nil
}

func UintToPtr(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func (r *repository) UpdateBarangMasuk(inputBarangMasuk RegisterBarangMasukInput, id int) (BarangMasuk, error) {

	var updatedBarangMasuk BarangMasuk
	errFind := r.db.Where("id = ?", id).First(&updatedBarangMasuk).Error

	if errFind != nil {
		return updatedBarangMasuk, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputBarangMasuk)

	if errorMarshal != nil {
		return updatedBarangMasuk, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedBarangMasuk, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedBarangMasuk).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedBarangMasuk, updateErr
	}

	return updatedBarangMasuk, nil
}

func (r *repository) DeleteBarangMasuk(id uint) (bool, error) {
	tx := r.db.Begin()
	var barangmasuks BarangMasuk

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&barangmasuks).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&barangmasuks).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}
