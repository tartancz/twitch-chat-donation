package twitch

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TwitchDonationHandler struct {
	regFind *regexp.Regexp
	client  *TwitchReader
	BotName string
}

type Donation struct {
	User     string
	Channel  string
	SendFrom string
	Amount   int64
	Text     string
}

func NewTwitchDonationHandler(regPattern string, client *TwitchReader, BotName string) *TwitchDonationHandler {
	return &TwitchDonationHandler{
		regFind: regexp.MustCompile(regPattern),
		client:  client,
		BotName: BotName,
	}
}

func (t *TwitchDonationHandler) GetDonation() (*Donation, error) {
	return t.GetDonationWithContext(context.Background())
}

func (t *TwitchDonationHandler) GetDonationWithContext(ctx context.Context) (*Donation, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			line, err := t.client.ReadLineWithContext(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to read line: %w", err)
			}

			message, err := ParseMessage(line)
			if err != nil {
				if errors.Is(err, ErrorNotValidFormat) {
					continue
				}
				return nil, fmt.Errorf("failed to parse message: %w", err)
			}

			//fmt.Printf("%#v\n", message)

			if !strings.EqualFold(message.User, t.BotName) {
				continue
			}

			value := t.regFind.FindString(message.Text)
			if value == "" {
				continue
			}

			amount, err := t.parseAmount(value)
			if err != nil {
				continue
			}

			return &Donation{
				User:     message.User,
				Text:     message.Text,
				Channel:  message.Channel,
				SendFrom: strings.Split(message.Text, " ")[0],
				Amount:   amount,
			}, nil
		}
	}
}

func (t *TwitchDonationHandler) parseAmount(value string) (int64, error) {
	don, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse donation amount: %w", err)
	}
	// Convert to cents to avoid floating point precision issues
	return int64(don), nil
}

func (t *TwitchDonationHandler) Close() error {
	return t.client.Close()
}
