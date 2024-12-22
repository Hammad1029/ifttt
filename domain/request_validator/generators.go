package requestvalidator

import (
	"fmt"

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
	var regex string

	if len(t.In) > 0 {
		return t.in(), nil
	}

	charClass := "["
	if t.Alpha {
		charClass += alphaRegex
	}
	if t.Numeric {
		charClass += numericRegex
	}
	if t.Special {
		charClass += specialRegex
	}
	charClass += "]"
	regex += charClass

	regex += t.minMax()

	return regex, nil
}

func (t *textValue) GetMap() map[string]any {
	return structs.Map(*t)
}

func (n *numberValue) Generate() (string, error) {
	var regex string
	regex += fmt.Sprintf("[%s]+", numericRegex)
	regex += n.in()
	return regex, nil
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
