package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var (
	WantHelp      bool
	WantVersion   bool
	Scheme        string
	UrlPrefix     string
	HandlingPath  string
	ListeningAddr string
)

const (
	SH            = " (shorthand)"
	VersionSerial = "GoFS Version 2021.11.24"
	UsageTmpl     = `
USAGE:
    {{.AppPath}} [-h] [-v] [--url-prefix <prefix>] [-s {http, https, ftp}] [-a <address>] [-p <path>]

OPTIONS:
    -h, --help
                    show usage
    -v, --version
                    show version
    --url-prefix <prefix>
                    url prefix
    -s {http, https, ftp}, --scheme {http, https, ftp}
                    scheme name (default: "{{.DefaultScheme}}")
    -a <ip:port>, --address <ip:port>
                    listening address (default: "{{.DefaultAddr}}")
    -p </path/to/file>,	--path </path/to/file>
                    handing path or directory (default: "{{.DefaultPath}}")

EXAMPLES:
    {{.AppPath}} -a 10.0.13.120:8080 -p /opt/share0/releases/
    {{.AppPath}} --url-prefix /share/releases/ -a 10.0.13.120:8080 -p /opt/share0/releases/
    {{.AppPath}} --url-prefix /share/releases/ -a=10.0.13.120:8080 -p=/opt/share0/releases/
    {{.AppPath}} --url-prefix=/share/releases/ --address 10.0.13.120:8080 --path /opt/share0/releases/
    {{.AppPath}} --url-prefix=/share/releases/ --address=10.0.13.120:8080 --path=/opt/share0/releases/
`
)

func WantUsage() bool {
	helpArgs := map[string]bool{
		"-h":     true,
		"--h":    true,
		"-help":  true,
		"--help": true,
	}
	for _, arg := range os.Args[1:] {
		if helpArgs[arg] == true {
			return true
		}
	}
	return false
}

func ShowUsage() {
	tmpl := template.Must(template.New("usage tmpl").Parse(UsageTmpl))

	data := struct {
		AppPath       string
		DefaultScheme string
		DefaultAddr   string
		DefaultPath   string
	}{
		AppPath:       os.Args[0],
		DefaultPath:   UserHomeDirMust(),
		DefaultScheme: "http",
		DefaultAddr:   "127.0.0.1:8080",
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		log.Fatalln(err)
	}
}

func ShowVersion() {
	fmt.Println(VersionSerial)
}

func InitArgs() {
	flag.StringVar(&UrlPrefix, "url-prefix", "", "url prefix")

	flag.StringVar(&Scheme, "s", "http", "scheme name"+SH)
	flag.StringVar(&Scheme, "scheme", "http", "scheme name")

	flag.BoolVar(&WantVersion, "v", false, "show version"+SH)
	flag.BoolVar(&WantVersion, "version", false, "show version")

	flag.StringVar(&HandlingPath, "p", UserHomeDirMust(), "handing path or directory"+SH)
	flag.StringVar(&HandlingPath, "path", UserHomeDirMust(), "handing path or directory")

	flag.StringVar(&ListeningAddr, "a", "127.0.0.1:8080", "listening address"+SH)
	flag.StringVar(&ListeningAddr, "address", "127.0.0.1:8080", "listening address")
}

// ParseArgs parses some values from arguments.
func ParseArgs() (addr string, dir string, filename string, prefix string, url string, err error) {
	if WantUsage() {
		ShowUsage()
		err = errors.New("just show usage")
		return
	}

	flag.Parse()
	if WantVersion {
		ShowVersion()
		err = errors.New("just show version")
		return
	}

	info, e := os.Stat(HandlingPath)
	if os.IsNotExist(err) {
		fmt.Printf("No such file or directory: %s\n", HandlingPath)
		err = e
		return
	}

	addr = ListeningAddr // NOTE: set addr
	HandlingPath = AbsPathMust(HandlingPath)
	if info.IsDir() {
		dir = HandlingPath
		filename = ""
	} else {
		dir = filepath.Dir(HandlingPath)       // NOTE: set dir
		filename = filepath.Base(HandlingPath) // NOTE: set filename
	}

	if UrlPrefix != "" {
		if UrlPrefix[0] == '/' {
			UrlPrefix = UrlPrefix[1:]
		}
		if UrlPrefix[len(UrlPrefix)-1] == '/' {
			UrlPrefix = UrlPrefix[:len(UrlPrefix)-1]
		}
		prefix = UrlPrefix
	} else {
		prefix = filepath.Base(dir) // NOTE: set prefix
	}

	url = fmt.Sprintf("%s://%s/%s/%s", "http", addr, prefix, filename) // NOTE: set url

	return
}
