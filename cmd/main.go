package main

import (
	"context"
	"github.com/hmuriyMax/DatabaseCW/pkg/httpservice"
	"github.com/hmuriyMax/DatabaseCW/pkg/sqlservice"
	"log"
	"time"
)

func main() {
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Lmicroseconds)

	sqlSvc := sqlservice.NewSQLService("postgres", "postgrespw",
		"127.0.0.1", "55000", "course_work", logger)
	sqlContext, cancelSql := context.WithTimeout(context.Background(), time.Second)
	defer cancelSql()
	err := sqlSvc.Start(sqlContext)
	if err != nil {
		log.Fatal(err)
	}

	httpSvc := httpservice.NewHTTPService(80, "127.0.0.1", logger, false)
	httpSvc.ConnectToDataBase(sqlSvc)
	httpSvc.Start()

	select {
	case err := <-httpSvc.GetErrChan():
		if err != nil {
			log.Fatal(err)
		}
	}
}
