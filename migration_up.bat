goose postgres "postgres://app_user:my-password@localhost:5432/testdb?sslmode=disable&search_path=public" -dir migration up 
