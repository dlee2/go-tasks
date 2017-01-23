package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
)

type Configuration struct {
    ServerPort string
}

var err error
var config Configuration

func ReadConfig(fileName string) (Configuration, error) {
    configFile, err := ioutil.ReadFile(fileName)
    if err != nil {
        log.Print("Unable to read config file, switching to flag mode")
        return Configuration{}, err
    }

    //log.Print(configFile)
    err = json.Unmarshal(configFile, &config)
    if err != nil {
        log.Print("Invalid JSON, expecting port from command line flag")
        return Configuration{}, err
    }
    return config, nil
}