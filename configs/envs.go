package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PublicHost      string
	Port            string
	DBUser          string
	DBPassword      string
	DBAdress        string
	DBName          string
	TokenSecretWord string
}
type djamel struct {
	name  string
	phone int16
}

var Env = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:      getEnv("PUBLIC_HOST", "http://localhost"),
		Port:            getEnv("PORT", "8080"),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", "Waelbvbusmh007."),
		DBAdress:        fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:          getEnv("DB_NAME", "marquino"),
		TokenSecretWord: getEnv("Token_Secret_Word", "waelo"),
	}

}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
