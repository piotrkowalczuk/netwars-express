package main

import (
	"encoding/xml"
	"github.com/piotrkowalczuk/netwars-backend/database"
	"github.com/piotrkowalczuk/netwars-backend/service"
	"log"
	"os"
)

type Config struct {
	Server  Server                 `xml:"server"`
	Redis   database.RedisConfig   `xml:"redis"`
	Postgre database.PostgreConfig `xml:"postgre"`
	Sentry  service.SentryConfig   `xml:"sentry"`
}

type Server struct {
	Host         string `xml:"host"`
	Port         string `xml:"port"`
	ReadTimeout  int    `xml:"read_timeout"`
	WriteTimeout int    `xml:"write_timeout"`
}

func openFile(filePath string) (file *os.File) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Cannot open configuration file: %v\n", err)
		os.Exit(1)
	}

	return
}

func ReadConfiguration(filePath string) (config *Config) {
	config = new(Config)
	file := openFile(filePath)
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.Decode(&config)

	return
}
