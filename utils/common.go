package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func GetScriptString(path string) string {
	absScriptPath, err := filepath.Abs(path)
	HandleError(nil, err, "Error getting absolute path for Lua script: ")
	dat, err := os.ReadFile(absScriptPath)
	HandleError(nil, err)
	return string(dat[:])
}

type fnT func()

func BenchmarkFn(fn fnT) {
	start := time.Now()
	fn()
	timeTaken := time.Since(start).Seconds()
	log.Printf("Time taken: %.9f seconds.\n", timeTaken)
}

func ConvertToMap(key string, data interface{}) map[string]interface{} {
	mapData := make(map[string]interface{})
	mapData[key] = data
	return mapData
}
