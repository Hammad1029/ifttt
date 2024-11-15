package requestvalidator

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-viper/mapstructure/v2"
)

type ParameterConfig interface {
	Generate() (string, error)
	GetMap() map[string]any
}

func (r *RequestParameter) Generate() (string, error) {
	var generator ParameterConfig
	switch r.DataType {
	case dataTypeText:
		generator = &textValue{}
	case dataTypeNumber:
		generator = &numberValue{}
	case dataTypeBoolean:
		generator = &booleanValue{}
	case dataTypeArray:
		generator = &arrayValue{}
	case dataTypeMap:
		generator = &mapValue{}
	default:
		return "", fmt.Errorf("datatype %s not found", r.DataType)
	}
	if err := mapstructure.Decode(&r.Config, generator); err != nil {
		return "", fmt.Errorf("datatype %s could not be decoded", r.DataType)
	}
	if generatedRegex, err := generator.Generate(); err != nil {
		return "", fmt.Errorf("error in generating regex for datatype %s: %s", r.DataType, err)
	} else {
		if r.DataType == dataTypeArray || r.DataType == dataTypeMap {
			r.Config = generator.GetMap()
		}
		return generatedRegex, nil
	}
}

func (t *textValue) Generate() (string, error) {
	var patterns []string

	var charClass strings.Builder
	charClass.WriteString("[")
	if t.Alpha {
		charClass.WriteString(alphaRegex)
	}
	if t.Numeric {
		charClass.WriteString(numericRegex)
	}

	if t.Special {
		charClass.WriteString(specialRegex)
	}
	charClass.WriteString("]")

	patterns = append(patterns, charClass.String())

	if t.Minimum > 0 || t.Maximum > 0 {
		patterns = append(patterns, t.minMax(t.Minimum, t.Maximum))
	}

	if len(t.In) > 0 {
		patterns = append(patterns, t.in(t.In))
	}

	if len(t.NotIn) > 0 {
		patterns = append(patterns, t.notIn(t.NotIn))
	}

	return strings.Join(patterns, ""), nil
}

func (t *textValue) GetMap() map[string]any {
	return structs.Map(*t)
}

func (n *numberValue) Generate() (string, error) {
	var patterns []string
	patterns = append(patterns, numberRegex)
	patterns = append(patterns, n.minMax(n.Minimum, n.Maximum))
	patterns = append(patterns, n.in(n.In))
	patterns = append(patterns, n.notIn(n.NotIn))
	return strings.Join(patterns, ""), nil
}

func (n *numberValue) GetMap() map[string]any {
	return structs.Map(*n)
}

func (b *booleanValue) Generate() (string, error) {
	return booleanRegex, nil
}

func (b *booleanValue) GetMap() map[string]any {
	return structs.Map(*b)
}

func (a *arrayValue) Generate() (string, error) {
	elementRegex, err := a.OfType.Generate()
	if err != nil {
		return "", err
	}
	a.OfType.Regex = elementRegex
	return "", nil
}

func (a *arrayValue) GetMap() map[string]any {
	return structs.Map(*a)
}

func (m *mapValue) Generate() (string, error) {
	for key, value := range *m {
		valueRegex, err := value.Generate()
		if err != nil {
			return "", err
		}
		value.Regex = valueRegex
		(*m)[key] = value
	}
	return "", nil
}

func (m *mapValue) GetMap() map[string]any {
	mapConverted := map[string]any{}
	for key, val := range *m {
		mapConverted[key] = val
	}
	return mapConverted
}
