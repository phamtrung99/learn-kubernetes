package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func GetMysqlClient(
	masters []string, workers []string, port string,
	userName string, dbname string, password string,
	maxIdleCons int, maxOpenCons int, logLevel string,
) *gorm.DB {
	var masterDialector gorm.Dialector

	masterDialectors := make([]gorm.Dialector, 0)
	workerDialectors := make([]gorm.Dialector, 0)
	dnsClient := "%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	for i, master := range masters {

		connectionString, err := truncateSprintf(dnsClient, userName, password, master, port, dbname)

		if err != nil {
			panic(err)
		}
		dialector := mysql.New(mysql.Config{DSN: connectionString})

		if i == 0 {
			masterDialector = dialector
		} else {
			masterDialectors = append(masterDialectors, dialector)
		}
	}

	for _, worker := range workers {
		if worker == "" {
			continue
		}

		connectionString, err := truncateSprintf(dnsClient, userName, password, worker, port, dbname)
		if err != nil {
			panic(err)
		}

		dialector := mysql.New(mysql.Config{DSN: connectionString})

		workerDialectors = append(workerDialectors, dialector)
	}

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		//Logger:                 logger.NewGorm(!cfg.Debug),
	}

	db, err := gorm.Open(masterDialector, gormConfig)
	if err != nil {
		panic(err)
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  masterDialectors,
		Replicas: workerDialectors,
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		panic(err)
	}

	if logLevel == "debug" {
		db = db.Debug()
	}

	rawDB, _ := db.DB()
	rawDB.SetMaxIdleConns(maxIdleCons)
	rawDB.SetMaxOpenConns(maxIdleCons)
	rawDB.SetConnMaxLifetime(time.Minute * 5)

	return db.Session(&gorm.Session{}).WithContext(context.TODO())
}

func truncateSprintf(str string, args ...interface{}) (string, error) {
	n := strings.Count(str, "%s")
	if n > len(args) {
		return "", errors.New("replace items more than args")
	}
	return fmt.Sprintf(str, args[:n]...), nil
}
