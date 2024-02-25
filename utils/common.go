package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetScriptString(path string) string {
	absScriptPath, err := filepath.Abs(path)
	HandleError(err, "Error getting absolute path for Lua script: ")
	dat, err := os.ReadFile(absScriptPath)
	HandleError(err)
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

func ConvertStringToInterfaceArray(obj []string) []interface{} {
	s := make([]interface{}, len(obj))
	for i, v := range obj {
		s[i] = v
	}
	return s
}

func GenerateIndexName(tableName string, columns []string) string {
	// tables_description_idx_index
	return fmt.Sprintf("%s_%s_idx", tableName[:30], strings.Join(columns, "_"))
}
