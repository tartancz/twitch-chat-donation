package main

import (
	"TwitchDonoCalculator/internal/db"
	"TwitchDonoCalculator/internal/discord"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

func (app *application) registerDiscordCommands() {
	discord.DefaultServer.AddHandler("donation", discord.DiscordMessageHandlerStruct{
		HandleFunc:  app.DiscordGetAllDonationsByStreamer,
		HelpMessage: "Get all donations by streamer within a date range.",
	})
}

func (app *application) DiscordGetAllDonationsByStreamer(args discord.DiscordMessageArgs, writer io.Writer) {
	f := flag.NewFlagSet(args.CommandName, flag.ContinueOnError)
	f.SetOutput(writer)

	from := f.String("from", "", "start from date format: YYYY-MM-DD")
	to := f.String("to", "", "end date format: YYYY-MM-DD")
	fmt.Println(*from, *to)
	if err := f.Parse(args.Args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
	}

	var params db.GetSumDonationByStreamerParams

	if *from == "" {
		params.FromTimestamp = time.Time{}
	} else {
		fromTime, err := time.Parse(time.DateOnly, *from)
		if err != nil {
			fmt.Fprintf(writer, "Invalid to date format. Use 'YYYY-MM-DD' format.\n")
			return
		}
		params.FromTimestamp = fromTime
	}

	if *to == "" {
		params.ToTimestamp = time.Now()
	} else {
		toTime, err := time.Parse(time.DateOnly, *to)
		if err != nil {
			fmt.Fprintf(writer, "Invalid to date format. Use 'YYYY-MM-DD' format.\n")
			return
		}
		params.ToTimestamp = toTime
	}

	res, err := app.db.GetSumDonationByStreamer(context.Background(), params)
	if err != nil {
		fmt.Fprintf(writer, "Error getting donations: %v\n", err)
		return
	}
	if len(res) == 0 {
		fmt.Fprintf(writer, "No donations found.\n")
		return
	}

	buf := &bytes.Buffer{}

	tb := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tb, "Channel\tAmount\tStartingDate\tEndingDate")

	for _, r := range res {
		fmt.Fprintf(tb, "%s\t%d\t%s\t%s\n", r.Channel, r.Amount, r.Startingdate, r.Endingdate)
	}

	tb.Flush()
	fmt.Fprintf(writer, "```%s```", buf.String())

}
