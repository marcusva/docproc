package main

import (
	"flag"
	"fmt"
	"github.com/marcusva/docproc/common/config"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/queue/processors"
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
	cfgFile  = "docproc.conf"
	usageMsg = `usage: %s [-hv] [-c file] [-l file]

A simple content processing command.

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
	fmt.Fprintln(os.Stdout, "Supported message handlers:")
	fmt.Fprintln(os.Stdout, "  ", strings.Join(processors.Types(), ", "))
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

	// Create the message queue to read from.
	inqtype := conf.GetOrPanic("in-queue", "type")
	inparams := map[string]string{
		"host":  conf.GetOrPanic("in-queue", "host"),
		"topic": conf.GetOrPanic("in-queue", "topic"),
	}
	inqueue, err := queue.CreateRQ(inqtype, inparams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create inbound message queue: %v\n", err)
		os.Exit(2)
	}
	// Create the message queue to write to, if any.
	outqueue := getOutQueue("out-queue", conf)
	// Create the message queue to use for errors, if any.
	errqueue := getOutQueue("error-queue", conf)

	var consumer queue.ProcConsumer
	if outqueue != nil || errqueue != nil {
		writer := queue.NewWriter(outqueue, errqueue)
		if err := writer.Open(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not create outbound message queue(s): %v\n", err)
			os.Exit(2)
		}
		consumer = writer
	} else {
		consumer = queue.NewSimpleConsumer()
	}

	// Setup the processors to be executed on new messages
	sections, err := conf.Array("execute", "handlers")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not retrieve message handlers: %v\n", err)
		os.Exit(2)
	}
	for _, sec := range sections {
		params, err := conf.AllFor(sec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve section '%s': %v\n", sec, err)
			os.Exit(2)
		}
		proc, err := processors.Create(params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create processor: %v\n", err)
			os.Exit(2)
		}
		consumer.Add(proc)
	}
	if err := inqueue.Open(consumer); err != nil {
		fmt.Fprintf(os.Stdout, "could not open inbound message queue: %v\n", err)
		os.Exit(2)
	}

	// Catch some signals to allow a graceful daemon shutdown
	log.Infof("Connecting signal handler")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigchan
	log.Infof("Shutdown on signal %s", s.String())

	exitval := 0
	if writer, ok := consumer.(*queue.Writer); ok {
		if err := writer.Close(); err != nil {
			log.Errorf("Could not close outbound message queue(s): %v\n", err)
			exitval = 2
		}
	}
	if err := inqueue.Close(); err != nil {
		log.Errorf("Could not close inbound message queue: %v\n", err)
		exitval = 2
	}

	os.Exit(exitval)
}

func getOutQueue(name string, conf *config.Config) queue.WriteQueue {
	if !conf.HasSection(name) {
		return nil
	}
	qtype := conf.GetOrPanic(name, "type")
	params := map[string]string{
		"host":  conf.GetOrPanic(name, "host"),
		"topic": conf.GetOrPanic(name, "topic"),
	}
	wq, err := queue.CreateWQ(qtype, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create outbound message queue: %v\n", err)
		os.Exit(2)
	}
	return wq
}
