migrate:
	goose -dir "migration" postgres "user=postgres dbname=course_work host=localhost password=postgrespw sslmode=disable" up
	goose -dir "migration" postgres "user=postgres dbname=course_work host=localhost password=postgrespw sslmode=disable" status
