package main

import (
	"TwitchDonoCalculator/internal/discord"
	"TwitchDonoCalculator/internal/twitch"
	"fmt"
	"os"
	"path"
)

func (app *application) GetStreamer(streamer string) *Streamer {
	return app.streamers[streamer]
}

func (app *application) LogStreamerMessage(message twitch.Message, streamer *Streamer) {
	if !streamer.LogMessage {
		return
	}

	//check if folder exists
	app.CreateLogFolder()
	if streamer.LogFile == nil {
		if file, err := app.CreateLogFile(streamer.ChannelName); err != nil {
			fmt.Fprintf(discord.DefaultServer,"ERROR when creating log file: %s", err)
			fmt.Println(err)
			return
		} else {
			streamer.LogFile = file
		}
	}
	fmt.Fprintln(streamer.LogFile, message.GetRaw())
}

func (app *application) LogUnknownMessage(RawMessage string) {
	if !app.cfg.LogAll {
		return
	}
	app.CreateLogFolder()
	if app.unknowLogFile == nil {
		if file, err := app.CreateLogFile("unknown"); err != nil {
			fmt.Fprintf(discord.DefaultServer,"ERROR when creating log file: %s", err)
			return
		} else {
			app.unknowLogFile = file
		}
	}
	fmt.Fprintln(app.unknowLogFile, RawMessage)

}

func (app *application) CreateLogFolder() {
	if _, err := os.Stat(app.cfg.LogFolder); os.IsNotExist(err) {
		os.Mkdir(app.cfg.LogFolder, os.ModePerm)
	}
}

func (app *application) CreateLogFile(fileName string) (*os.File, error) {
	return os.OpenFile(
		path.Join(app.cfg.LogFolder, fmt.Sprintf("%s.log", fileName)),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
}

func (app *application) CloseLogFiles() {
	for _, streamer := range app.streamers {
		if streamer.LogFile != nil {
			streamer.LogFile.Close()
		}
	}
	if app.unknowLogFile != nil {
		app.unknowLogFile.Close()
	}
}
