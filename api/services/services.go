package services

type Service struct {
	Database Databaser
}

type Databaser interface {
	InsertNewAlert(id, pool string) error
}

func NewService(db Databaser) *Service {
	return &Service{
		Database: db,
	}
}

func (s *Service) AlertSubscribe(id string, pool string) error {
	return s.Database.InsertNewAlert(id, pool)
}
