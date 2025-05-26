package main

import (
	"TwitchDonoCalculator/internal/config"
	"TwitchDonoCalculator/internal/db"
	"TwitchDonoCalculator/internal/discord"
	"TwitchDonoCalculator/internal/twitch"
	"log"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type application struct {
	db            *db.Queries
	streamers     map[string]*Streamer
	cfg           *config.Config
	logger        *slog.Logger
	unknowLogFile *os.File
}

func main() {
	go discord.DefaultServer.RunServer("TwitchDonoCalculator")
	defer discord.DefaultServer.Close()
	cfg := config.Load()

	// Initialize database
	database, err := db.OpenDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()
	db.RunMigrations(database)

	c := twitch.NewAnonymousClient()

	var app = &application{
		db:        db.New(database),
		streamers: NewStreamersFromMap(cfg.Streamers),
		cfg:       cfg,
		logger:    slog.Default(),
	}
	defer app.CloseLogFiles()

	c.SetOnChatMessage(app.HandleChatMessage)
	c.SetOnChatNotice(app.HandleChatNotice)
	c.SetOnAnyMessage(app.HandleAnyMessage)
	c.SetOnUnknowMessage(app.HandleUnknowMessage)

	app.registerDiscordCommands()

	for k := range cfg.Streamers {
		c.AddStreamers(k)
	}

	err = c.Listen()
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
}
