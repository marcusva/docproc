package main

import (
	"flag"
	"fmt"
	"github.com/marcusva/docproc/common/config"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/docproc.webinput/input"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	version = "undefined"

	flagconfig  = flag.String("c", "", "Configuration file to load")
	flaghelp    = flag.Bool("h", false, "Print the help")
	flaglogfile = flag.String("l", "", "Logfile to write information to")
	flaginfo    = flag.Bool("v", false, "Print version details and information about the build configuration.")
)

const (
	cfgFile  = "docproc-webinput.conf"
	usageMsg = `usage: %s [-hv] [-c file] [-l file]

A simple web-based input message converter service.

Options:

  -c <file>   Load the configuration from the passed file.
  -h          Print this help.
  -l <file>   Log information to the passed file.
  -v          Print version details and information about the build configuration.
`
)

func usage() {
	fmt.Fprintf(os.Stderr, usageMsg, os.Args[0])
	os.Exit(2)
}

func info() {
	fmt.Fprintln(os.Stdout, "Version:", version)
	fmt.Fprintln(os.Stdout, "Supported message queue types:")
	fmt.Fprintln(os.Stdout, "  reading:", strings.Join(queue.ReadTypes(), ", "))
	fmt.Fprintln(os.Stdout, "  writing:", strings.Join(queue.WriteTypes(), ", "))
	fmt.Fprintln(os.Stdout, "Supported web handlers:")
	fmt.Fprintln(os.Stdout, "  ", strings.Join(input.WebHandlers(), ", "))
	os.Exit(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *flaghelp {
		flag.Usage()
	}
	if *flaginfo {
		info()
	}

	// Read the configuration
	cfgfile := cfgFile
	if *flagconfig != "" {
		cfgfile = *flagconfig
	}
	conf, err := config.LoadFile(cfgfile, config.NoValidate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Create a logger
	level, err := log.GetLogLevel(conf.GetDefault("log", "level", "Error"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	logfile := conf.GetDefault("log", "file", "")
	if logfile == "" {
		log.Init(os.Stderr, level, true)
	} else {
		if err := log.InitFile(logfile, level, true); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
	}

	// Create the message queue to write to.
	qtype := conf.GetOrPanic("out-queue", "type")
	params := map[string]string{
		"host":  conf.GetOrPanic("out-queue", "host"),
		"topic": conf.GetOrPanic("out-queue", "topic"),
	}
	wq, err := queue.CreateWQ(qtype, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create outbound message queue: %v\n", err)
		os.Exit(2)
	}

	// Create the web service
	address := conf.GetOrPanic("input", "address")
	ws := input.NewWebService(address)
	// Retrieve the input handlers
	sections, err := conf.Array("input", "handlers")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not retrieve input handlers: %v\n", err)
		os.Exit(2)
	}
	for _, sec := range sections {
		params, err := conf.AllFor(sec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve section '%s': %v\n", sec, err)
			os.Exit(2)
		}
		endpoint, ok := params["endpoint"]
		if !ok {
			fmt.Fprintf(os.Stderr, "web handler '%s' misses 'endpoint'", sec)
			os.Exit(2)
		}
		handler, err := input.Create(wq, params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "creating the web handler failed: %v\n", err)
			os.Exit(2)
		}
		ws.Bind(endpoint, handler)
	}
	if err := wq.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "opening the message queue failed: %v\n", err)
		os.Exit(2)
	}
	ws.Start()

	// Catch some signals to allow a graceful daemon shutdown
	log.Infof("Connecting signal handler")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigchan
	log.Infof("Shutdown on signal %s", s.String())

	exitval := 0
	if err := ws.Stop(); err != nil {
		log.Errorf("Could not stop the webservice properly: %v\n", err)
		exitval = 2
	}
	if err := wq.Close(); err != nil {
		log.Errorf("Could not close outbound message queue(s): %v\n", err)
		exitval = 2
	}
	os.Exit(exitval)
}
