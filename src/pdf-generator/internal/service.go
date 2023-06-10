package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Limpid-LLC/saiService"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type InternalService struct {
	Context   *saiService.Context
	Logger    *zap.Logger
	AwsConfig *aws.Config
	FileNum   int
}

func (is *InternalService) Init() {
	is.SetLogger()
	is.awsConfig()

	fileNum := is.Context.GetConfig("file_num", 50).(int)

	fileserverPort := is.Context.GetConfig("common.http.fileserver_port", "8083").(int)

	fileserverHandler := http.FileServer(http.Dir("./files"))
	fsMux := http.NewServeMux()
	fsMux.Handle("/", fileserverHandler)

	is.Logger.Debug("Fileserver started", zap.String("directory", "files"), zap.Int("port", fileserverPort), zap.Int("file slots", fileNum))

	go http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(fileserverPort)), fsMux)

}

// SetLogger set service logger
func (is *InternalService) SetLogger() {

	var (
		logger *zap.Logger
		err    error
	)
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	mode := is.Context.GetConfig("common.log_mode", "debug").(string)
	if mode == "debug" {
		option := zap.AddStacktrace(zap.DPanicLevel)
		logger, err = config.Build(option)
		if err != nil {
			log.Fatal("error creating logger : ", err.Error())
		}
		logger.Debug("Logger started", zap.String("mode", "debug"))
	} else {
		option := zap.AddStacktrace(zap.DPanicLevel)
		logger, err = config.Build(option)
		if err != nil {
			log.Fatal("error creating logger : ", err.Error())
		}
		logger.Info("Logger started", zap.String("mode", "production"))
	}

	is.Logger = logger
}
