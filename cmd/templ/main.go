package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/a-h/templ"
	"github.com/a-h/templ/cmd/templ/fmtcmd"
	"github.com/a-h/templ/cmd/templ/generatecmd"
	"github.com/a-h/templ/cmd/templ/lspcmd"
	"github.com/a-h/templ/cmd/templ/migratecmd"
)

func main() {
	run(os.Stdout, os.Args)
}

func run(w io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(w, usageText)
		return 0
	}
	switch args[1] {
	case "generate":
		generateCmd(args[2:])
		return
	case "migrate":
		migrateCmd(args[2:])
		return
	case "fmt":
		fmtCmd(args[2:])
		return
	case "lsp":
		lspCmd(args[2:])
		return
	case "version":
		fmt.Fprintln(w, templ.Version)
		return
	case "--version":
		fmt.Fprintln(w, templ.Version)
		return
	}
	fmt.Fprint(w, usageText)
	return 0
}

const usageText = `usage: templ <command> [parameters]
To see help text, you can run:
  templ generate --help
  templ fmt --help
  templ lsp --help
  templ migrate --help
  templ version
examples:
  templ generate
`

func generateCmd(args []string) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	fileNameFlag := cmd.String("f", "", "Optionally generates code for a single file, e.g. -f header.templ")
	pathFlag := cmd.String("path", ".", "Generates code for all files in path.")
	sourceMapVisualisations := cmd.Bool("sourceMapVisualisations", false, "Set to true to generate HTML files to visualise the templ code and its corresponding Go code.")
	includeVersionFlag := cmd.Bool("include-version", true, "Set to false to skip inclusion of the templ version in the generated code.")
	includeTimestampFlag := cmd.Bool("include-timestamp", false, "Set to true to include the current time in the generated code.")
	watchFlag := cmd.Bool("watch", false, "Set to true to watch the path for changes and regenerate code.")
	cmdFlag := cmd.String("cmd", "", "Set the command to run after generating code.")
	proxyFlag := cmd.String("proxy", "", "Set the URL to proxy after generating code and executing the command.")
	proxyPortFlag := cmd.Int("proxyport", 7331, "The port the proxy will listen on.")
	workerCountFlag := cmd.Int("w", runtime.NumCPU(), "Number of workers to run in parallel.")
	pprofPortFlag := cmd.Int("pprof", 0, "Port to start pprof web server on.")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		cmd.PrintDefaults()
		return
	}
	err = generatecmd.Run(generatecmd.Arguments{
		FileName:                        *fileNameFlag,
		Path:                            *pathFlag,
		Watch:                           *watchFlag,
		Command:                         *cmdFlag,
		Proxy:                           *proxyFlag,
		ProxyPort:                       *proxyPortFlag,
		WorkerCount:                     *workerCountFlag,
		GenerateSourceMapVisualisations: *sourceMapVisualisations,
		IncludeVersion:                  *includeVersionFlag,
		IncludeTimestamp:                *includeTimestampFlag,
		PPROFPort:                       *pprofPortFlag,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func migrateCmd(args []string) {
	cmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	fileName := cmd.String("f", "", "Optionally migrate a single file, e.g. -f header.templ")
	path := cmd.String("path", ".", "Migrates code for all files in path.")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		cmd.PrintDefaults()
		return
	}
	err = migratecmd.Run(migratecmd.Arguments{
		FileName: *fileName,
		Path:     *path,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func fmtCmd(args []string) {
	cmd := flag.NewFlagSet("fmt", flag.ExitOnError)
	helpFlag := cmd.Bool("help", false, "Print help and exit.")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		cmd.PrintDefaults()
		return
	}
	err = fmtcmd.Run(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func lspCmd(args []string) {
	cmd := flag.NewFlagSet("lsp", flag.ExitOnError)
	log := cmd.String("log", "", "The file to log templ LSP output to, or leave empty to disable logging.")
	goplsLog := cmd.String("goplsLog", "", "The file to log gopls output, or leave empty to disable logging.")
	goplsRPCTrace := cmd.Bool("goplsRPCTrace", false, "Set gopls to log input and output messages.")
	helpFlag := cmd.Bool("help", false, "Print help and exit.")
	pprofFlag := cmd.Bool("pprof", false, "Enable pprof web server (default address is localhost:9999)")
	httpDebugFlag := cmd.String("http", "", "Enable http debug server by setting a listen address (e.g. localhost:7474)")
	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		cmd.PrintDefaults()
		return
	}
	err = lspcmd.Run(lspcmd.Arguments{
		Log:           *log,
		GoplsLog:      *goplsLog,
		GoplsRPCTrace: *goplsRPCTrace,
		PPROF:         *pprofFlag,
		HTTPDebug:     *httpDebugFlag,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
