package main

import (
	"fmt"
	"io/ioutil"
	"log"

	M "github.com/XMatrixStudio/Coffee/App/model"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Mongo M.Mongo
}

func main() {
	data, err := ioutil.ReadFile("Config/config.yaml")
	if err != nil {
		log.Printf("ConfigFile.Get err   #%v ", err)
	}
	conf := Config{}
	yaml.Unmarshal(data, &conf)
	fmt.Println(conf)
}
