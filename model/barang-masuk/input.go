package barangmasuk

type RegisterBarangMasukInput struct {
	AssetID     int    `json:"asset_id"`
	Tanggal     string `json:"tanggal"`
	JumlahMasuk int    `json:"jumlah_masuk"`
}

type SortFilterBarangMasuk struct {
	Field       string
	Sort        string
	AssetID     string `json:"asset_id"`
	Tanggal     string `json:"tanggal"`
	JumlahMasuk string `json:"jumlah_masuk"`
}
