package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/phamtrung99/learn-kubernetes/web/utils"
)

var serverID string

type User struct {
	ID   string `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"name"`
}

func main() {
	cfg := utils.GetConfig()
	// Init db connection
	db := utils.GetMysqlClient(
		cfg.MySQL.Masters, cfg.MySQL.Workers, cfg.MySQL.Port,
		cfg.MySQL.User, cfg.MySQL.DBName, cfg.MySQL.Pass,
		cfg.MySQL.MaxIdleConns, cfg.MySQL.MaxOpenConns, "debug",
	)

	// get users in db
	users := make([]*User, 0)
	_ = db.Find(&users).Error

	// Init APIs
	rand.Seed(time.Now().UTC().UnixNano())
	logger := log.New(os.Stdout, "", log.Lmicroseconds)

	logger.Printf("Listening on port %s ...", cfg.Port)
	logger.Printf("Users: %v", users[0].Name)

	http.HandleFunc("/health", Middle(logger, Health))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil); err != nil {
		logger.Panicln(err)
	}
}

func Middle(l *log.Logger, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		f(w, r)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s: ok", serverID)
}
