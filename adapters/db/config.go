package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	cardtoken "payment-layer-card-api/entities/card_token"
	"payment-layer-card-api/entities/cards"
	customertoken "payment-layer-card-api/entities/customer_token"
	"payment-layer-card-api/entities/customers"

	"github.com/adhfoundation/layer-tools/log"
	_ "github.com/lib/pq"

	database "github.com/adhfoundation/layer-tools/database/standard"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSL      string
}

func New(config Config) (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSL)

	descriptionMap := map[string]map[string]string{
		"customers":       customers.CustomerAttributes,
		"cards":           cards.CardsAttributes,
		"card_tokens":     cardtoken.CardTokensAttributes,
		"customer_tokens": customertoken.CustomerTokensAttributes,
	}

	db, err := database.NewDBWithHooks(conn, &descriptionMap, nil)
	if err != nil {
		return nil, err
	}

	maxOpenConnString := os.Getenv("MAX_OPEN_CONNECTIONS")
	maxOpenIddleConnString := os.Getenv("MAX_IDLE_CONNECTIONS")

	maxOpenConn, err := strconv.Atoi(maxOpenConnString)
	if err != nil {
		log.Fatal(context.Background()).Msg("Error converting MAX_OPEN_CONNECTIONS to int")
	}

	maxOpenIddleConn, err := strconv.Atoi(maxOpenIddleConnString)
	if err != nil {
		log.Fatal(context.Background()).Msg("Error converting MAX_IDLE_CONNECTIONS to int")
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxOpenIddleConn)

	log.Info(context.Background()).Msg("Banco de dados conectado com sucesso, DBname: " + config.DBName)

	return db, nil
}
