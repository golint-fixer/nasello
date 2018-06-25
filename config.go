package nasello

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Resolvers map[string]resolver
	Rules     map[string]rule
}

type resolver struct {
	Servers  []string
	Port     int
	Protocol string
}

type rule struct {
	Description string
	Match       string
	Resolver    string
}

// ReadConfig reads a JSON file and returns a Configuration object
// containing the raw elements.
func ReadConfig(filename string) Configuration {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't open config file: %s\n", err.Error())
	}

	var jsonConfig Configuration
	toml.Unmarshal(file, &jsonConfig)

	// Safety checks
	if len(jsonConfig.Resolvers) == 0 || len(jsonConfig.Rules) == 0 {
		log.Fatal("Configuration contains no 'filters' section")
	}

	return jsonConfig
}
