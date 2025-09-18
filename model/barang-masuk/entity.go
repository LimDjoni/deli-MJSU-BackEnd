package barangmasuk

import (
	"mjsubackend/model/asset"

	"gorm.io/gorm"
)

type BarangMasuk struct {
	gorm.Model
	AssetID     int    `json:"asset_id"`
	Tanggal     string `json:"tanggal"`
	JumlahMasuk int    `json:"jumlah_masuk"`

	Asset asset.Asset `json:"Asset"`
}
