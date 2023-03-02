package main

import (
	"context"
	"github.com/hmuriyMax/DatabaseCW/internal/httpservice"
	"github.com/hmuriyMax/DatabaseCW/internal/sqlservice"
	"github.com/hmuriyMax/DatabaseCW/internal/testservice"
	"log"
	"time"
)

func main() {
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Lmicroseconds)

	sqlSvc := sqlservice.NewSQLService("postgres", "postgrespw",
		"127.0.0.1", "5432", "course_work", logger)
	sqlContext, cancelSql := context.WithTimeout(context.Background(), time.Second)
	defer cancelSql()
	err := sqlSvc.Start(sqlContext)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := testservice.NewTestService(logger)
	if err != nil {
		logger.Fatal(err)
	}

	httpSvc := httpservice.NewHTTPService(80, "127.0.0.1", logger, false, ts)
	httpSvc.ConnectToDataBase(sqlSvc)
	httpSvc.Start()

	select {
	case err := <-httpSvc.GetErrChan():
		if err != nil {
			log.Fatal(err)
		}
	}
}
