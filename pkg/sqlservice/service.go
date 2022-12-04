package sqlservice

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type SQLService struct {
	username string
	password string
	host     string
	port     string
	dbName   string
	db       *sqlx.DB
	lg       *log.Logger
}

func NewSQLService(username, password, host, port, dbName string, lg *log.Logger) *SQLService {
	return &SQLService{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
		lg:       lg,
	}
}

func (s *SQLService) Start(ctx context.Context) error {
	err := s.connect(ctx)
	if err != nil {
		return fmt.Errorf("start connection: %v", err)
	}
	return nil
}

func (s *SQLService) connect(ctx context.Context) (err error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		s.username, s.password, s.host, s.port, s.dbName)
	s.db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("open SQL connection: %v", err)
	}

	err = s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("ping: %v", err)
	}
	return nil
}
