package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const bannerMsg = `
____ ____ _    _    ___  __   ___    __   ____ __   ____ 
| __\|   ||\/\ |\/\ |  \ | \|\|  \   | |  |___\| \|\| __\
| \__| . ||   \|   \| . \|  \|| . \  | |__| /  |  \||  ]_
|___/|___/|/v\/|/v\/|/\_/|/\_/|___/  |___/|/   |/\_/|___/
`

var configFile string

func init() {
	// set config defaults
	viper.SetDefault("garbage-collect", false)
	viper.SetConfigType("yml")
	viper.SetEnvPrefix("commandline")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// flags
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode. Do not display banner messages")
	RootCmd.PersistentFlags().DurationP("shutdown-wait-time", "", time.Duration(500)*time.Millisecond, "Shutdown wait time")

	// config
	viper.BindPFlag("quiet", RootCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("shutdown-wait-time", RootCmd.PersistentFlags().Lookup("shutdown-wait-time"))

	// local flags;
	RootCmd.Flags().StringVar(&configFile, "config", "", "/path/to/config.yml")

}

// RootCmd is the main command to run the application
var RootCmd = &cobra.Command{
	Use:   "commandline",
	Short: "Command Line Boilerplate",
	Long:  bannerMsg,
	Run:   run,

	// parse the config if one is provided, or use the defaults. Set the backend
	// driver to be used
	PersistentPreRun: preRun,
}

func run(cmd *cobra.Command, args []string) {
	if !viper.GetBool("quiet") {
		fmt.Println(bannerMsg)
		fmt.Println()
		fmt.Printf("Started with process id: %d\n", os.Getpid())
	}

	// track signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	var waitTime time.Duration
	for {
		waitTime = viper.GetDuration("shutdown-wait-time")
		if exit := start(sig, waitTime); exit {
			return
		}
	}

}

func preRun(ccmd *cobra.Command, args []string) {
	// if --config is passed, attempt to parse the config file
	if configFile != "" {
		// get the filepath
		abs, err := filepath.Abs(configFile)
		if err != nil {
			log.Fatalf("Error reading filepath: %s", err)
			os.Exit(1)
		}

		// get the config name
		base := filepath.Base(abs)

		// get the path
		path := filepath.Dir(abs)

		//
		viper.SetConfigName(strings.Split(base, ".")[0])
		viper.AddConfigPath(path)

		// Find and read the config file; Handle errors reading the config file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Failed to read config file: %s", err)
			os.Exit(1)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s\n", e.Name)
			// Send a HUP signal to restart
			if p, err := os.FindProcess(os.Getpid()); err == nil {
				p.Signal(os.Interrupt)
			}
		})
	}
}

func start(sig <-chan os.Signal, waittime time.Duration) bool {
	log.Println("Start!")

	var exitApplication bool

	shutdown := make(chan interface{})
	exit := make(chan interface{})

	go myApplicationMain(shutdown, exit)

	switch <-sig {
	case syscall.SIGINT, syscall.SIGTERM:
		exitApplication = true

	// case syscall.SIGHUP:
	default:
		log.Println("Reload!!")
	}

	close(shutdown)
	time.Sleep(waittime)
	close(exit)

	log.Println("Exit!")

	return exitApplication
}

func myApplicationMain(shutdown, exit <-chan interface{}) {
	log.Println("myApplicationMain")
	for {
		select {
		case <-shutdown:
			log.Println("myApplicationMain: shutdown started...")
			// some long process....
			time.Sleep(time.Second * 5)
			log.Println("myApplicationMain: shutdown completed.")
			return

		case <-exit:
			log.Println("myApplicationMain: exit now.")
			return

		}
	}
}
