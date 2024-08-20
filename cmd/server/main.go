package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"PandoraHelper/cmd/server/wire"
	"PandoraHelper/pkg/config"
	"PandoraHelper/pkg/log"
	"go.uber.org/zap"

	// 导入 MySQL 驱动
	_ "github.com/go-sql-driver/mysql"
)

// @title           Pandora Helper API
// @version         1.0.0
// @description     This is the API server for Pandora Helper.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:9000
// @BasePath  /api/v1

// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	var envConf = flag.String("conf", "data/", "config path, eg: -conf ./config/local.yml")
	flag.Parse()

	conf := config.NewConfig(*envConf)
	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()

	if err != nil {
		logger.Fatal("Failed to initialize application", zap.Error(err))
	}

	pwd := conf.GetString("security.admin_password")
	if pwd == "" || len(pwd) < 8 {
		logger.Fatal("Invalid admin password: must be at least 8 characters long")
	}

	host := conf.GetString("http.host")
	port := conf.GetInt("http.port")

	logger.Info("Server starting", 
		zap.String("host", fmt.Sprintf("http://%s:%d", host, port)),
		zap.String("docs", fmt.Sprintf("http://%s:%d/swagger/index.html", host, port)),
	)

	if err = app.Run(context.Background()); err != nil {
		logger.Fatal("Server failed to run", zap.Error(err))
	}
}
