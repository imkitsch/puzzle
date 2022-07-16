package subfind

import (
	"Allin/gologger"
	"errors"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"os"
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
	// at the end of options parsing from default via goflags.
	if err := options.Providers.UnmarshalFrom(location); isFatalErr(err) && !errors.Is(err, os.ErrNotExist) {
		gologger.Fatalf("Could not read providers from %s: %s\n", location, err)
	}
}

func isFatalErr(err error) bool {
	return err != nil && !errors.Is(err, io.EOF)
}
