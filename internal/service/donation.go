package service

import (
	"TwitchDonoCalculator/internal/config"
	"TwitchDonoCalculator/internal/db"
	"TwitchDonoCalculator/internal/twitch"
	"context"
	"log"
	"time"
)

type DonationService struct {
	repo     *db.Queries
	config   config.TwitchConfig
	streamer config.StreamerConfig
}

func NewDonationService(repo *db.Queries, config config.TwitchConfig, streamer config.StreamerConfig) *DonationService {
	return &DonationService{
		repo:     repo,
		config:   config,
		streamer: streamer,
	}
}

func (s *DonationService) StartMonitoring(ctx context.Context) error {
	auth := twitch.TwitchAuth{
		Oauth: s.config.OAuth,
		Nick:  s.config.Nick,
	}
	t, err := twitch.NewTwitchReader(
		auth,
		s.streamer.ChannelName,
	)
	if err != nil {
		return err
	}
	handler := twitch.NewTwitchDonationHandler(
		s.streamer.ValueRegex,
		t,
		s.streamer.BotName,
	)

	log.Printf("Starting donation monitoring for channel: %s", s.streamer.ChannelName)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping donation monitoring")
			return ctx.Err()
		default:
			donation, err := handler.GetDonation()
			if err != nil {
				log.Printf("Error getting donation: %v", err)
				time.Sleep(time.Second) // Brief pause before retrying
				continue
			}

			if err := s.saveDonation(ctx, donation); err != nil {
				log.Printf("Error saving donation: %v", err)
			} else {
				log.Printf("Saved donation: %+v", donation)
			}
		}
	}
}

func (s *DonationService) saveDonation(ctx context.Context, donation *twitch.Donation) error {
	params := db.CreateDonationParams{
		User:     donation.User,
		Channel:  donation.Channel,
		SendFrom: donation.SendFrom,
		Amount:   donation.Amount,
		Text:     donation.Text,
	}

	_, err := s.repo.CreateDonation(ctx, params)
	return err
}
