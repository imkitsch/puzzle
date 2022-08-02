package subfind

import (
	"Allin/gologger"
	"Allin/util"
	"bytes"
	"context"
	"errors"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	defaultConfigLocation         = filepath.Join(util.GetRunDir() + "/config/subfinder/config.yaml")
	defaultProviderConfigLocation = filepath.Join(util.GetRunDir() + "/config/subfinder/provider-config.yaml")
)

func loadProvidersFrom(location string, options *runner.Options) {
	if len(options.AllSources) == 0 {
		options.AllSources = passive.DefaultAllSources
	}
	if len(options.Recursive) == 0 {
		options.Recursive = passive.DefaultRecursiveSources
	}
	// todo: move elsewhere
	if len(options.Resolvers) == 0 {
		options.Recursive = resolve.DefaultResolvers
	}
	if len(options.Sources) == 0 {
		options.Sources = passive.DefaultSources
	}

	options.Providers = &runner.Providers{}
	// We skip bailing out if file doesn't exist because we'll create it
	// at the end of options parsing from default via flags.
	if err := options.Providers.UnmarshalFrom(location); isFatalErr(err) && !errors.Is(err, os.ErrNotExist) {
		gologger.Fatalf("Could not read providers from %s: %s\n", location, err)
	}
}

func DoSubFinder(domain string) []string {
	// Parse the command line flags and read config files
	options := &runner.Options{
		Verbose:            false,
		NoColor:            false,
		Silent:             true,
		RemoveWildcard:     true,
		All:                false,
		OnlyRecursive:      false,
		Threads:            50,
		Timeout:            30,
		MaxEnumerationTime: 10,
		OutputFile:         "",
		OutputDirectory:    "",
		ResolverList:       "",
		Proxy:              "",
		Version:            false,
		ExcludeIps:         false,
		CaptureSources:     false,
		HostIP:             false,
		Config:             defaultConfigLocation,
		ProviderConfig:     defaultProviderConfigLocation,
		JSON:               false,
		RateLimit:          0,
		ExcludeSources:     []string{},
		Resolvers:          []string{},
		Sources:            []string{},
		DomainsFile:        ""}
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
	err = newRunner.EnumerateSingleDomain(context.Background(), domain, []io.Writer{&buf})
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
