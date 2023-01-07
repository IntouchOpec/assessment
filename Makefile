.EXPORT_ALL_VARIABLES:

PORT=:2565
DB_HOST=localhost
DB_PORT=5432
DB_USER=pre_test_assessment
DB_PASSWORD=pS5h140Evri1
DB_NAME=dev_assessment
DATABASE_URL=host=localhost port=5432 user=pre_test_assessment password=pS5h140Evri1 dbname=dev_assessment

run:
	PORT=:2565 DATABASE_URL="host=$$DB_HOST port=$$DB_PORT user=$$DB_USER password=$$DB_PASSWORD dbname=$$DB_NAME sslmode=disable" go run server.go
