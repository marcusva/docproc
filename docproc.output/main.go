package main

import (
	"flag"
	"fmt"
	"github.com/marcusva/docproc/common/config"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/docproc.output/output"
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
	cfgFile  = "docproc-output.conf"
	usageMsg = `usage: %s [-hv] [-c file] [-l file]

A simple document output service

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
	fmt.Fprintln(os.Stdout, "Supported output handlers:")
	fmt.Fprintln(os.Stdout, "  ", strings.Join(output.Writers(), ", "))
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

	// Create the message queue to read from
	inqtype := conf.GetOrPanic("in-queue", "type")
	inparams := map[string]string{
		"host":  conf.GetOrPanic("in-queue", "host"),
		"topic": conf.GetOrPanic("in-queue", "topic"),
	}
	rq, err := queue.CreateRQ(inqtype, inparams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create inbound message queue: %v\n", err)
		os.Exit(2)
	}

	handler := output.NewHandler()
	sections, err := conf.Array("output", "handlers")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not retrieve output handlers: %v\n", err)
		os.Exit(2)
	}

	for _, sec := range sections {
		params, err := conf.AllFor(sec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve section '%s': %v\n", sec, err)
			os.Exit(2)
		}
		proc, err := output.Create(params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create renderer: %v\n", err)
			os.Exit(2)
		}
		handler.AddProcessor(proc)
	}
	if err := rq.Open(handler); err != nil {
		fmt.Fprintf(os.Stdout, "could not open inbound message queue: %v\n", err)
		os.Exit(2)
	}

	// Create an error queue, if provided
	if conf.HasSection("error-queue") {
		errqtype := conf.GetOrPanic("error-queue", "type")
		errparams := map[string]string{
			"host":  conf.GetOrPanic("error-queue", "host"),
			"topic": conf.GetOrPanic("error-queue", "topic"),
		}
		_, err := queue.CreateWQ(errqtype, errparams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create error message queue: %v\n", err)
			os.Exit(2)
		}
		// TODO
	}

	// Catch some signals to allow a graceful daemon shutdown
	log.Infof("Connecting signal handler")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigchan
	log.Infof("Shutdown on signal %s", s.String())

	/* FIXME: graceful shutdown of everything */
	if err := rq.Close(); err != nil {
		log.Errorf("Could not close message queue: %v\n", err)
		os.Exit(2)
	}

	os.Exit(0)
}
