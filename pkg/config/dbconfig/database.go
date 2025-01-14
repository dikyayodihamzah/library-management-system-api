package dbconfig

import (
	"context"
	"time"

	"github.com/dikyayodihamzah/library-management-api/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	host     = utils.GetString("DB_HOST")
	port     = utils.GetString("DB_PORT")
	username = utils.GetString("DB_USERNAME")
	password = utils.GetString("DB_PASSWORD")
	name     = utils.GetString("DB_NAME")
	minConns = utils.GetInt("DB_POOL_MIN")
	maxConns = utils.GetInt("DB_POOL_MAX")
	timeout  = utils.GetInt("DB_TIMEOUT")

	logger = utils.NewLogger()
)

// QuestDB
var (
	questDBHost     = utils.GetString("QUEST_DB_HOST")
	questDBPort     = utils.GetString("QUEST_DB_PORT")
	questDBUser     = utils.GetString("QUEST_DB_USERNAME")
	questDBPassword = utils.GetString("QUEST_DB_PASSWORD")
	questDBName     = utils.GetString("QUEST_DB_NAME")
)

type dbConfig struct {
	username string
	password string
	host     string
	port     string
	name     string
}

func NewDBConfig(dbType string) *dbConfig {
	switch dbType {
	case "quest":
		return &dbConfig{
			username: questDBUser,
			password: questDBPassword,
			host:     questDBHost,
			port:     questDBPort,
			name:     questDBName,
		}
	default:
		if dbType != "" {
			name = dbType
		}
		return &dbConfig{
			username: username,
			password: password,
			host:     host,
			port:     port,
			name:     name,
		}
	}
}

func NewPool(dbType ...string) *pgxpool.Pool {
	// define connection string
	if len(dbType) == 0 {
		dbType = append(dbType, "")
	}

	// db := NewDBConfig(dbType[0])
	dsn := utils.GetString("DB_URL")
	utils.Debug(dsn)

	// create pgxPool config object
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatalw("Failed to parse configuration", "dsn", dsn)
	}

	config.MinConns = int32(minConns)
	config.MaxConns = int32(maxConns)

	// set pool configuration
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatalw("Failed to apply pool configuration", "error", err)
	}

	c, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// ping database
	if err := pool.Ping(c); err != nil {
		logger.Fatalw("Failed to ping database", "error", err)
	}

	logger.Infow("Database connected", "dsn", dsn)
	return pool
}
