package db

import (
	"context"
	"time"
)

type ArgsGetLastDonationsByStreamer struct {
	Streamer string
	Page     int64
	PageSize int64
	From     time.Time
	To       time.Time
}

// GetLastDonationsByStreamer returns the last donations made to a streamer.
// The donations are paginated, with page 0 being the most recent donations.
// streamer: the name of the streamer to get donations for if "" then all streamers are returned
//
//	page: the page to get
func (q *Queries) GetLastDonationsByStreamer(ctx context.Context, args ArgsGetLastDonationsByStreamer) (donations []Donation, err error) {
	if args.Page < 0 {
		args.Page = 0
	}
	if args.PageSize <= 0 || args.PageSize > 100 {
		args.PageSize = 20 // default page size
	}

	pageOffset := args.Page * args.PageSize

	if args.Streamer == "" {
		donations, err = q.getLastDonations(ctx, getLastDonationsParams{
			Offset: pageOffset,
			Limit:  args.PageSize,
		})
	} else {
		donations, err = q.getLastDonationsByStreamer(ctx, getLastDonationsByStreamerParams{
			Channel: args.Streamer,
			Offset:  pageOffset,
			Limit:   args.PageSize,
		})
	}
	return
}
