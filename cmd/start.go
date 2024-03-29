package cmd

import (
	"context"
	"fmt"
	"github.com/danesparza/fxtrigger/internal/data"
	"github.com/danesparza/fxtrigger/internal/trigger"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danesparza/fxtrigger/api"
	_ "github.com/danesparza/fxtrigger/docs" // swagger docs location
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API and UI services",
	Long:  `Start the API and UI services`,
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Debug().Str("configFile", viper.ConfigFileUsed()).Msg("Using config file")
	} else {
		log.Debug().Msg("No config file found")
	}

	systemdb := viper.GetString("datastore.system")
	dndschedule := viper.GetString("trigger.dndschedule")
	dndstarttime := viper.GetString("trigger.dndstart")
	dndendtime := viper.GetString("trigger.dndend")

	//	Emit what we know:
	log.Info().
		Str("systemdb", systemdb).
		Str("dndschedule", dndschedule).
		Str("dndstarttime", dndstarttime).
		Str("dndendtime", dndendtime).
		Msg("Config")

	//	Create a DBManager object and associate with the api.Service
	db, err := data.NewManager(systemdb)
	if err != nil {
		log.Err(err).Msg("Problem trying to open the system database")
		return
	}
	defer db.Close()

	//	Create a background service object
	backgroundService := trigger.BackgroundProcess{
		FireTrigger:   make(chan data.Trigger),
		AddMonitor:    make(chan data.Trigger),
		RemoveMonitor: make(chan string),
		DB:            db,
	}

	//	Create an api service object
	apiService := api.Service{
		FireTrigger:   backgroundService.FireTrigger,
		AddMonitor:    backgroundService.AddMonitor,
		RemoveMonitor: backgroundService.RemoveMonitor,
		DB:            db,
		StartTime:     time.Now(),
	}

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

	//	Log that the system has started:
	log.Info().Msg("System started")

	//	Create a router and setup our REST endpoints...
	restRouter := mux.NewRouter()

	//	TRIGGER ROUTES
	restRouter.HandleFunc("/v1/triggers", apiService.CreateTrigger).Methods("POST")        // Create a trigger
	restRouter.HandleFunc("/v1/triggers", apiService.UpdateTrigger).Methods("PUT")         // Update a trigger
	restRouter.HandleFunc("/v1/triggers", apiService.ListAllTriggers).Methods("GET")       // List all triggers
	restRouter.HandleFunc("/v1/triggers/{id}", apiService.DeleteTrigger).Methods("DELETE") // Delete a trigger

	restRouter.HandleFunc("/v1/trigger/fire/{id}", apiService.FireSingleTrigger).Methods("POST") // Fire a trigger

	//	SWAGGER ROUTES
	restRouter.PathPrefix("/v1/swagger").Handler(httpSwagger.WrapHandler)

	//	Create background processes to
	//	- listen for triggers events
	//	- handle requests to fire a trigger:
	go backgroundService.ListenForEvents(ctx)
	go backgroundService.HandleAndProcess(ctx)

	//	Initialize monitoring
	backgroundService.InitializeMonitors()

	//	Setup the CORS options:
	log.Info().Str("CORS origins", viper.GetString("server.allowed-origins")).Msg("CORS config")

	uiCorsRouter := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("server.allowed-origins"), ","),
		AllowCredentials: true,
	}).Handler(restRouter)

	//	Format the bound interface:
	formattedServerPort := fmt.Sprintf(":%v", viper.GetString("server.port"))

	//	Start the service and display how to access it
	log.Info().Str("server", formattedServerPort).Msg("Started REST service")
	log.Err(http.ListenAndServe(formattedServerPort, uiCorsRouter)).Msg("HTTP API service error")
}

func handleSignals(ctx context.Context, sigs <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigs:
		switch sig {
		case os.Interrupt:
			log.Info().Msg("SIGINT")
		case syscall.SIGTERM:
			log.Info().Msg("SIGTERM")
		}

		log.Info().Msg("Shutting down ...")
		cancel()
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
