package infra

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Mysql struct {
	MySQL *sql.DB
}

func newMysql(viper *viper.Viper) *Mysql {

	config := viper.Sub("mysql")
	host := config.GetString("host")
	port := config.GetInt("port")
	username := config.GetString("username")
	password := config.GetString("password")
	database := config.GetString("database")
	// connection pool
	maxIdleLifeTime := config.GetInt("conn_max_life_time")
	maxIdleTime := config.GetInt("conn_max_idle_time")
	maxIdle := config.GetInt("max_idle_conn")
	maxOpen := config.GetInt("max_open_conn")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host, port, database)
	db, err := sql.Open("mysql", dsn)

	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)
	db.SetConnMaxLifetime(time.Duration(maxIdleLifeTime) * time.Minute)
	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)

	if err != nil {
		panic(err)
	}

	return &Mysql{
		MySQL: db,
	}

}
