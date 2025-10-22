package util

import "log"

func Info(v ...any) {
	log.Println("[INFO]", v)
}

func Success(v ...any) {
	log.Println("[SUCCESS]", v)
}

func Error(v ...any) {
	log.Println("[ERROR]", v)
}
