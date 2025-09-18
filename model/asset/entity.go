package asset

import (
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	AssetType string `json:"asset_type"`
	Ukuran    string `json:"ukuran"`
	Stock     int    `json:"stock"`
}
