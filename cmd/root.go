package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/hashicorp/logutils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile               string
	problemWithConfigFile bool
	loglevel              string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fxtrigger",
	Short: "REST service for GPIO / Sensor triggers",
	Long:  `REST based management for GPIO/Sensor -> endpoint triggers`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/fxtrigger.yaml)")
	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "INFO", "Log level: DEBUG/INFO/WARN/ERROR")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)        // adding home directory as first search path
		viper.AddConfigPath(".")         // also look in the working directory
		viper.SetConfigName("fxtrigger") // name the config file (without extension)
	}

	viper.AutomaticEnv() // read in environment variables that match

	//	Set our defaults
	viper.SetDefault("loglevel", "INFO")
	viper.SetDefault("datastore.system", path.Join(home, "fxtrigger", "db", "system.db"))
	viper.SetDefault("datastore.retentiondays", 30)
	viper.SetDefault("trigger.dndschedule", false) //	Use a 'Do not disturb' schedule
	viper.SetDefault("trigger.dndstart", "8:00pm") //	Do not disturb scheduled start time
	viper.SetDefault("trigger.dndend", "6:00am")   //	Do not disturb scheduled end time
	viper.SetDefault("server.port", 3020)
	viper.SetDefault("server.allowed-origins", "*")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		problemWithConfigFile = true
	}

	//	Set the log level from config (if we have it)
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(viper.GetString("loglevel")),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

// GetOutboundIP gets the preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
