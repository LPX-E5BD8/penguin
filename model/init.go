package model

import (
	"log"
	"os"
	"path"

	"github.com/panjf2000/ants"
)

// base dir
var WorkDir = path.Join(os.Getenv("HOME"), ".penguin")

// LogDir to save log
var LogDir = path.Join(WorkDir, "log")

// CacheDir to save cache files
var CacheDir = path.Join(WorkDir, "cache")

// Logger for app
var Logger *log.Logger
var LogFile *os.File

var Pool, _ = ants.NewPool(100)

func init() {
	dirValidation()
	loggerInit()
}

func dirValidation() {
	dirs := []string{WorkDir, LogDir, CacheDir}
	for _, p := range dirs {
		exists, err := PathExists(p)
		if err != nil {
			log.Fatal(err)
		}

		if !exists {
			if err = os.MkdirAll(p, 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func loggerInit() {
	var err error
	LogFile, err = os.OpenFile(path.Join(LogDir, "penguin.log"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("fail to create log file:", err)
	}

	Logger = log.New(LogFile, "", log.LstdFlags|log.Lshortfile)
}
