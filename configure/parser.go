package configure

import (
	"reflect"
	"strings"
)

type ConfigParser struct {
	path 		string
	configType 	string
}

func NewConfigParser(path string, format string) *ConfigParser {
	return &ConfigParser{
		path: path,
		configType: format,
	}
}

func (cp *ConfigParser) Parser(target interface{}) error {
	if strings.Trim(cp.path, " ") == "" {
		return &ConfigParseError{EmptyConfigPath}
	}
	if strings.Trim(cp.configType, " ") == "" {
		return &ConfigParseError{EmptyConfigFormat}
	}
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &ConfigParseError{TargetShouldBePointer}
	}

	switch cp.configType {
	case "json":
		return cp.parserJson(cp.path)
	case "yaml":
		return cp.parserYaml(cp.path)
	case "toml":
		return cp.parserToml(cp.path)
	case "ini":
		return cp.parserIni(cp.path)
	case "conf":
		return cp.parserConf(cp.path)
	default:
		return &ConfigParseError{UnsupportedType}
	}
}

func (cp *ConfigParser) parserJson(path string) error {
	return nil
}

func (cp *ConfigParser) parserYaml(path string) error {
	return nil
}

func (cp *ConfigParser) parserToml(path string) error {
	return nil
}

func (cp *ConfigParser) parserIni(path string) error {
	return nil
}

func (cp *ConfigParser) parserConf(path string) error {
	return nil
}

const (
	UnsupportedType = iota
	TargetShouldBePointer
	EmptyConfigPath
	EmptyConfigFormat
)

// config parse error
type ConfigParseError struct {
	ErrType 		int
}

func (ce *ConfigParseError) Error() string {
	base := "config parse error,"
	switch ce.ErrType {
	case UnsupportedType:
		return base + "unsupported type"
	case TargetShouldBePointer:
		return base + "target should be pointer"
	case EmptyConfigPath:
		return base + "empty config path"
	case EmptyConfigFormat:
		return base + "empty config format"
	default:
		return base + "unknown"
	}
}
