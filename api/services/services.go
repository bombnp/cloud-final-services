package services

type Service struct {
	Database Databaser
}

type Databaser interface {
	InsertNewSubscribe(id, pool, t string) error
}

func NewService(db Databaser) *Service {
	return &Service{
		Database: db,
	}
}

func (s *Service) AlertSubscribe(id string, pool string) error {
	return s.Database.InsertNewSubscribe(id, pool, "alert")
}
