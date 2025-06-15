package main

import (
	"TwitchDonoCalculator/internal/db"
	"TwitchDonoCalculator/internal/discord"
	"TwitchDonoCalculator/internal/validator"
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

func (app *application) newArgsParser(args discord.DiscordMessageArgs, writer io.Writer) *flag.FlagSet {
	f := flag.NewFlagSet(args.CommandName, flag.ContinueOnError)
	f.SetOutput(writer) // Suppress output to avoid cluttering the console
	return f
}

func (app *application) DiscordGetAllDonationsByStreamer(args discord.DiscordMessageArgs, writer io.Writer) {
	f := app.newArgsParser(args, writer)

	from := f.String("from", "", "start from date format: YYYY-MM-DD")
	to := f.String("to", "", "end date format: YYYY-MM-DD")

	if err := f.Parse(args.Args); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(writer, "Error parsing arguments: %v\n", err)
		}
		return
	}

	var argsStruct struct {
		From time.Time
		To   time.Time
		validator.Validator
	}

	validator.HandleDateRange(&argsStruct.Validator, *from, *to, &argsStruct.From, &argsStruct.To)

	if !argsStruct.Valid() {
		fmt.Fprintln(writer, argsStruct.Error())
		return
	}

	var params db.GetSumDonationByStreamerParams
	params.FromTimestamp = argsStruct.From
	params.ToTimestamp = argsStruct.To

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

func (app *application) getLastDonationsFromStreamer(args discord.DiscordMessageArgs, writer io.Writer) {
	f := app.newArgsParser(args, writer)

	var argsStruct struct {
		From     time.Time
		To       time.Time
		Streamer string
		validator.Validator
	}

	from := f.String("from", "", "start from date format: YYYY-MM-DD")
	to := f.String("to", "", "end date format: YYYY-MM-DD")
	f.StringVar(&argsStruct.Streamer, "streamer", "", "Streamer name, if not empty, will use all streamers")

	if err := f.Parse(args.Args); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(writer, "Error parsing arguments: %v\n", err)
		}
		return
	}

	validator.HandleDateRange(&argsStruct.Validator, *from, *to, &argsStruct.From, &argsStruct.To)

	if !argsStruct.Valid() {
		fmt.Fprintln(writer, argsStruct.Error())
		return
	}

	res, err := app.db.GetLastDonationsByStreamer(context.Background(), db.ArgsGetLastDonationsByStreamer{
		Streamer: argsStruct.Streamer,
		From:     argsStruct.From,
		To:       argsStruct.To,
	})
	if err != nil {
		fmt.Fprintf(writer, "Error getting donations: %v\n", err)
		return
	}
	if len(res) == 0 {
		fmt.Fprintf(writer, "No donations found for streamer %s.\n", argsStruct.Streamer)
		return
	}

	db.MakeTable()
	
}

