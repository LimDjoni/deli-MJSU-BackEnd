package barangkeluar

import (
	"mjsubackend/model/asset"
	"mjsubackend/model/employee"

	"gorm.io/gorm"
)

type BarangKeluar struct {
	gorm.Model
	AssetID      int    `json:"asset_id"`
	EmployeeID   int    `json:"employee_id"`
	Tanggal      string `json:"tanggal"`
	JumlahKeluar int    `json:"jumlah_keluar"`

	Asset    asset.Asset       `json:"Asset"`
	Employee employee.Employee `json:"employee"`
}
