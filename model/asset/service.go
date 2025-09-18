package asset

type Service interface {
	CreateAsset(assets RegisterAssetInput) (Asset, error)
	FindAsset() ([]Asset, error)
	FindAssetById(id uint) (Asset, error)
	FindAssetByName(assetName string) (Asset, error)
	GetListAsset(page int, sortFilter SortFilterAsset) (Pagination, error)
	UpdateAsset(inputAsset RegisterAssetInput, id int) (Asset, error)
	DeleteAsset(id uint) (bool, error)
	ExportReportAsset(sortFilter AssetSummary) ([]AssetSummary, error)
	ListReportAsset(page int, sortFilter AssetSummary) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateAsset(assets RegisterAssetInput) (Asset, error) {
	newAsset, err := s.repository.CreateAsset(assets)

	return newAsset, err
}

func (s *service) FindAsset() ([]Asset, error) {
	assets, err := s.repository.FindAsset()

	return assets, err
}

func (s *service) FindAssetById(id uint) (Asset, error) {
	assets, err := s.repository.FindAssetById(id)

	return assets, err
}

func (s *service) FindAssetByName(assetName string) (Asset, error) {
	assets, err := s.repository.FindAssetByName(assetName)

	return assets, err
}

func (s *service) GetListAsset(page int, sortFilter SortFilterAsset) (Pagination, error) {
	listListAsset, listListAssetErr := s.repository.ListAsset(page, sortFilter)

	return listListAsset, listListAssetErr
}

func (s *service) UpdateAsset(inputAsset RegisterAssetInput, id int) (Asset, error) {
	updateAsset, updateAssetErr := s.repository.UpdateAsset(inputAsset, id)

	return updateAsset, updateAssetErr
}

func (s *service) DeleteAsset(id uint) (bool, error) {
	isDeletedAsset, isDeletedAssetErr := s.repository.DeleteAsset(id)

	return isDeletedAsset, isDeletedAssetErr
}

func (s *service) ExportReportAsset(sortFilter AssetSummary) ([]AssetSummary, error) {
	listListAsset, listListAssetErr := s.repository.ExportReportAsset(sortFilter)

	return listListAsset, listListAssetErr
}

func (s *service) ListReportAsset(page int, sortFilter AssetSummary) (Pagination, error) {
	listListAsset, listListAssetErr := s.repository.ListReportAsset(page, sortFilter)

	return listListAsset, listListAssetErr
}
