package patcher

import (
	"bytes"

	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/logger"
	"go.uber.org/zap"
)

type Options struct {
	SkipTags           string
	debug              bool
	AllowUnescapedHTML bool
}

func NewOptions(skipTags string, allowUnescapedHTML, debug bool) *Options {

	return &Options{
		SkipTags:           skipTags,
		AllowUnescapedHTML: allowUnescapedHTML,
		debug:              debug,
	}
}

// Apply list of patch config
func Apply(cfg *config.Config, opts *Options) error {

	logger := logger.GetZapLogger(opts.debug)

	for _, p := range cfg.Patches {

		// handle skips
		r, err := Skip(p, opts)
		if err != nil {
			logger.Debug("error processing", zap.String("source", p.Source), zap.String("destination", p.Destination))
			return err
		}

		if r.Skip {
			logger.Info("skipping patch", zap.String("source", p.Source), zap.String("destination", p.Destination), zap.String("reason", r.Reason))
			continue
		}

		logger.Debug("processing", zap.String("source", p.Source), zap.String("destination", p.Destination))

		// execute patch
		mdoc, err := Patch(p, opts)
		if err != nil {
			return err
		}

		// write patch
		err = WriteDestination(p.Destination, mdoc, opts)
		if err != nil {
			return err
		}

	}

	return nil

}

// Patch a source
func Patch(patch *config.Patch, opts *Options) ([]byte, error) {

	doc, err := readSource(patch.Source, opts)
	if err != nil {
		return nil, err
	}

	mdoc, err := patch.DecodedPatch.Apply(doc)
	if err != nil {
		return nil, err
	}

	// revert HTML escape
	if opts.AllowUnescapedHTML {
		mdoc = bytes.Replace(mdoc, []byte("\\u0026"), []byte("&"), -1)
		mdoc = bytes.Replace(mdoc, []byte("\\u003c"), []byte("<"), -1)
		mdoc = bytes.Replace(mdoc, []byte("\\u003e"), []byte(">"), -1)
	}

	return mdoc, nil
}
