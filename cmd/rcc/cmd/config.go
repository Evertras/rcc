package cmd

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var config struct {
	DynamoDBTable      string `mapstructure:"dynamodb-table"`
	FileStorageBaseDir string `mapstructure:"file-storage-base-dir"`
	StorageType        string `mapstructure:"storage-type"`
}

func initConfig() {
	viper.SetEnvPrefix("RCC")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("Failed to unmarshal config: %s", err.Error())
	}
}
