package main

import (
	"context"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"bitbucket.org/bexstech/temis-compliance/src/_init"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/prometheus"

	_ "bitbucket.org/bexstech/temis-compliance/src/docs"
	"bitbucket.org/bexstech/temis-compliance/src/infra"
	"bitbucket.org/bexstech/temis-compliance/src/repository"
)

const (
	sslMode                           = false
	startCommandComplianceAsync       = "ASYNC_COMPLIANCE"
	startCommandComplianceExternalApi = "EXTERNAL_COMPLIANCE_API"
	startCommandComplianceInternalApi = "INTERNAL_COMPLIANCE_API"
)

var (
	//DB
	host     = os.Getenv("DATABASE_HOST")
	dbPort   = os.Getenv("DATABASE_PORT")
	user     = os.Getenv("DATABASE_USERNAME")
	dbName   = os.Getenv("DATABASE_NAME")
	password = os.Getenv("DATABASE_PASSWORD")

	startCommand = os.Getenv("START_COMMAND")
)

// @title Temis Compliance API
// @version 1.0
// @description API Documentation for Temis Compliance Microservice.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes https
// @host api-dev.bexs.com.br
// @BasePath /compliance
func main() {
	_ = os.Setenv("TZ", "UTC")

	if !_init.PrettyLog() {
		formatter := &logrus.JSONFormatter{}
		formatter.TimestampFormat = "2006-01-02T15:04:05.999999999Z"
		logrus.SetFormatter(formatter)
	}

	logrus.SetOutput(&infra.LogWriter{})

	logrus.Infof("started Temis Compliance with %s command", startCommand)

	prometheus.ExposeMetrics()

	wg := new(sync.WaitGroup)

	ctx := context.Background()
	db := getDb()
	_init.InitPubsubClients(ctx)
	_init.InitWebClients()
	_init.InitRepositories(db)

	switch startCommand {
	case startCommandComplianceAsync:
		wg.Add(1)
		go _init.StartComplianceAsync(ctx, wg)
		break
	case startCommandComplianceExternalApi:
		wg.Add(1)
		go _init.StartComplianceApis(ctx, wg, true, db)
		break
	case startCommandComplianceInternalApi:
		wg.Add(1)
		go _init.StartComplianceApis(ctx, wg, false, db)
		break
	default:
		wg.Add(2)
		go _init.StartComplianceAsync(ctx, wg)
		go _init.StartComplianceApis(ctx, wg, false, db)
		break
	}
	wg.Wait()
}

func getDb() *gorm.DB {
	db, err := infra.OpenConnection(host, dbPort, user, dbName, password, sslMode)
	if err != nil {
		logrus.Fatal(errors.WithStack(err))
	}
	err = repository.Migrate(db)
	if err != nil {
		logrus.Fatal(errors.WithStack(err))
	}
	return db
}
