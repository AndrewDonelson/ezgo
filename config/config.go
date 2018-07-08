package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/camelcase"
	"github.com/hashicorp/hcl"
	yaml "gopkg.in/yaml.v2"
)

//Config is the primary configuration structure for your application. Modifiy the properties as you see fit.
type Config struct {
	aString  string
	aInteger int
	aFloat   float64
	aBoolean bool
	aStruct  credentials
}

type credentials struct {
	username string
	password string
}

//AddFlags makes each configuration property available on commandline
func (cfg *Config) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.aString, "string", cfg.aString, "Framework user to register with the Mesos master")
	fs.IntVar(&cfg.aInteger, "integer", cfg.aInteger, "Framework name to register with the Mesos master")
	fs.Float64Var(&cfg.aFloat, "float", cfg.aFloat, "Framework role to register with the Mesos master")
	fs.BoolVar(&cfg.aBoolean, "boolean", cfg.aBoolean, "Codec to encode/decode scheduler API communications [protobuf, json]")
	fs.StringVar(&cfg.aStruct.username, "aStruct.username", cfg.aStruct.username, "Username for Mesos authentication")
	fs.StringVar(&cfg.aStruct.password, "aStruct.passwordFile", cfg.aStruct.password, "Path to file that contains the password for Mesos authentication")
}

//DefaultConfig crates adn returns a new configuration object
func DefaultConfig() Config {
	return Config{
		aString:  env("A_STRING_VALUE", "testString"),
		aInteger: envInt("AN_INTEGER_VALUE", "5"),
		aFloat:   envFloat("A_FLOAT_VALUE", "3.12569"),
		aBoolean: true,
		aStruct: credentials{
			username: env("AUTH_USER", ""),
			password: env("AUTH_PASSWORD_FILE", ""),
		},
	}
}

// LoadConfig reads configuration from path. The format is deduced from the file extension
//	* .json    - is decoded as json
//	* .yml     - is decoded as yaml
//	* .toml    - is decoded as toml
//  * .hcl	   - is decoded as hcl
func LoadConfig(name, path string) (*Config, error) {
	//If no filename given use default 'config'
	if name == "" {
		name = "config"
	}

	//if no path is given use default current folder (exe)
	if path == "" {
		path = "./"
	}
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	switch filepath.Ext(path) {
	case ".json":
		jerr := json.Unmarshal(data, cfg)
		if jerr != nil {
			return nil, jerr
		}
	case ".toml":
		_, terr := toml.Decode(string(data), cfg)
		if terr != nil {
			return nil, terr
		}
	case ".yml":
		yerr := yaml.Unmarshal(data, cfg)
		if yerr != nil {
			return nil, yerr
		}
	case ".hcl":
		obj, herr := hcl.Parse(string(data))
		if herr != nil {
			return nil, herr
		}
		if herr = hcl.DecodeObject(&cfg, obj); herr != nil {
			return nil, herr
		}
	default:
		return nil, fmt.Errorf("EZGO-Config: Config file format [%s] not supported", filepath.Ext(path))
	}

	err = cfg.syncEnv()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SyncEnv overrides c field's values that are set in the environment.
//
// The environment variable names are derived from config fields by underscoring, and uppercasing
// the name. E.g. AppName will have a corresponding environment variable APP_NAME
//
// NOTE only int, string and bool fields are supported and the corresponding values are set.
// when the field value is not supported it is ignored.
func (cfg *Config) syncEnv() error {
	c := reflect.ValueOf(cfg).Elem()
	cTyp := c.Type()

	for k := range make([]struct{}, cTyp.NumField()) {
		field := cTyp.Field(k)

		cm := getEnvName(field.Name)
		env := os.Getenv(cm)
		if env == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			c.FieldByName(field.Name).SetString(env)
		case reflect.Int:
			v, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("EZGO-Config: Loading config field %s %v", field.Name, err)
			}
			c.FieldByName(field.Name).Set(reflect.ValueOf(v))
		case reflect.Bool:
			b, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("EZGO-Config: Loading config field %s %v", field.Name, err)
			}
			c.FieldByName(field.Name).SetBool(b)
		}

	}
	return nil
}

// getEnvName returns all upper case and underscore separated string, from field.
// field is a camel case string.
//
// example
//	AppName will change to APP_NAME
func getEnvName(field string) string {
	camSplit := camelcase.Split(field)
	var rst string
	for k, v := range camSplit {
		if k == 0 {
			rst = strings.ToUpper(v)
			continue
		}
		rst = rst + "_" + strings.ToUpper(v)
	}
	return rst
}
