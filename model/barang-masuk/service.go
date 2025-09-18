package barangmasuk

type Service interface {
	CreateBarangMasuk(barangmasuks RegisterBarangMasukInput) (BarangMasuk, error)
	FindBarangMasuk() ([]BarangMasuk, error)
	FindBarangMasukById(id uint) (BarangMasuk, error)
	GetListBarangMasuk(page int, sortFilter SortFilterBarangMasuk) (Pagination, error)
	UpdateBarangMasuk(inputBarangMasuk RegisterBarangMasukInput, id int) (BarangMasuk, error)
	DeleteBarangMasuk(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateBarangMasuk(barangmasuks RegisterBarangMasukInput) (BarangMasuk, error) {
	newBarangMasuk, err := s.repository.CreateBarangMasuk(barangmasuks)

	return newBarangMasuk, err
}

func (s *service) FindBarangMasuk() ([]BarangMasuk, error) {
	barangmasuks, err := s.repository.FindBarangMasuk()

	return barangmasuks, err
}

func (s *service) FindBarangMasukById(id uint) (BarangMasuk, error) {
	barangmasuks, err := s.repository.FindBarangMasukById(id)

	return barangmasuks, err
}

func (s *service) GetListBarangMasuk(page int, sortFilter SortFilterBarangMasuk) (Pagination, error) {
	listListBarangMasuk, listListBarangMasukErr := s.repository.ListBarangMasuk(page, sortFilter)

	return listListBarangMasuk, listListBarangMasukErr
}

func (s *service) UpdateBarangMasuk(inputBarangMasuk RegisterBarangMasukInput, id int) (BarangMasuk, error) {
	updateBarangMasuk, updateBarangMasukErr := s.repository.UpdateBarangMasuk(inputBarangMasuk, id)

	return updateBarangMasuk, updateBarangMasukErr
}

func (s *service) DeleteBarangMasuk(id uint) (bool, error) {
	isDeletedBarangMasuk, isDeletedBarangMasukErr := s.repository.DeleteBarangMasuk(id)

	return isDeletedBarangMasuk, isDeletedBarangMasukErr
}
