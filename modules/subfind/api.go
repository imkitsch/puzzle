package subfind

import (
	"bytes"
	"errors"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"os"
	"path/filepath"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
)

var (
	defaultConfigLocation         = filepath.Join(util.GetRunDir() + "/config/subfinder/config.yaml")
	defaultProviderConfigLocation = filepath.Join(util.GetRunDir() + "/config/subfinder/provider-config.yaml")
)

func loadProvidersFrom(location string, options *runner.Options) {
	if len(options.Resolvers) == 0 {
		options.Resolvers = resolve.DefaultResolvers
	}

	// We skip bailing out if file doesn't exist because we'll create it
	// at the end of options parsing from default via goflags.
	if err := runner.UnmarshalFrom(location); isFatalErr(err) && !errors.Is(err, os.ErrNotExist) {
		gologger.Fatalf("Could not read providers from %s: %s", location, err)
	}
}

func DoSubFinder(domain string) []string {
	// Parse the command line flags and read config files
	options := &runner.Options{
		Verbose:            false,
		NoColor:            false,
		Silent:             true,
		RemoveWildcard:     false,
		OnlyRecursive:      false,
		All:                true,
		Threads:            50,
		Timeout:            30,
		MaxEnumerationTime: 10,
		OutputFile:         "",
		OutputDirectory:    "",
		ResolverList:       "",
		Version:            false,
		ExcludeIps:         false,
		CaptureSources:     false,
		HostIP:             false,
		Sources:            []string{},
		Config:             defaultConfigLocation,
		ProviderConfig:     defaultProviderConfigLocation,
		JSON:               false,
		RateLimit:          0,
		Resolvers:          resolve.DefaultResolvers,
		DomainsFile:        "",
	}
	if util.FileExists(options.ProviderConfig) {
		gologger.Infof("Loading provider config file %s", options.ProviderConfig)
		loadProvidersFrom(options.ProviderConfig, options)
	} else {
		gologger.Infof("Loading the default")
		loadProvidersFrom(defaultProviderConfigLocation, options)
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

func isFatalErr(err error) bool {
	return err != nil && !errors.Is(err, io.EOF)
}
