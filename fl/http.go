package main

import (
	"fl_ru/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPathHttp string
)

func init() {
	flag.StringVar(&configPathHttp, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPathHttp, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config, false); err != nil {
		log.Fatal(err)
	}
}
