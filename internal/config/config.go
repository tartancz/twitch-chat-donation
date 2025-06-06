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
	Env               string
	LogFolder         string
	LogAll            bool
	LogUnknownMessage bool
	DB                DBConfig
	Twitch            TwitchConfig
	Streamers         map[string]*StreamerConfig
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
	ValueRegex        string
	LineFilterContain string
	LogMessage        bool
}

func Load() *Config {
	return &Config{
		Env:               getEnv("ENV", "development"),
		LogFolder:         getEnv("LOG_FOLDER", "./logs/"),
		LogAll:            getEnvBool("LOG_ALL", true),
		LogUnknownMessage: getEnvBool("LOG_UNKNOWN_MESSAGE", true),
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

func GetStreamersConfig() map[string]*StreamerConfig {
	file, err := os.Open(JsonFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			streamers := GetNewStreamerConfig()
			CreateStreamersFile(streamers)
		}
	}
	defer file.Close()
	var streamers map[string]*StreamerConfig
	err = json.NewDecoder(file).Decode(&streamers)
	if err != nil {
		panic(err)
	}
	return streamers
}

func GetNewStreamerConfig() map[string]*StreamerConfig {
	streamerConfig := &StreamerConfig{}
	fmt.Print("Enter bot name to watch in Twitch (e.g. streamelements): ")
	fmt.Scanln(&streamerConfig.BotName)
	fmt.Print("Enter channel name to watch: ")
	var channelName string
	fmt.Scanln(&channelName)
	fmt.Print("Enter string chat message should contain: ")
	fmt.Scanln(&streamerConfig.LineFilterContain)
	fmt.Print("Enter Regex for float number: ")
	fmt.Scanln(&streamerConfig.ValueRegex)

	streamerConfigs := make(map[string]*StreamerConfig)
	streamerConfigs[channelName] = streamerConfig
	return streamerConfigs
}

func CreateStreamersFile(s map[string]*StreamerConfig) error {
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

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
