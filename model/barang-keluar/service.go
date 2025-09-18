package barangkeluar

type Service interface {
	CreateBarangKeluar(barangkeluars RegisterBarangKeluarInput) (BarangKeluar, error)
	FindBarangKeluar() ([]BarangKeluar, error)
	FindBarangKeluarById(id uint) (BarangKeluar, error)
	GetListBarangKeluar(page int, sortFilter SortFilterBarangKeluar) (Pagination, error)
	UpdateBarangKeluar(inputBarangKeluar RegisterBarangKeluarInput, id int) (BarangKeluar, error)
	DeleteBarangKeluar(id uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateBarangKeluar(barangkeluars RegisterBarangKeluarInput) (BarangKeluar, error) {
	newBarangKeluar, err := s.repository.CreateBarangKeluar(barangkeluars)

	return newBarangKeluar, err
}

func (s *service) FindBarangKeluar() ([]BarangKeluar, error) {
	barangkeluars, err := s.repository.FindBarangKeluar()

	return barangkeluars, err
}

func (s *service) FindBarangKeluarById(id uint) (BarangKeluar, error) {
	barangkeluars, err := s.repository.FindBarangKeluarById(id)

	return barangkeluars, err
}

func (s *service) GetListBarangKeluar(page int, sortFilter SortFilterBarangKeluar) (Pagination, error) {
	listListBarangKeluar, listListBarangKeluarErr := s.repository.ListBarangKeluar(page, sortFilter)

	return listListBarangKeluar, listListBarangKeluarErr
}

func (s *service) UpdateBarangKeluar(inputBarangKeluar RegisterBarangKeluarInput, id int) (BarangKeluar, error) {
	updateBarangKeluar, updateBarangKeluarErr := s.repository.UpdateBarangKeluar(inputBarangKeluar, id)

	return updateBarangKeluar, updateBarangKeluarErr
}

func (s *service) DeleteBarangKeluar(id uint) (bool, error) {
	isDeletedBarangKeluar, isDeletedBarangKeluarErr := s.repository.DeleteBarangKeluar(id)

	return isDeletedBarangKeluar, isDeletedBarangKeluarErr
}
