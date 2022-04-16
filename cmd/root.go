package cmd

import (
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/Bipolar-Penguin/svc-notificator/pkg/service"
	"github.com/Bipolar-Penguin/svc-notificator/pkg/transport/amqp"
	"github.com/Bipolar-Penguin/svc-notificator/pkg/transport/http"
)

const (
	defaultHTTPort int = 8000
)

var (
	cfgHTTPPort                 int
	cfgRabbitmqURL              string
	cfgNotificatorTelegramToken string
	cfgNotificatorSMSToken      string
	cfgNotificatorEmailUsername string
	cfgNotificatorEmailPassword string
)

var rootCmd = &cobra.Command{
	Use:   "svc-notificator",
	Short: "Notificator microservice",
	Long:  "Svc-notificator is a microservice that sends notifications to users according to their strategy",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVar(&cfgHTTPPort, "http-port", defaultHTTPort, "HTTP port to listen")
	rootCmd.PersistentFlags().StringVar(&cfgRabbitmqURL, "rabbitmq-url", "", "RabbitMQ connection URL")
	rootCmd.PersistentFlags().StringVar(&cfgRabbitmqURL, "rabbitmq-excange", "", "RabbitMQ excange name")
	rootCmd.PersistentFlags().StringVar(&cfgNotificatorTelegramToken, "notificator-telegram-token", "", "Telegram bot API hash ID")
	rootCmd.PersistentFlags().StringVar(&cfgNotificatorSMSToken, "notificator-sms-token", "", "SMS Aero API hash ID")
	rootCmd.PersistentFlags().StringVar(&cfgNotificatorEmailUsername, "notificator-mail-user", "", "Yandex email")
	rootCmd.PersistentFlags().StringVar(&cfgNotificatorEmailPassword, "notificator-mail-password", "", "Yandex email password")

	viper.BindPFlag("http-port", rootCmd.PersistentFlags().Lookup("http-port"))
	viper.BindPFlag("rabbitmq-url", rootCmd.PersistentFlags().Lookup("rabbitmq-url"))
	viper.BindPFlag("rabbitmq-exchange", rootCmd.PersistentFlags().Lookup("rabbitmq-exchange"))
	viper.BindPFlag("notificator-telegram-token", rootCmd.PersistentFlags().Lookup("notificator-telegram-token"))
	viper.BindPFlag("notificator-sms-token", rootCmd.PersistentFlags().Lookup("notificator-sms-token"))
	viper.BindPFlag("notificator-mail-user", rootCmd.PersistentFlags().Lookup("notificator-mail-user"))
	viper.BindPFlag("notificator-mail-password", rootCmd.PersistentFlags().Lookup("notificator-mail-password"))
}

func initConfig() {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
}

func reader(conn *websocket.Conn) {
}

func run() {
	//var err error

	// Create the logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger.Log("app", os.Args[0], "event", "starting")
	}

	// Declare and validate the parameters
	httpPort := viper.GetInt("http-port")
	if httpPort == 0 {
		logger.Log("error", "http-port argument was not provided")
	}

	rabbitmqURL := viper.GetString("rabbitmq-url")
	if rabbitmqURL == "" {
		logger.Log("error", "rabbitmq-url argument was not provided")
		os.Exit(1)
	}

	rabbitmqExchangeName := viper.GetString("rabbitmq-exchange")
	if rabbitmqExchangeName == "" {
		logger.Log("error", "rabbitmq-url argument was not provided")
		os.Exit(1)
	}

	notificatorTelegramToken := viper.GetString("notificator-telegram-token")
	if notificatorTelegramToken == "" {
		logger.Log("error", "notificator-telegram-token argument was not provided")
		os.Exit(1)
	}

	notificatorSMSToken := viper.GetString("notificator-sms-token")
	if notificatorTelegramToken == "" {
		logger.Log("error", "notificator-sms-token argument was not provided")
		os.Exit(1)
	}

	notificatorEmailUsername := viper.GetString("notificator-mail-user")
	if notificatorEmailUsername == "" {
		logger.Log("error", "notificator-mail-user argument was not provided")
		os.Exit(1)
	}

	notificatorEmailPassword := viper.GetString("notificator-mail-password")
	if notificatorEmailPassword == "" {
		logger.Log("error", "notificator-mail-password argument was not provided")
		os.Exit(1)
	}

	// Create TG notificator
	telegramNotificator := service.NewTelegramNotificator(notificatorTelegramToken, logger)

	// Create SMS notificator
	smsNotificator := service.NewSMSNotificator(notificatorSMSToken, logger)

	// Create Email notificator
	mailNotificator := service.NewEmailNotificator(notificatorEmailUsername, notificatorEmailPassword, logger)

	// Create service
	service := service.NewService(*telegramNotificator, *smsNotificator, *mailNotificator, logger)

	tradingUpdates := make(chan domain.Event)

	// Create http websocket server
	websocketServer := http.NewWebsocketServer(httpPort, tradingUpdates, logger)
	websocketServer.StreamUpdatesToWebsocket()

	// Create AMQP listener
	var consumer amqp.Consumer
	{
		rabbitConfig := amqp.Config{
			ExchangeName: rabbitmqExchangeName,
			ExchangeType: "topic",
			ConnString:   rabbitmqURL,
			Durable:      true,
			AutoDelete:   false,
			Internal:     false,
			Args:         nil,
			AutoAck:      false,
			Exclusive:    false,
			NoLocal:      false,
			NoWait:       false,
		}

		broker := amqp.NewRabbitBroker(logger, &rabbitConfig, service, tradingUpdates)

		consumer = amqp.Consumer(broker)
		consumer.Consume()
	}
}
