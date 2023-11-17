package config

type ApiEnvConfig struct {
	// LOCAL, DEV, STG, PRD
	Env string
	// server traffic on this port
	Port string
	// path to VERSION file
	Version string
}

const DEV_ENV = "DEV"
const STAGE_ENV = "STAGE"
const PROD_ENV = "PROD"
