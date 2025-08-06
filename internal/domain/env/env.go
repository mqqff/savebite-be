package env

import (
	"github.com/mqqff/savebite-be/pkg/log"
	"github.com/spf13/viper"
	"time"
)

type Env struct {
	AppName string `mapstructure:"APP_NAME"`
	AppHost string `mapstructure:"APP_HOST"`
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`

	APIKey string `mapstructure:"API_KEY"`

	DBName string `mapstructure:"DB_NAME"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`

	JWTSecretKey   string        `mapstructure:"JWT_SECRET_KEY"`
	JWTExpiredTime time.Duration `mapstructure:"JWT_EXPIRED_TIME"`

	OAuthState         string `mapstructure:"OAUTH_STATE"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`

	SupabaseURL    string `mapstructure:"SUPABASE_URL"`
	SupabaseSecret string `mapstructure:"SUPABASE_SECRET"`

	GeminiAPIKey string `mapstructure:"GEMINI_API_KEY"`
	GeminiModel  string `mapstructure:"GEMINI_MODEL"`
}

var AppEnv = getEnv()

func getEnv() *Env {
	env := &Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[Env][getEnv] failed to read env file")
	}

	err = viper.Unmarshal(env)
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[Env][getEnv] failed to unmarshal")
	}

	return env
}
