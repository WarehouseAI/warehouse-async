package dependencies

import (
	"crypto/tls"
	"os"
	"os/signal"
	"syscall"

	randomAdpt "github.com/warehouse/mail-service/internal/adapter/random"
	timeAdpt "github.com/warehouse/mail-service/internal/adapter/time"
	internalBroker "github.com/warehouse/mail-service/internal/broker"
	"github.com/warehouse/mail-service/internal/config"
	"github.com/warehouse/mail-service/internal/handler/broker"
	"github.com/warehouse/mail-service/internal/pkg/logger"
	"github.com/warehouse/mail-service/internal/server"
	mailSvc "github.com/warehouse/mail-service/internal/service/mail"
	"gopkg.in/gomail.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Dependencies interface {
		Close()
		Cfg() *config.Config
		Internal() dependencies
		WaitForInterrupr()

		BrokerServer() server.Server
	}

	dependencies struct {
		cfg        *config.Config
		log        logger.Logger
		mailDialer *gomail.Dialer

		rabbitClient *internalBroker.RabbitClient

		mailHandler  broker.MailHandler
		mailService  mailSvc.Service
		brokerServer server.Server

		timeAdapter   timeAdpt.Adapter
		randomAdapter randomAdpt.Adapter

		shutdownChannel chan os.Signal
		closeCallbacks  []func()
	}
)

func NewDependencies(cfgPath string) (Dependencies, error) {
	cfg, err := config.NewConfig(cfgPath)
	if err != nil && err.Error() == "Config File \"config\" Not Found in \"[]\"" {
		cfg, err = config.NewConfig("./configs/local")
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "1"
	encoderCfg.TimeKey = "t"

	z := zap.New(
		&logger.WarehouseZapCore{
			Core: zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.Lock(os.Stdout),
				zap.NewAtomicLevel(),
			),
		},
		zap.AddCaller(),
	)

	return &dependencies{
		cfg:             cfg,
		log:             logger.NewLogger(z),
		shutdownChannel: make(chan os.Signal),
	}, nil
}

func (d *dependencies) Close() {
	for i := len(d.closeCallbacks) - 1; i >= 0; i-- {
		d.closeCallbacks[i]()
	}
	d.log.Zap().Sync()
}

func (d *dependencies) Internal() dependencies {
	return *d
}

func (d *dependencies) Cfg() *config.Config {
	return d.cfg
}

func (d *dependencies) MailDialer() *gomail.Dialer {
	if d.mailDialer == nil {
		d.mailDialer = gomail.NewDialer(
			d.Cfg().Mail.Host,
			d.Cfg().Mail.Port,
			d.Cfg().Mail.Email,
			d.Cfg().Mail.Password,
		)
		d.mailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return d.mailDialer
}

func (d *dependencies) BrokerServer() server.Server {
	if d.brokerServer == nil {
		var err error
		msg := "initialize broker server"
		if d.brokerServer, err = server.NewBrokerServer(
			d.log,
			*d.cfg,
			d.RabbitClient(),
			d.MailHandler(),
		); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}

		d.closeCallbacks = append(d.closeCallbacks, func() {
			msg := "shutting down broker server"
			if err := d.brokerServer.Stop(); err != nil {
				d.log.Zap().Warn(msg, zap.Error(err))
				return
			}
			d.log.Zap().Info(msg)
		})
	}

	return d.brokerServer
}

func (d *dependencies) WaitForInterrupr() {
	signal.Notify(d.shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	d.log.Zap().Info("Wait for receive interrupt signal")
	<-d.shutdownChannel // ждем когда сигнал запишется в канал и сразу убираем его, значит, что сигнал получен
	d.log.Zap().Info("Receive interrupt signal")
}
