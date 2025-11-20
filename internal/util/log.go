package utils

import (
	"fmt"
	"time"
)

var DebugEnabled = false
var IsProduction = true // NEW

func Info(msg string, args ...interface{}) {
	if IsProduction {
		return // hide info logs in production
	}
	fmt.Printf("🟦 INFO %s | %s\n",
		time.Now().Format(time.RFC3339),
		fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...interface{}) {
	fmt.Printf("🟨 WARN %s | %s\n",
		time.Now().Format(time.RFC3339),
		fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...interface{}) {
	fmt.Printf("🟥 ERROR %s | %s\n",
		time.Now().Format(time.RFC3339),
		fmt.Sprintf(msg, args...))
}

func Debug(msg string, args ...interface{}) {
	if DebugEnabled {
		fmt.Printf("🟪 DEBUG %s | %s\n",
			time.Now().Format(time.RFC3339),
			fmt.Sprintf(msg, args...))
	}
}
