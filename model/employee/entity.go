package employee

import (
	"mjsubackend/model/master/apd"
	"mjsubackend/model/master/bank"
	"mjsubackend/model/master/bpjskesehatan"
	"mjsubackend/model/master/bpjsketenagakerjaan"
	"mjsubackend/model/master/department"
	"mjsubackend/model/master/doh"
	"mjsubackend/model/master/history"
	"mjsubackend/model/master/jabatan"
	"mjsubackend/model/master/kartukeluarga"
	"mjsubackend/model/master/ktp"
	"mjsubackend/model/master/laporan"
	"mjsubackend/model/master/mcu"
	"mjsubackend/model/master/npwp"
	"mjsubackend/model/master/pendidikan"
	"mjsubackend/model/master/position"
	"mjsubackend/model/master/role"
	"mjsubackend/model/master/sertifikat"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	NomorKaryawan         string `json:"nomor_karyawan"`
	DepartmentId          uint   `json:"department_id"`
	Firstname             string `json:"firstname"`
	Lastname              string `json:"lastname"`
	PhoneNumber           string `json:"phone_number"`
	Email                 string `json:"email"`
	Level                 string `json:"level"`
	Status                string `json:"status"`
	RoleId                uint   `json:"role_id"`
	KartuKeluargaId       uint   `json:"kartu_keluarga_id"`
	KTPId                 uint   `json:"ktp_id"`
	PendidikanId          uint   `json:"pendidikan_id"`
	HiredBy               string `json:"hire_by"`
	LaporanId             uint   `json:"laporan_id"`
	APDId                 *uint  `json:"apd_id"`
	BankId                uint   `json:"bank_id"`
	BPJSKesehatanId       *uint  `json:"bpjs_kesehatan_id"`
	BPJSKetenagakerjaanId *uint  `json:"bpjs_ketenagakerjaan_id"`
	NPWPId                *uint  `json:"npwp_id"`
	PositionId            uint   `json:"position_id"`
	DateOfHire            string `json:"date_of_hire"`

	Department          department.Department                    `json:"Department"`
	Role                role.Role                                `json:"Role"`
	KartuKeluarga       kartukeluarga.KartuKeluarga              `json:"KartuKeluarga"`
	KTP                 ktp.KTP                                  `json:"KTP"`
	Pendidikan          pendidikan.Pendidikan                    `json:"Pendidikan"`
	DOH                 []doh.DOH                                `json:"DOH"`
	Jabatan             []jabatan.Jabatan                        `json:"Jabatan"`
	Sertifikat          []sertifikat.Sertifikat                  `json:"Sertifikat"`
	MCU                 []mcu.MCU                                `json:"MCU"`
	Laporan             laporan.Laporan                          `json:"Laporan"`
	APD                 *apd.APD                                 `json:"APD"`
	NPWP                *npwp.NPWP                               `json:"NPWP"`
	Bank                bank.Bank                                `json:"Bank"`
	BPJSKesehatan       *bpjskesehatan.BPJSKesehatan             `json:"BPJSKesehatan"`
	BPJSKetenagakerjaan *bpjsketenagakerjaan.BPJSKetenagakerjaan `json:"BPJSKetenagakerjaan"`
	History             *[]history.History                       `json:"History"`
	Position            position.Position                        `json:"Position"`
}
