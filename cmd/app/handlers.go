package main

import (
	"TwitchDonoCalculator/internal/db"
	"TwitchDonoCalculator/internal/discord"
	"TwitchDonoCalculator/internal/twitch"
	"context"
	"fmt"
	"strings"
)

const donationLimitNotification = 10_000

func (app *application) HandleAnyMessage(m twitch.Message) {
	app.LogAnyMessage(m.GetRaw())
}



func (app *application) HandleChatMessage(m *twitch.MessagePrivate) {
	streamer := app.streamers[m.Streamer]
	app.LogStreamerMessage(m, streamer)
	if streamer.BotName != m.Sender {
		return
	}

	value := streamer.FindDonation(m.Text)
	if value == 0 {
		return
	}
	if value >= donationLimitNotification {
		fmt.Fprintf(discord.DefaultServer, "%s just got  %d donation", m.Streamer, value)
	}

	app.db.CreateDonation(context.Background(), db.CreateDonationParams{
		User:     m.Sender,
		Channel:  m.Streamer,
		SendFrom: strings.Split(m.Text, " ")[0],
		Amount:   value,
		Text:     m.Text,
	})

}

func (app *application) HandleChatNotice(m *twitch.MessageNotice) {
	streamer := app.streamers[m.Streamer]
	value := streamer.FindDonation(m.Text)
	if value == 0 {
		return
	}
	if value >= donationLimitNotification {
		fmt.Fprintf(discord.DefaultServer, "%s just got  %d donation", m.Streamer, value)
	}
	app.db.CreateDonation(context.Background(), db.CreateDonationParams{
		User:     "",
		Channel:  m.Streamer,
		SendFrom: strings.Split(m.Text, " ")[0],
		Amount:   value,
		Text:     m.Text,
	})
}

func (app *application) HandleUnknowMessage(m *twitch.UnknowMessage) {
	app.LogUnknownMessage(m.Raw)
}
