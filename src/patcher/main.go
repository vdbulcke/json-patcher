package patcher

import (
	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/logger"
	"go.uber.org/zap"
)

// Apply list of patch config
func Apply(cfg *config.Config, skipTags string, debug bool) error {

	logger := logger.GetZapLogger(debug)

	for _, p := range cfg.Patches {

		// handle skips
		r, err := skip(p, skipTags)
		if err != nil {
			logger.Debug("error processing", zap.String("source", p.Source), zap.String("destination", p.Destination))
			return err
		}

		if r.Skip {
			logger.Info("skipping patch", zap.String("source", p.Source), zap.String("destination", p.Destination), zap.String("reason", r.Reason))
			continue
		}

		logger.Debug("processing", zap.String("source", p.Source), zap.String("destination", p.Destination))

		err = Patch(p)
		if err != nil {
			return err
		}

	}

	return nil

}

// Patch a source into a destination
func Patch(patch *config.Patch) error {

	doc, err := readSource(patch.Source)
	if err != nil {
		return err
	}

	mdoc, err := patch.DecodedPatch.Apply(doc)
	if err != nil {
		return err
	}

	return writeDestination(patch.Destination, mdoc)

}
