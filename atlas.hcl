env "dev" {
  src = "file://internal/db/migrations/schema.sql"
  dir = "file://migrations"
  dev = "docker://postgres/17/dev?search_path=public"  
  url = "postgres://app_user:my-password@localhost:5432/testdb?sslmode=disable"
}

