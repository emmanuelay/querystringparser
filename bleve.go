package querystringparser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// ToBleveQuery returns a Bleve-compatible search query for the parsed parameters
func (p *Parser) ToBleveQuery() (string, error) {

	output := []string{}

	for _, parameter := range p.Parameters {

		if !parameter.Parsed || !parameter.IncludeInOutput {
			continue
		}

		bleveQuery, err := parameter.ToBleveQuery()
		if err != nil {
			return "", err
		}

		if len(bleveQuery) > 0 {
			output = append(output, bleveQuery)
		}
	}

	return strings.Join(output, " "), nil
}

// ToBleveQuery returns a Bleve-compatible query parameter
func (p *Parameter) ToBleveQuery() (string, error) {

	conditionalModifier := ""
	switch p.OutputCondition {
	case Must:
		conditionalModifier = "+"
	case Not:
		conditionalModifier = "-"
	}

	switch p.Type {

	case Integer:
		return fmt.Sprintf("%v%v:%v", conditionalModifier, p.OutputName, p.IntValue), nil

	case Boolean:
		return fmt.Sprintf("%v%v:%v", conditionalModifier, p.OutputName, p.BoolValue), nil

	case IntegerRange:
		return fmt.Sprintf("%v%v:>=%v %v%v:<=%v", conditionalModifier, p.OutputName, p.MinValue, conditionalModifier, p.Name, p.MaxValue), nil

	case Strings:
		{
			if p.OutputCondition == Should {
				return fmt.Sprintf("%v%v:%v", conditionalModifier, p.OutputName, strings.Join(p.StringsValue, ",")), nil
			}

			var query bytes.Buffer
			for _, stringValue := range p.StringsValue {
				if query.Len() > 0 {
					query.WriteString(" ")
				}
				condition := fmt.Sprintf("%v%v:%v", conditionalModifier, p.OutputName, stringValue)
				query.WriteString(condition)
			}
			return query.String(), nil
		}

	case SearchString:
		{
			fieldName := p.OutputName
			if len(fieldName) > 0 {
				fieldName = fmt.Sprintf("%v:", fieldName)
			}

			switch p.Position {

			case Prefix:
				return fmt.Sprintf("%v%v%v*", conditionalModifier, fieldName, p.StringValue), nil

			case Suffix:
				return fmt.Sprintf("%v%v*%v", conditionalModifier, fieldName, p.StringValue), nil

			default:
				return fmt.Sprintf("%v%v*%v*", conditionalModifier, fieldName, p.StringValue), nil
			}
		}
	}
	return "N/A", nil
}

// ErrInvalidParameter ...
var ErrInvalidParameter = errors.New("Invalid parameter type, expected 'SortStrings'")

// ToBleveSortSlice returns a slice of strings to be used as a sort modifier by Bleve
func (p *Parameter) ToBleveSortSlice() ([]string, error) {
	if p.Type != SortStrings {
		return nil, ErrInvalidParameter
	}

	output := []string{}
	for idx, val := range p.StringsValue {
		sortDirection := p.SortDirections[idx]
		var item string
		if sortDirection == false {
			item = fmt.Sprintf("%v%v", p.SortModifierCharacter, val)
		} else {
			item = val
		}
		output = append(output, item)
	}

	return output, nil
}

// ToBleveSortSlice retrieves a string-slice in a Bleve-compatible format (SortBy)
func (p *Parser) ToBleveSortSlice(sortParameterName string) ([]string, error) {
	parameter, err := p.getParameter(sortParameterName)
	if err != nil {
		return nil, err
	}

	sortSlice, err := parameter.ToBleveSortSlice()
	return sortSlice, err
}
