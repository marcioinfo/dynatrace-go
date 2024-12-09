package main

import (
	"context"
	"database/sql"
	"os"
	postgres "payment-layer-card-api/adapters/db"
	"payment-layer-card-api/cmd/internal/factories"
	"payment-layer-card-api/common/variables"

	"github.com/adhfoundation/layer-tools/apmtracer"
	"github.com/adhfoundation/layer-tools/log"
	validateEnvironment "github.com/adhfoundation/layer-tools/validate-environments"
	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	return os.Getenv(key)
}

func buildEnvMap(keys []string) map[string]string {
	envs := make(map[string]string)
	for _, key := range keys {
		envs[key] = getEnvVariable(key)
	}
	return envs
}

func main() {
	ctx := context.Background()

	godotenv.Load()
	log.Init(ctx)

	envs := buildEnvMap(variables.EnvKeys)

	err_in_validate_envs := validateEnvironment.ValidateEnvVariables(envs)
	if err_in_validate_envs != nil {
		log.Fatal(ctx).Msg("Erro ao validar varíaveis de ambiente: " + err_in_validate_envs.Error())
		os.Exit(1)
	}

	err := apmtracer.Init(ctx)
	if err != nil {
		log.Fatal(ctx).Msg("Erro ao inicializar o ApmTracer: " + err.Error())
		os.Exit(1)
	}

	db, err := connectToDB()
	if err != nil {
		log.Fatal(ctx).Msg("Erro ao conectar ao banco de dados: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	httpServer := factories.NewServerFactory(db)
	httpServer.Start()
}

func connectToDB() (*sql.DB, error) {

	dbConfig := postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSL:      os.Getenv("DB_SSL"),
	}

	return postgres.New(dbConfig)
}
