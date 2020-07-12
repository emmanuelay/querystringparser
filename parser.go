package querystringparser

import (
	"errors"
	"fmt"
	"strings"
)

// Parser ...
type Parser struct {
	Parameters         []Parameter
	ParameterSeparator string
	KeyValueSeparator  string
}

const (
	querySeparator              = "?"
	parameterSeparatorCharacter = "&"
	keyValueSeparatorCharacter  = "="
	wildCardCharacter           = "*"
	rangeSeparatorCharacter     = "-"
	listSeparatorCharacter      = ","
	sortModifierCharacter       = "-"
)

var (
	// ErrNoParameters ...
	ErrNoParameters = errors.New("No parameters to parse")

	// ErrNoParameter ...
	ErrNoParameter = errors.New("Could not find parameter")

	// ErrNoQueryString ...
	ErrNoQueryString = errors.New("No querystring to parse")

	// ErrInvalidType ...
	ErrInvalidType = errors.New("Invalid type")

	// ErrInvalidRange ...
	ErrInvalidRange = errors.New("Invalid range parameter")
)

// NewParser creates a Parser-instance
func NewParser() *Parser {
	return &Parser{
		ParameterSeparator: parameterSeparatorCharacter,
		KeyValueSeparator:  keyValueSeparatorCharacter,
	}
}

// AddParameter adds a parameter to the parser
func (p *Parser) AddParameter(parameter Parameter) {
	if p.Parameters == nil {
		parameters := []Parameter{parameter}
		p.Parameters = parameters
		return
	}

	parameters := append(p.Parameters, parameter)
	p.Parameters = parameters
}

// Parse performs a parse of the queryString
func (p *Parser) Parse(queryString string) error {

	if len(p.Parameters) == 0 {
		return ErrNoParameters
	}

	paramString := queryString

	// http://www.domain.com/search <-> ?parameter=value
	if strings.Contains(queryString, querySeparator) {
		query := strings.Split(queryString, querySeparator)
		if len(query) == 1 || len(query[1]) == 0 {
			return ErrNoQueryString
		}
		paramString = query[1]
	}

	// http://www.domain.com/search? parameter=value <-> &parameter2=value
	queryParameters := strings.Split(paramString, p.ParameterSeparator)
	if len(queryParameters) == 0 {
		return ErrNoParameters
	}

	for _, queryParameter := range queryParameters {
		keyValue := strings.Split(queryParameter, p.KeyValueSeparator)

		if len(keyValue) < 2 {
			continue
		}

		key := keyValue[0]
		value := keyValue[1]

		parameter, err := p.getParameter(key)
		if err != nil {
			continue
		}

		err = parameter.Parse(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetIntValue returns the integer value for the parameter with name 'key'
func (p *Parser) GetIntValue(key string) (int, error) {
	parameter, err := p.getParameter(key)
	if err != nil {
		return -1, err
	}

	if parameter.Type != Integer {
		return -1, fmt.Errorf("Invalid parameter type for parameter '%v' (expected Integer)", parameter.Name)
	}

	if !parameter.Parsed {
		parameter.IntValue = parameter.DefaultIntValue
	}

	return parameter.IntValue, nil
}

func (p *Parser) getParameter(key string) (*Parameter, error) {
	for idx, parameter := range p.Parameters {
		if parameter.Name == key {
			// Return reference
			return &p.Parameters[idx], nil
		}
	}
	return nil, ErrNoParameter
}

// ParsedParameterCount returns the number of parsed parameters
func (p *Parser) ParsedParameterCount() int {
	counter := 0
	for _, parameter := range p.Parameters {
		if parameter.Parsed {
			counter++
		}
	}
	return counter
}
