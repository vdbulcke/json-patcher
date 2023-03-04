package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
	k8syaml "sigs.k8s.io/yaml"
)

var (
	STDIN  = "STDIN"
	STDOUT = "STDOUT"
	NEW    = "NEW"
)

type Config struct {
	Patches []*Patch `yaml:"patches" validate:"required"`
}

type Patch struct {
	Source      string `yaml:"source" validate:"required"`
	Destination string `yaml:"destination" validate:"required"`
	JSONPatch   string `yaml:"json_patch" validate:"required"`

	DecodedPatch jsonpatch.Patch
}

// ValidateConfig validate config
func ValidateConfig(config *Config) bool {

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	errs := validate.Struct(config)

	if errs == nil {
		return true
	}

	for _, e := range errs.(validator.ValidationErrors) {
		fmt.Println(e)
	}

	return false

}

// ParseConfig Parse config file
func ParseConfig(configFile string) (*Config, error) {

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := Config{}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	// validate config
	if !ValidateConfig(&config) {
		return nil, fmt.Errorf("validation Error %s", configFile)
	}

	// process jsonpatch
	for _, p := range config.Patches {

		trimedPatch := strings.TrimSpace(p.JSONPatch)
		p.DecodedPatch, err = jsonPatchFromString(trimedPatch)
		if err != nil {
			return nil, err
		}

	}

	// return Parse config struct
	return &config, nil

}

// jsonPatchFromString loads a Json 6902 patch
// source: https://github.com/kubernetes-sigs/kustomize/blob/master/plugin/builtin/patchtransformer/PatchTransformer.go
func jsonPatchFromString(ops string) (jsonpatch.Patch, error) {

	if ops == "" {
		return nil, fmt.Errorf("empty json patch operations")
	}

	if ops[0] != '[' {
		jsonOps, err := k8syaml.YAMLToJSON([]byte(ops))
		if err != nil {
			return nil, err
		}
		ops = string(jsonOps)
	}
	return jsonpatch.DecodePatch([]byte(ops))
}
