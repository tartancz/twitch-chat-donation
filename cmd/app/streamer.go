package main

import (
	"TwitchDonoCalculator/internal/config"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Streamer struct {
	ChannelName      string
	BotName           string
	RegFind           *regexp.Regexp
	LineFilterContain string
	LogMessage        bool
	LogFile           *os.File
}

func NewStreamer(streamerConfig config.StreamerConfig, channelName string) *Streamer {
	return &Streamer{
		BotName:           streamerConfig.BotName,
		RegFind:           regexp.MustCompile(streamerConfig.ValueRegex),
		LineFilterContain: streamerConfig.LineFilterContain,
		LogMessage:        streamerConfig.LogMessage,
		ChannelName:       channelName,
	}
}

func NewStreamersFromMap(streamers map[string]*config.StreamerConfig) map[string]*Streamer {
	streamersMap := make(map[string]*Streamer)
	for k, v := range streamers {
		streamersMap[k] = NewStreamer(*v, k)
	}
	return streamersMap
}

func (s *Streamer) FindDonation(message string) int64 {
	if !strings.Contains(message, s.LineFilterContain) {
		return 0
	}

	value := s.RegFind.FindString(message)
	if value == "" {
		return 0
	}
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return int64(amount)
}
