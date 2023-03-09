package patcher

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vdbulcke/json-patcher/src/config"
)

type SkipReason struct {
	Skip   bool
	Reason string
}

// sourceExists return if source file exists
func sourceExists(source string) bool {

	_, err := os.Stat(source)
	return !errors.Is(err, os.ErrNotExist)

}

// readSource read source file or STDIN or return '{}' if NEW as bytes
func readSource(source string, opts *Options) ([]byte, error) {

	if source == config.STDIN {
		return io.ReadAll(os.Stdin)

	}

	if source == config.NEW {
		return []byte("{}"), nil
	}

	return os.ReadFile(source)
}

// writeDestination write bytes to file or STDOUT
func writeDestination(destination string, data []byte, opts *Options) error {

	// revert HTML escape
	if opts.allowUnescapedHTML {
		data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
		data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
		data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	}

	var writer *os.File
	if destination == config.STDOUT {
		writer = os.Stdout

		_, err := writer.Write(data)
		if err != nil {
			return err
		}

		// for STDOUT add line break after each data
		_, err = writer.WriteString("\n")
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

func skip(patch *config.Patch, skipTags string) (*SkipReason, error) {

	// handle skip tags
	for _, t := range strings.Split(skipTags, ",") {

		// for each skipTags if at least one
		// is matching a patch tag then skip
		if stringInSlice(t, patch.Tags) {
			return &SkipReason{
				Skip:   true,
				Reason: fmt.Sprintf("skip tag '%s'", t),
			}, nil
		}
	}

	// handle source
	if patch.Source != config.NEW && patch.Source != config.STDIN &&
		!sourceExists(patch.Source) {

		switch patch.SourceNotExists {
		case config.CONTINUE:

			return &SkipReason{
				Skip:   true,
				Reason: "continue if source does not exists",
			}, nil

		default:
			return &SkipReason{
				Skip:   false,
				Reason: "fail if source does not exists",
			}, fmt.Errorf("source %s does not exists", patch.Source)

		}

	}

	return &SkipReason{
		Skip:   false,
		Reason: "",
	}, nil

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
