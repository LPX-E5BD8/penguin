package model

import (
	"log"
	"os"
	"path"
)

// base dir
var WorkDir = path.Join(os.Getenv("HOME"), ".penguin")
// log dir
var LogDir = path.Join(WorkDir, "log")
// cache dir
var CacheDir = path.Join(WorkDir, "cache")
// Logger
var Logger *log.Logger

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
	logFile, err := os.OpenFile(path.Join(LogDir, "penguin.log"), os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("fail to create log file:", err)
	}

	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}
