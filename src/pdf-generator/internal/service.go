package internal

import (
	"fmt"
	"github.com/Limpid-LLC/saiService"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"strconv"
)

type InternalService struct {
	Context   *saiService.Context
	Logger    *zap.Logger
	AwsConfig *aws.Config
	FileNum   int //file slots from config
}

func (is *InternalService) Init() {
	is.SetLogger()
	is.awsConfig()

	if is.Context.GetConfig("file_server.enabled", false).(bool) {
		fileNum := is.Context.GetConfig("file_num", 50).(int)
		is.FileNum = fileNum

		fileServerPort := is.Context.GetConfig("file_server.port", "8083").(int)
		fileServerHandler := http.FileServer(http.Dir("./files"))
		fsMux := http.NewServeMux()
		fsMux.Handle("/", fileServerHandler)

		go func() {
			is.Logger.Debug("FileServer started", zap.String("directory", "files"), zap.Int("port", fileServerPort), zap.Int("file slots", fileNum))
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(fileServerPort)), fsMux))
		}()
	}
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
