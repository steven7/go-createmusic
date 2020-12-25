package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port     int    		`json:"port"`
	Env      string 		`json:"env"`
	Pepper   string 		`json:"pepper"`
	HMACKey  string 		`json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
	Dropbox  OAuthConfig    `json:"dropbox"`
}

type PostgresConfig struct {
	Host	 string `json:"host"`
	Port 	 int    `json:"port"`
	User	 string `json:"user"`
	Password string `json:"password"`
	Name 	 string `json:"name"`
}


func (c PostgresConfig) Dialect () string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	// We are going to provide two potential connection info
	// strings based on whether a password is present
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s " +
			"sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s " +
		"dbname=%s sslmode=disable", c.Host, c.Port, c.User,
		c.Password, c.Name)
}

func DefaultConfig() Config {
	return Config{
		Port:     5000,
		Env:      "dev",
		Pepper:   "secret-random-string",
		HMACKey:  "secret-hmac-key",
		// Database: DefaultPostgresConfig(),
		Database: DefaultPostgresConfig_Docker_Dev(),
	}
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "create_music",
	}
}
func DefaultPostgresConfig_Docker_Dev() PostgresConfig {
	return PostgresConfig{
		Host:     "go-createmusic_db_1",
		//Host:     "db",
		Port:     5432,
		//Name:     "createmusic_dev",
		Name:     "postgres",
		User:     "postgres",
		Password: "postgres",
	}
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

type OAuthConfig struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	AuthURL  string `json:"auth_url"`
	TokenURL string `json:"token_url"`
}





func (c Config) IsProd() bool {
	return c.Env == "prod"
}

//func DefaultConfigProd() Config {
//	return Config{
//		Port:     3000,
//		Env:      "prod",
//		Pepper:   "pepper",
//		HMACKey:  "hmac_key",
//		Database: DefaultPostgresConfigProd(),
//	}
//}



func LoadConfig(configReq bool) Config {
	// Open the config file
	// f, err := os.Open(".config")
	f, err := os.Open(".configdev")
	if err != nil {
		if configReq {
			panic(err)
		}
		// If there was an error opening the file, print out a
		// message saying we are using the default config and
		// return it.
		fmt.Println("Using the default config...")
		return DefaultConfig()
		// return DefaultConfigProd() // this is a stopgap for the moment
	}
	fmt.Println("Using the .config file...")
	// If we opened the config file successfully we are going
	// to create a Config variable to load it into.
	var c Config
	// We also need a JSON decoder, which will read from the
	// file we opened when decoding
	dec := json.NewDecoder(f)
	// We then decode the file and place the results in c, the
	// Config variable we created for the results. The decoder
	// knows how to decode the data because of the struct tags
	// (eg `json:"port"`) we added to our Config and
	// PostgresConfig fields, much like GORM uses struct tags
	// to know which database column each field maps to.
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	// If all goes well, return the loaded config.
	fmt.Println("Sucessfully loaded .config")
	return c
}

