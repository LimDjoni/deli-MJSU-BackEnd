package asset

import (
	"encoding/json"
	"math"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	CreateAsset(assets RegisterAssetInput) (Asset, error)
	FindAsset() ([]Asset, error)
	FindAssetById(id uint) (Asset, error)
	FindAssetByName(assetName string) (Asset, error)
	ListAsset(page int, sortFilter SortFilterAsset) (Pagination, error)
	UpdateAsset(inputAsset RegisterAssetInput, id int) (Asset, error)
	DeleteAsset(id uint) (bool, error)
	ExportReportAsset(sortFilter AssetSummary) ([]AssetSummary, error)
	ListReportAsset(page int, sortFilter AssetSummary) (Pagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateAsset(AssetInput RegisterAssetInput) (Asset, error) {
	var newAsset Asset

	newAsset.AssetType = AssetInput.AssetType
	newAsset.Ukuran = AssetInput.Ukuran
	newAsset.Stock = AssetInput.Stock

	err := r.db.Create(&newAsset).Error
	if err != nil {
		return newAsset, err
	}

	return newAsset, nil
}

func (r *repository) FindAsset() ([]Asset, error) {
	var assets []Asset

	queryFilter := "assets.id > 0"

	errFind := r.db.
		Where(queryFilter).
		Order("id ASC").Find(&assets).Error

	return assets, errFind
}

func (r *repository) FindAssetById(id uint) (Asset, error) {
	var assets Asset

	errFind := r.db.
		Where("id = ?", id).
		Order("id ASC").First(&assets).Error

	return assets, errFind
}

func (r *repository) FindAssetByName(assetName string) (Asset, error) {
	var assets Asset

	errFind := r.db.
		Where("asset_type = ?", assetName).
		Order("id ASC").First(&assets).Error

	return assets, errFind
}

func (r *repository) ListAsset(page int, sortFilter SortFilterAsset) (Pagination, error) {
	var listAsset []Asset
	var pagination Pagination

	pagination.Limit = 10
	pagination.Page = page
	queryFilter := "assets.id > 0"
	querySort := "assets.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	if sortFilter.AssetType != "" {
		queryFilter = queryFilter + " AND asset_type ILIKE '%" + sortFilter.AssetType + "%'"
	}

	if sortFilter.Ukuran != "" {
		queryFilter = queryFilter + " AND ukuran ILIKE '%" + sortFilter.Ukuran + "%'"
	}

	if sortFilter.Stock != "" {
		queryFilter = queryFilter + " AND cast(stock AS TEXT) ILIKE '%" + sortFilter.Stock + "%'"
	}

	errFind := r.db.Where(queryFilter).Order(querySort).Scopes(paginateData(listAsset, &pagination, r.db, queryFilter)).Find(&listAsset).Error
	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listAsset

	return pagination, nil
}

func UintToPtr(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func (r *repository) UpdateAsset(inputAsset RegisterAssetInput, id int) (Asset, error) {

	var updatedAsset Asset
	errFind := r.db.Where("id = ?", id).First(&updatedAsset).Error

	if errFind != nil {
		return updatedAsset, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputAsset)

	if errorMarshal != nil {
		return updatedAsset, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedAsset, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedAsset).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedAsset, updateErr
	}

	return updatedAsset, nil
}

func (r *repository) DeleteAsset(id uint) (bool, error) {
	tx := r.db.Begin()
	var assets Asset

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&assets).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&assets).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) ExportReportAsset(sortFilter AssetSummary) ([]AssetSummary, error) {
	var results []AssetSummary

	var filters []string

	// ✅ Only non-deleted assets
	filters = append(filters, "a.deleted_at IS NULL")

	// ✅ Build main query
	tx := r.db.
		Table("assets a").
		Select(`
			a.asset_type,
			a.ukuran,
			a.stock,
			COALESCE(bm.total_masuk, 0) AS jumlah_barang_masuk,
			COALESCE(bk.total_keluar, 0) AS jumlah_barang_keluar
		`).
		// ✅ Subquery for barang_masuks (aggregate first)
		Joins(`
			LEFT JOIN (
				SELECT asset_id, SUM(jumlah_masuk) AS total_masuk
				FROM barang_masuks
				WHERE TO_CHAR(tanggal::DATE, 'YYYY-MM') = ? 
				AND deleted_at IS NULL
				GROUP BY asset_id
			) bm ON a.id = bm.asset_id
		`, sortFilter.Month).
		// ✅ Subquery for barang_keluars (aggregate first)
		Joins(`
			LEFT JOIN (
				SELECT asset_id, SUM(jumlah_keluar) AS total_keluar
				FROM barang_keluars
				WHERE TO_CHAR(tanggal::DATE, 'YYYY-MM') = ? 
				AND deleted_at IS NULL
				GROUP BY asset_id
			) bk ON a.id = bk.asset_id
		`, sortFilter.Month).
		Where(strings.Join(filters, " AND "))

	// ✅ Execute query and fetch data
	if err := tx.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *repository) ListReportAsset(page int, sortFilter AssetSummary) (Pagination, error) {
	var results []AssetSummary
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page

	var filters []string
	var args []interface{}

	// ✅ Only non-deleted assets
	filters = append(filters, "a.deleted_at IS NULL")

	// ✅ Build main query
	tx := r.db.
		Table("assets a").
		Select(`
			a.asset_type,
			a.ukuran,
			a.stock,
			COALESCE(bm.total_masuk, 0) AS jumlah_barang_masuk,
			COALESCE(bk.total_keluar, 0) AS jumlah_barang_keluar
		`).
		// ✅ Subquery for barang_masuks (aggregate first)
		Joins(`
			LEFT JOIN (
				SELECT asset_id, SUM(jumlah_masuk) AS total_masuk
				FROM barang_masuks
				WHERE TO_CHAR(tanggal::DATE, 'YYYY-MM') = ? AND deleted_at IS NULL
				GROUP BY asset_id
			) bm ON a.id = bm.asset_id
		`, sortFilter.Month).
		// ✅ Subquery for barang_keluars (aggregate first)
		Joins(`
			LEFT JOIN (
				SELECT asset_id, SUM(jumlah_keluar) AS total_keluar
				FROM barang_keluars
				WHERE TO_CHAR(tanggal::DATE, 'YYYY-MM') = ? AND deleted_at IS NULL
				GROUP BY asset_id
			) bk ON a.id = bk.asset_id
		`, sortFilter.Month).
		Where(strings.Join(filters, " AND "), args...).
		Group("a.asset_type, a.ukuran, a.stock, bm.total_masuk, bk.total_keluar")

	// ✅ Count total rows
	var count int64
	if err := r.db.
		Table("assets a").
		Where("a.deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return pagination, err
	}

	pagination.TotalRows = count
	pagination.TotalPages = int(math.Ceil(float64(count) / float64(pagination.Limit)))

	// ✅ Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	if err := tx.Limit(pagination.Limit).Offset(offset).Scan(&results).Error; err != nil {
		return pagination, err
	}

	pagination.Data = results
	return pagination, nil
}
