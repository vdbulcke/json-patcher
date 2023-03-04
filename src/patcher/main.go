package patcher

import (
	"io"
	"os"

	"github.com/vdbulcke/json-patcher/src/config"
	"github.com/vdbulcke/json-patcher/src/logger"
	"go.uber.org/zap"
)

// Apply list of patch config
func Apply(config *config.Config, debug bool) error {

	logger := logger.GetZapLogger(debug)

	for _, p := range config.Patches {

		logger.Debug("processing", zap.String("source", p.Source), zap.String("destination", p.Destination))
		err := Patch(p)
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

// readSource read source file or STDIN or return '{}' if NEW as bytes
func readSource(source string) ([]byte, error) {

	if source == config.STDIN {
		return io.ReadAll(os.Stdin)

	}

	if source == config.NEW {
		return []byte("{}"), nil
	}

	return os.ReadFile(source)
}

// writeDestination write bytes to file or STDOUT
func writeDestination(destination string, data []byte) error {
	var writer *os.File
	if destination == config.STDOUT {
		writer = os.Stdout

		_, err := writer.Write(data)
		if err != nil {
			return err
		}

	} else {
		writer, err := os.Create(destination)
		//nolint
		defer writer.Close()
		if err != nil {

			return err
		}

		_, err = writer.Write(data)
		if err != nil {
			return err
		}

	}

	return nil
}
