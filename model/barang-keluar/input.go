package barangkeluar

type RegisterBarangKeluarInput struct {
	AssetID      int    `json:"asset_id"`
	EmployeeID   int    `json:"employee_id"`
	Tanggal      string `json:"tanggal"`
	JumlahKeluar int    `json:"jumlah_keluar"`
}

type SortFilterBarangKeluar struct {
	Field        string
	Sort         string
	AssetID      string `json:"asset_id"`
	EmployeeID   string `json:"employee_id"`
	Department   string `json:"department"`
	Tanggal      string `json:"tanggal"`
	JumlahKeluar string `json:"jumlah_keluar"`
}
