package subfind

import (
	"bytes"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"path/filepath"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
)

var (
	defaultConfigLocation         = filepath.Join(util.GetRunDir() + "/config/subfinder/config.yaml")
	defaultProviderConfigLocation = filepath.Join(util.GetRunDir() + "/config/subfinder/provider-config.yaml")
)

func DoSubFinder(domain string) []string {
	// Parse the command line flags and read config files
	options := &runner.Options{
		Verbose:            false,
		NoColor:            false,
		Silent:             false,
		RemoveWildcard:     false,
		OnlyRecursive:      false,
		All:                true,
		Threads:            50,
		Timeout:            30,
		MaxEnumerationTime: 10,
		Version:            false,
		CaptureSources:     false,
		HostIP:             false,
		Config:             defaultConfigLocation,
		ProviderConfig:     defaultProviderConfigLocation,
		JSON:               false,
		RateLimit:          0,
		Resolvers:          resolve.DefaultResolvers,
	}

	newRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatalf("newRunner, err := runner.NewRunner Could not create runner: %s\n", err)
	}
	buf := bytes.Buffer{}
	err = newRunner.EnumerateSingleDomain(domain, []io.Writer{&buf})
	if err != nil {
		gologger.Fatalf("Could not run enumeration: %s\n", err)
	}
	data, err := io.ReadAll(&buf)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	domains := strings.Split(string(data), "\n")
	return domains
}
