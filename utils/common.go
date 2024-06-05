package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

func BenchmarkFn(fn func()) {
	start := time.Now()
	fn()
	fmt.Printf("execution time: %+v\n", time.Since(start))
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

func StringifyMapInt(m map[string]interface{}) (map[string]interface{}, error) {
	for key, val := range m {
		if reflect.TypeOf(val).Kind() == reflect.Map {
			if data, err := json.Marshal(val); err != nil {
				return nil, err
			} else {
				m[key] = string(data[:])
			}
		} else {
			m[key] = fmt.Sprint(val)
		}
	}
	return m, nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
