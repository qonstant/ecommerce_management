package app

import (
	"context"
	_ "ecommerce_management/docs"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"ecommerce_management/internal/config"
	"ecommerce_management/internal/database"
	"ecommerce_management/internal/handlers"
	"ecommerce_management/internal/provider/epay"
	"ecommerce_management/pkg/log"
	"ecommerce_management/pkg/server"
)

// Run initializes the whole application
func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.LoadConfig(".")
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	// Debug print to verify configuration values
	fmt.Printf("Loaded configuration: %+v\n", configs)
	fmt.Printf("EPAY Configuration: \n* URL: %s\n* Login: %s\n* Password: %s\n* OAuthURL: %s\n* PaymentPageURL: %s\n",
		configs.EPAYURL, configs.EPAYLogin, configs.EPAYPassword, configs.EPAYOAuthURL, configs.EPAYPaymentPageURL)

	database.InitDB()

	// Initialize the ePay client
	epayClient, err := epay.New(epay.Credentials{
		URL:            configs.EPAYURL,
		Login:          configs.EPAYLogin,
		Password:       configs.EPAYPassword,
		OAuthURL:       configs.EPAYOAuthURL,
		PaymentPageURL: configs.EPAYPaymentPageURL,
		ShopID:         configs.ShopID,
		TerminalID:     configs.TerminalID,
	})
	if err != nil {
		logger.Error("ERR_INIT_EPAY_CLIENT", zap.Error(err))
		return
	}

	handlers, err := handlers.New(
		handlers.Dependencies{
			DB:         database.DB,
			Configs:    configs,
			EpayClient: epayClient, // Add epayClient to dependencies
		},
		handlers.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.ServerAddress))
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}

	// Run our server in a goroutine so that it doesn't block
	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		return
	}
	logger.Info("http server started on http://localhost" + configs.ServerAddress + "/swagger/index.html")

	// Graceful Shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // Create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("gracefully shutting down...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the httpServer gracefully
	}

	fmt.Println("running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("server was successfully shutdown.")
}
