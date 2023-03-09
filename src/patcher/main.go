package patcher

import (
	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/logger"
	"go.uber.org/zap"
)

type Options struct {
	skipTags           string
	debug              bool
	allowUnescapedHTML bool
}

func NewOptions(skipTags string, allowUnescapedHTML, debug bool) *Options {

	return &Options{
		skipTags:           skipTags,
		allowUnescapedHTML: allowUnescapedHTML,
		debug:              debug,
	}
}

// Apply list of patch config
func Apply(cfg *config.Config, opts *Options) error {

	logger := logger.GetZapLogger(opts.debug)

	for _, p := range cfg.Patches {

		// handle skips
		r, err := skip(p, opts.skipTags)
		if err != nil {
			logger.Debug("error processing", zap.String("source", p.Source), zap.String("destination", p.Destination))
			return err
		}

		if r.Skip {
			logger.Info("skipping patch", zap.String("source", p.Source), zap.String("destination", p.Destination), zap.String("reason", r.Reason))
			continue
		}

		logger.Debug("processing", zap.String("source", p.Source), zap.String("destination", p.Destination))

		err = Patch(p, opts)
		if err != nil {
			return err
		}

	}

	return nil

}

// Patch a source into a destination
func Patch(patch *config.Patch, opts *Options) error {

	doc, err := readSource(patch.Source, opts)
	if err != nil {
		return err
	}

	mdoc, err := patch.DecodedPatch.Apply(doc)
	if err != nil {
		return err
	}

	return writeDestination(patch.Destination, mdoc, opts)

}
