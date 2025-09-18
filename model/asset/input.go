package asset

type RegisterAssetInput struct {
	AssetType string `json:"asset_type"`
	Ukuran    string `json:"ukuran"`
	Stock     int    `json:"stock"`
}

type SortFilterAsset struct {
	Field     string
	Sort      string
	AssetType string `json:"asset_type"`
	Ukuran    string `json:"ukuran"`
	Stock     string `json:"stock"`
}

type AssetSummary struct {
	Month              string `json:"month"`
	AssetType          string `json:"asset_type"`
	Ukuran             string `json:"ukuran"`
	Stock              int    `json:"stock"`
	JumlahBarangMasuk  int    `json:"jumlah_barang_masuk"`
	JumlahBarangKeluar int    `json:"jumlah_barang_keluar"`
}
