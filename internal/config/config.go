package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const JsonFilePath = "./streamers.json"

type Config struct {
	Env       string
	DB        DBConfig
	Twitch    TwitchConfig
	Streamers []StreamerConfig
}

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type TwitchConfig struct {
	OAuth string
	Nick  string
}

type StreamerConfig struct {
	BotName           string
	ChannelName       string
	ValueRegex        string
	LineFilterContain string
}

func Load() *Config {
	return &Config{
		Env: getEnv("ENV", "development"),
		DB: DBConfig{
			DSN:          getEnv("DB_DSN", "db.db"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 50),
			MaxIdleTime:  getEnvDuration("DB_MAX_IDLE_TIME", time.Minute*15),
		},
		Twitch: TwitchConfig{
			OAuth: getEnv("TWITCH_OAUTH", ""),
			Nick:  getEnv("TWITCH_NICK", ""),
		},
		Streamers: GetStreamersConfig(),
	}
}

func GetStreamersConfig() []StreamerConfig {
	file, err := os.Open(JsonFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			streamer := GetNewStreamerConfig()
			CreateStreamersFile([]StreamerConfig{*streamer})
		}
	}
	defer file.Close()
	var streamers []StreamerConfig
	err = json.NewDecoder(file).Decode(&streamers)
	if err != nil {
		panic(err)
	}
	return streamers
}

func GetNewStreamerConfig() *StreamerConfig {
	var s StreamerConfig
	fmt.Println("Insert bot name that should be watched in twitch (most probably streamelements):")
	fmt.Scanln(&s.BotName)
	fmt.Println("Insert channel name to be watched:")
	fmt.Scanln(&s.ChannelName)
	fmt.Println("Insert string that chat message should contain")
	fmt.Scanln(&s.LineFilterContain)
	fmt.Println("Insert Regex for float number")
	fmt.Scanln(&s.ValueRegex)
	return &s
}

func CreateStreamersFile(s []StreamerConfig) error {
	file, err := os.Create(JsonFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	return enc.Encode(s)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
