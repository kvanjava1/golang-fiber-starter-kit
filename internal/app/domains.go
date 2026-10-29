package app

import (
	"golang-starter-kit/internal/domain/example"
	"golang-starter-kit/internal/domain/health"
	"golang-starter-kit/internal/domain/user"
)

// RegisterAppDomains adalah tempat tunggal pendaftaran seluruh modul domain aplikasi.
// Setiap domain dirakit secara mandiri menggunakan provider NewModule() dalam 1 baris kode yang bersih.
func RegisterAppDomains(c *Container) {
	// 🩺 1. Domain Health Check & Observability
	healthCtrl := health.NewController(c.DB, nil)
	c.RegisterDomain(health.NewRouteHandler(healthCtrl))

	// Pastikan koneksi database utama aktif
	if c.DB == nil || c.DB.DefaultSQL() == nil {
		return
	}
	sqlDB := c.DB.DefaultSQL()

	// 📦 2. Katalog Modul Bisnis (Cukup 1 Baris per Modul):
	c.RegisterDomain(example.NewModule(sqlDB, c.Cache))
	c.RegisterDomain(user.NewModule(sqlDB, c.Cache))

	// Tambahkan modul baru Anda di sini secara rapi (contoh: product, category, order):
	// c.RegisterDomain(product.NewModule(sqlDB, c.Cache))
	// c.RegisterDomain(category.NewModule(sqlDB, c.Cache))
}
