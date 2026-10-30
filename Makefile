.PHONY: dev dev-web worker test templ migrate migrate-status migrate-rollback seed seed-mysql seed-pg build clean tailwind-install css css-watch

# Jalankan HTTP REST API & Web Server (Default: API-First, ringan, tanpa CSS watcher)
dev:
	templ generate
	go run cmd/server/main.go

# Download standalone Tailwind CSS CLI tanpa butuh Node.js / npm
tailwind-install:
	go run cmd/make/main.go tailwind:install

# Kompilasi Tailwind CSS sekali jalan (minified output)
css:
	@if [ ! -f bin/tailwindcss ]; then go run cmd/make/main.go tailwind:install; fi
	bin/tailwindcss -i ./resources/css/input.css -o ./public/css/output.css --minify

# Pantau perubahan template dan rebuild CSS secara live
css-watch:
	@if [ ! -f bin/tailwindcss ]; then go run cmd/make/main.go tailwind:install; fi
	bin/tailwindcss -i ./resources/css/input.css -o ./public/css/output.css --watch

# Jalankan Web development (Kompilasi template & CSS, lalu start server)
dev-web:
	@if [ ! -f bin/tailwindcss ]; then go run cmd/make/main.go tailwind:install; fi
	templ generate
	bin/tailwindcss -i ./resources/css/input.css -o ./public/css/output.css --minify
	@echo "🌐 Web server berjalan. Jalankan 'make css-watch' di tab terminal lain jika aktif mengubah styling UI."
	go run cmd/server/main.go

# Jalankan Asynq Background Job Worker
worker:
	go run cmd/worker/main.go

# Jalankan seluruh test suite dengan race detector
test:
	templ generate
	go test -v -race ./...

# Kompilasi template HTML a-h/templ
templ:
	templ generate

# Migrasi Database (Default: sqlite1)
migrate:
	go run cmd/make/main.go migrate:up --db=sqlite1

# Status Migrasi Database
migrate-status:
	go run cmd/make/main.go migrate:status --db=sqlite1

# Rollback Migrasi Terakhir
migrate-rollback:
	go run cmd/make/main.go migrate:down --db=sqlite1

# Database Seeding (Default: sqlite1)
seed:
	go run cmd/make/main.go db:seed --db=sqlite1

# Database Seeding MySQL saja (sqlite & db lain tidak disentuh)
seed-mysql:
	go run cmd/make/main.go db:seed --db=mysql1

# Database Seeding PostgreSQL saja
seed-pg:
	go run cmd/make/main.go db:seed --db=pg1

# Build single binary mandiri di direktori bin/
build:
	templ generate
	@if [ -f bin/tailwindcss ] && [ -f resources/css/input.css ]; then \
		bin/tailwindcss -i ./resources/css/input.css -o ./public/css/output.css --minify; \
	fi
	mkdir -p bin
	go build -o bin/server cmd/server/main.go
	go build -o bin/worker cmd/worker/main.go
	go build -o bin/make cmd/make/main.go

# Bersihkan build artifacts aplikasi (mempertahankan standalone binary seperti tailwindcss)
clean:
	rm -f bin/server bin/worker bin/make
