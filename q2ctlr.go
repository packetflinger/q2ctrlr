package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ThinkFunc func()

type Servers struct {
	Sv []Server `json:"servers"`
}
type Server struct {
	Name   string   `json:"name"`
	IP     string   `json:"ip"`
	Port   int      `json:"port"`
	Groups []string `json:"groups"`
}

type MVDFilelist struct {
	Filename []string
}

var svs Servers

func main() {

}

func init() {
	// load config from ~/.config/q2ctrlr/config.json
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}

	cfgfile := fmt.Sprintf(
		"%s%c.config%cq2ctrlr%cconfig.json",
		dirname,
		os.PathSeparator,
		os.PathSeparator,
		os.PathSeparator)

	log.Printf("Loading config file: %s\n", cfgfile)
	configdata, err := os.ReadFile(cfgfile)
	if err != nil {
		log.Println(err)
		return
	}

	json.Unmarshal(configdata, &svs)

	for _, srv := range svs.Sv {
		log.Println(srv)
	}
}
