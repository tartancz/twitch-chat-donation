package main

import (
	"TwitchDonoCalculator/internal/config"
	"TwitchDonoCalculator/internal/db"
	"TwitchDonoCalculator/internal/service"
	"TwitchDonoCalculator/migrations"
	"context"
	"database/sql"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.OpenDB(cfg.DB)
	RunMigrations(database)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Initialize repositories and services (using auto-generated SQLC code)
	donationRepo := db.New(database)

	var wg sync.WaitGroup

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range len(cfg.Streamers) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			donationService := service.NewDonationService(donationRepo, cfg.Twitch, cfg.Streamers[i])

			// Start donation monitoring
			if err := donationService.StartMonitoring(ctx); err != nil {
				log.Fatalf("Failed to start monitoring: %v", err)
			}
		}()

	}

	wg.Wait()

	// // Handle shutdown signals
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	<-sigChan
	// 	log.Println("Shutting down gracefully...")
	// 	cancel()
	// }()

}

func RunMigrations(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("failed to create sqlite driver: %v", err)
	}

	source, err := iofs.New(migrations.MigrationsFiles, ".")
	if err != nil {
		log.Fatalf("failed to create source: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite3", driver)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("Migration applied successfully.")
}
