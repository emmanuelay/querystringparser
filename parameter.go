package querystringparser

import (
	"fmt"
	"strconv"
	"strings"
)

// Type is an enum to denote different types of parameters
type Type int

const (
	// Strings type is a delimited array of strings
	// Ex: interest=friends,soccer,software
	Strings Type = iota

	// SearchString type is a wildcard pre/sufixed string
	SearchString

	// SortStrings type is a delimited array of strings that can have a directional modifier
	// Ex: sort=name,-age,height -> This sorts 'name' ascending, 'age' descending, 'height' ascending
	SortStrings

	// IntegerRange type is a parameter that restricts input to an integer range
	// Ex: age=18-30 -or- age=-30 -or- age=18-
	IntegerRange

	// Integer type is an integer with restrictions
	Integer
)

// MatchPosition denotes where in a search string the wildcard is located
type MatchPosition int

const (
	// Prefix is in the beginning of the search string (ex. *search)
	Prefix MatchPosition = iota

	// Suffix is in the end of the search string (ex. search*)
	Suffix

	// Surrounded is in the beginning and the end of the search string (ex. *search*)
	Surrounded
)

// Parameter ...
type Parameter struct {
	Name       string
	OutputName string
	Type       Type
	Output     bool

	// Range specific variables
	RangeSeparatorCharacter string

	// Defaults
	DefaultIntValue int
	DefaultMinValue int
	DefaultMaxValue int

	// String specific variables
	StringValue       string
	Position          MatchPosition
	WildCardCharacter string
	MinLength         int
	MaxLength         int
	OutputNames       []string

	// Strings specific variables
	ListSeparatorCharacter string
	SortModifierCharacter  string
	StringsValue           []string
	SortDirections         []bool // true = ascending, false = descending

	// Integer specific variables
	IntValue int
	MinValue int
	MaxValue int
}

// NewParameter creates a new parameter with default configuration
func NewParameter(parameter string) Parameter {
	return Parameter{
		Name:                    parameter,
		Output:                  true,
		OutputName:              parameter,
		WildCardCharacter:       wildCardCharacter,
		RangeSeparatorCharacter: rangeSeparatorCharacter,
		ListSeparatorCharacter:  listSeparatorCharacter,
		SortModifierCharacter:   sortModifierCharacter,
		MinLength:               1,
		MaxLength:               100,
	}
}

// Parse performs a parameter parse of a key/value-pair
func (p *Parameter) Parse(key, value string) error {

	switch p.Type {
	case Strings:
		return p.parseStrings(key, value)
	case SortStrings:
		return p.parseSortStrings(key, value)
	case IntegerRange:
		return p.parseIntegerRange(key, value)
	case SearchString:
		return p.parseSearchString(key, value)
	case Integer:
		return p.parseInteger(key, value)
	default:
		return ErrInvalidType
	}
}

func (p *Parameter) parseStrings(key, value string) error {
	items := strings.Split(value, p.ListSeparatorCharacter)
	if len(items) == 1 && len(items[0]) == 0 {
		p.StringsValue = []string{}
	} else {
		p.StringsValue = items
	}
	return nil
}

func (p *Parameter) parseSortStrings(key, value string) error {
	items := strings.Split(value, p.ListSeparatorCharacter)
	if len(items) == 1 && len(items[0]) == 0 {
		p.StringsValue = []string{}
		p.SortDirections = []bool{}
	}

	outputItems := []string{}
	outputDirections := []bool{}

	for _, item := range items {
		if strings.HasPrefix(item, p.SortModifierCharacter) {
			outputItems = append(outputItems, strings.ReplaceAll(item, p.SortModifierCharacter, ""))
			outputDirections = append(outputDirections, false)
		} else {
			outputItems = append(outputItems, item)
			outputDirections = append(outputDirections, true)
		}
	}

	p.StringsValue = outputItems
	p.SortDirections = outputDirections

	return nil
}

func (p *Parameter) parseIntegerRange(key, value string) error {
	rangePair := strings.Split(value, p.RangeSeparatorCharacter)
	if len(rangePair) == 1 {
		return ErrInvalidRange
	}

	minRange := rangePair[0]
	maxRange := rangePair[1]

	if len(minRange) > 0 {
		min, err := strToint(minRange)
		if err == ErrInvalidType {
			return fmt.Errorf("Invalid type in min-range value '%v' for parameter '%v'", minRange, p.Name)
		}
		p.MinValue = min
	}

	if len(maxRange) > 0 {
		max, err := strToint(maxRange)
		if err == ErrInvalidType {
			return fmt.Errorf("Invalid type in min-range value '%v' for parameter '%v'", maxRange, p.Name)
		}
		p.MaxValue = max
	}

	if p.MinValue > p.MaxValue {
		rMin := p.MinValue
		rMax := p.MaxValue
		p.MinValue = rMax
		p.MaxValue = rMin
	}

	return nil
}

func (p *Parameter) parseSearchString(key, value string) error {
	// Determine position
	hasPrefix := strings.HasPrefix(value, p.WildCardCharacter)
	hasSuffix := strings.HasSuffix(value, p.WildCardCharacter)

	// Default to surrounded
	p.Position = Surrounded

	if hasPrefix && !hasSuffix {
		p.Position = Prefix
	}

	if !hasPrefix && hasSuffix {
		p.Position = Suffix
	}

	strValue := strings.ReplaceAll(value, p.WildCardCharacter, "")
	if p.MaxLength > 0 && len(strValue) > p.MaxLength {
		p.StringValue = strValue[:p.MaxLength]
		return fmt.Errorf("Invalid length (%v) for parameter '%v' (max %v)", len(strValue), p.Name, p.MaxLength)
	}

	if len(strValue) < p.MinLength {
		p.StringValue = strValue
		return fmt.Errorf("Invalid length (%v) for parameter '%v' (min %v)", len(strValue), p.Name, p.MinLength)
	}

	p.StringValue = strValue

	return nil
}

func (p *Parameter) parseInteger(key, value string) error {
	val, err := strToint(value)
	if err == ErrInvalidType {
		return fmt.Errorf("Invalid type in integer value '%v' for parameter '%v'", val, p.Name)
	}

	if p.MaxValue != 0 && val > p.MaxValue {
		val = p.MaxValue
	}

	if val < p.MinValue {
		val = p.MinValue
	}

	p.IntValue = val

	return nil
}

func strToint(input string) (int, error) {
	intValue, ok := strconv.Atoi(input)
	if ok != nil {
		return 0, ErrInvalidType
	}
	return intValue, nil
}
