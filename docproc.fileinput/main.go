package main

import (
	"flag"
	"fmt"
	"github.com/marcusva/docproc/common/config"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/service"
	"github.com/marcusva/docproc/docproc.fileinput/input"
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
	cfgFile  = "docproc-fileinput.conf"
	usageMsg = `usage: %s [-hv] [-c file] [-l file]

A simple file-based input message converter service

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
	fmt.Fprintln(os.Stdout, "File input configuration:")
	fmt.Fprintln(os.Stdout, "  transformers:", strings.Join(input.FileTransfomers(), ", "))
	fmt.Fprintln(os.Stdout, "  default check interval:", input.CheckInterval)
	os.Exit(0)
}

// basic configuration validation
func validateConfig(cfg *config.Config) error {
	// Check, that handlers do exist
	handlers, err := cfg.Array("input", "handlers")
	if err != nil {
		return err
	}
	if len(handlers) == 0 {
		return fmt.Errorf("no input handlers configured")
	}
	// Check for conflicting file patterns
	patterns := make(map[string][]string)
	for _, h := range handlers {
		dir, err := cfg.Get(h, "folder.in")
		if err != nil {
			return err
		}
		pattern, err := cfg.Get(h, "pattern")
		if err != nil {
			return err
		}
		for k, v := range patterns {
			alike := strings.Contains(pattern, v[1]) || strings.Contains(v[1], pattern)
			if v[0] == dir && alike {
				return fmt.Errorf("file patterns for '%s' and '%s' overlap", h, k)
			}
		}
		patterns[h] = []string{dir, pattern}
	}
	return nil
}

func start(conf *config.Config) {
	// Create a logger
	level, err := log.GetLogLevel(conf.GetDefault("log", "level", "Error"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	if logfile := conf.GetDefault("log", "file", ""); logfile == "" {
		log.Init(os.Stderr, level, true)
	} else {
		if err := log.InitFile(logfile, level, true); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
	}

	// Get the queue to write to
	qtype := conf.GetOrPanic("out-queue", "type")
	params := map[string]string{
		"host":  conf.GetOrPanic("out-queue", "host"),
		"topic": conf.GetOrPanic("out-queue", "topic"),
	}
	wqueue, err := queue.CreateWQ(qtype, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create message queue: %v\n", err)
		os.Exit(2)
	}

	// Retrieve the input handlers
	sections, err := conf.Array("input", "handlers")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not retrieve input handlers: %v\n", err)
		os.Exit(2)
	}
	watchers := []*service.FileWatcher{}
	for _, sec := range sections {
		params, err := conf.AllFor(sec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve section '%s': %v\n", sec, err)
			os.Exit(2)
		}
		fw, err := input.Create(wqueue, params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "creating the filewatcher failed: %v\n", err)
			os.Exit(2)
		}
		watchers = append(watchers, fw)
	}
	// Open the queue
	if err := wqueue.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "could not open message queue: %v\n", err)
		os.Exit(2)
	}

	// Start all watchers
	for _, fw := range watchers {
		go fw.Watch()
	}

	// Catch some signals to allow a graceful daemon shutdown
	log.Infof("Connecting signal handler")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigchan
	log.Infof("Shutdown on signal %s", s.String())

	/* FIXME: graceful shutdown of everything */
	for _, fw := range watchers {
		fw.Stop()
	}
	wqueue.Close()
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
	conf, err := config.LoadFile(cfgfile, validateConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	start(conf)
}
