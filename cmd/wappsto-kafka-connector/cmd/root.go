package cmd

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"wappsto-kafka-connector/internal/config"
	"wappsto-kafka-connector/internal/wappsto"
)

var version string
var cfgFile string
var interrupt chan os.Signal

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Wappsto Connector version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var rootCmd = &cobra.Command{
	Use:   "wappsto-kafka-connector",
	Short: "Wappsto connection to Kafka service",
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		tasks := []func() error{
			printStartMessage,
			setupWappsto,
			handleWappstoStream,
		}

		for _, t := range tasks {
			if err := t(); err != nil {
				log.Fatal(err)
			}
		}

		interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully
		exitChan := make(chan struct{})
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // Notify the interrupt channel for SIGINT

		log.WithField("signal", <-interrupt).Info("signal received")
		go func() {
			log.Warning("stopping wappsto-kafka-connector")
			exitChan <- struct{}{}
		}()

		select {
		case <-exitChan:
		case s := <-interrupt:
			log.WithField("signal", s).Info("signal received, stopping immediately")
		}

		return nil
	},
}

func init() {

	cobra.OnInitialize(initConfig)
	viper.SetDefault("wappsto.server", "wappsto.com")

	rootCmd.AddCommand(versionCmd)
}
func printStartMessage() error {
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "https://documentation.wappsto.com",
	}).Info("starting Wappsto connector")

	return nil
}

func setupWappsto() error {
	if err := wappsto.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup wappsto stream error")
	}
	return nil
}

func handleWappstoStream() error {
	// TODO:
	go wappsto.HandleWappstoStream()
	return nil
}

func initConfig() {
	//
	if cfgFile != "" {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
		viper.SetConfigType("toml")
		if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
	} else {
		log.Println("Reading configuration file")
		confName := "wappsto-kafka-connector"
		viper.SetConfigName(confName)
		viper.AddConfigPath("$HOME/.config/" + confName)
		viper.AddConfigPath("/etc/" + confName)
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				log.WithError(err).WithField("config", cfgFile).Fatal("configuration file: " + confName + ".toml not found")
			default:
				log.WithError(err).Fatal("read configuration file " + confName + ".toml error")
			}
		}

		log.Printf("All keys: %s", viper.AllKeys())
	}

	viperBindEnvs(config.C)

	if err := viper.Unmarshal(&config.C); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}

}

/**
If environmental variable is defined for configuration, such environental variable will be preoritized
instead of key value pair in configuration file.
*/
func viperBindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			tv = strings.ToLower(t.Name)
		}
		if tv == "-" {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			viperBindEnvs(v.Interface(), append(parts, tv)...)
		default:
			// Bash doesn't allow env variable names with a dot so
			// bind the double underscore version.
			keyDot := strings.Join(append(parts, tv), ".")
			keyUnderscore := strings.Join(append(parts, tv), "__")
			viper.BindEnv(keyDot, strings.ToUpper(keyUnderscore))
		}
	}
}

func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
