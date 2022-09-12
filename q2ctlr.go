package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type ThinkFunc func()

type Servers struct {
	Sv []Server `json:"servers"`
}
type Server struct {
	Name        string   `json:"name"`
	IP          string   `json:"ip"`
	Port        int      `json:"port"`
	Groups      []string `json:"groups"`
	ThinkOffset int64    `json:"thinkoffset"`
	NextThink   int64
}

type MVDFilelist struct {
	NumFiles int
	Filename []string
}

var svs Servers

func (sv *Server) Think() {
	// not time yet
	if sv.NextThink > time.Now().Unix() {
		return
	}

	log.Printf("%s thinking\n", sv.Name)

	sv.NextThink = time.Now().Unix() + sv.ThinkOffset
	url := fmt.Sprintf("http://%s:%d/GetMVDFiles", sv.IP, sv.Port)

	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatal(response.StatusCode)
		return
	}

	bodybytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	files := MVDFilelist{}
	json.Unmarshal(bodybytes, &files)
}

func main() {

	for {
		for i := range svs.Sv {
			svs.Sv[i].Think()
		}
		time.Sleep(10 * time.Second)
	}
	/*
		url := "http://de:27999/GetMVDFiles"

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			log.Fatal(response.StatusCode)
			return
		}

		bodybytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(string(bodybytes))
		files := MVDFilelist{}
		json.Unmarshal(bodybytes, &files)

		for _, f := range files.Filename {
			log.Println(f)
		}
	*/
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

	for i := range svs.Sv {
		svs.Sv[i].NextThink = time.Now().Unix() + 120 + int64(rand.Intn(120))
		log.Println(svs.Sv[i])
	}
}
