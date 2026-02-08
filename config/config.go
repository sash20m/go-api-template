package config

import (
	"flag"
	"go-api-template/internal/libs/utils"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type EnvType string

const DEV_ENV EnvType = "DEV"
const PROD_ENV EnvType = "PROD"

// CONFIG contains all the values in .env.
var CONFIG ConfigSchema = ConfigSchema{}

const DefaultEnvFilename = ".env"

// ConfigSchema should have the same structure as .env.
type ConfigSchema struct {
	Env     EnvType
	Port    string
	Version string

	AllowedOrigins []string

	DevBaseURL  string
	ProdBaseURL string

	// Database
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string

	// RabbitMQ
	RabbitMQEnabled  bool
	RabbitMQURL      string
	RabbitMQPrefetch int

	// OpenAI
	OpenaiApiKey string
	// Anthropic
	AnthropicKey     string
	AnthropicVersion string
	AnthropicApiUrl  string
	// Google
	GeminiApiKey string
	// Groq
	GroqApiKey string

	// JWT
	JwtSecret             string
	JwtRefreshTokenSecret string
	TokenExpiration       string

	// Mailgun
	MailgunDomain string
	MailgunApiKey string

	// Emails
	EmailHello string

	// Stripe
	StripeKey            string
	StripeEndpointSecret string

	// Google
	GoogleApi          string
	GoogleClientId     string
	GoogleClientSecret string

	// Sentry
	SentryDsn string

	// GCP
	GcpCredentialsToken string
}

func LoadConfig() (ConfigSchema, error) {
	env := flag.String("envfilename", DefaultEnvFilename, "environment")
	flag.Parse()

	err := godotenv.Load(*env)
	if err != nil {
		logrus.WithError(err).Error("Can't load config from .env. Problem with .env, or the server is in production environment.")
	}

	config := ConfigSchema{
		Env:         EnvType(strings.ToUpper(os.Getenv("ENV"))),
		Port:        os.Getenv("PORT"),
		Version:     os.Getenv("VERSION"),
		DevBaseURL:  os.Getenv("DEV_BASE_URL"),
		ProdBaseURL: os.Getenv("PROD_BASE_URL"),

		AllowedOrigins: utils.ParseCsvLine(os.Getenv("ALLOWED_ORIGINS")),

		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),

		RabbitMQEnabled:  strings.ToLower(os.Getenv("RABBITMQ_ENABLED")) == "true",
		RabbitMQURL:      os.Getenv("RABBITMQ_URL"),
		RabbitMQPrefetch: utils.ParseIntWithDefault(os.Getenv("RABBITMQ_PREFETCH"), 20),

		OpenaiApiKey:     os.Getenv("OPENAI_API_KEY"),
		AnthropicKey:     os.Getenv("ANTHROPIC_KEY"),
		AnthropicVersion: os.Getenv("ANTHROPIC_VERSION"),
		AnthropicApiUrl:  os.Getenv("ANTHROPIC_API_URL"),
		GeminiApiKey:     os.Getenv("GEMINI_API_KEY"),
		GroqApiKey:       os.Getenv("GROQ_API_KEY"),

		JwtSecret:             os.Getenv("JWT_SECRET"),
		JwtRefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
		TokenExpiration:       os.Getenv("TOKEN_EXPIRATION"),

		MailgunDomain: os.Getenv("MAILGUN_DOMAIN"),
		MailgunApiKey: os.Getenv("MAILGUN_API_KEY"),
		EmailHello:    os.Getenv("EMAIL_HELLO"),

		StripeKey:            os.Getenv("STRIPE_KEY"),
		StripeEndpointSecret: os.Getenv("STRIPE_ENDPOINT_SECRET"),

		GoogleApi:          os.Getenv("GOOGLE_API"),
		GoogleClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),

		SentryDsn: os.Getenv("SENTRY_DSN"),

		GcpCredentialsToken: os.Getenv("GCP_CREDENTIALS_TOKEN"),
	}

	// All config variables in the app can be accessed by config.CONFIG anywhere in the app.
	CONFIG = config

	return config, nil
}
