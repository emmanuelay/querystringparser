package querystringparser

import (
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
	switch p.Type {

	case Integer:
		return fmt.Sprintf("%v:%v", p.OutputName, p.IntValue), nil

	case IntegerRange:
		return fmt.Sprintf("%v:>=%v %v:<=%v", p.OutputName, p.MinValue, p.Name, p.MaxValue), nil

	case Strings:
		return fmt.Sprintf("%v:%v", p.OutputName, strings.Join(p.StringsValue, ",")), nil

	case SearchString:
		{
			fieldName := p.OutputName
			if len(fieldName) > 0 {
				fieldName = fmt.Sprintf("%v:", fieldName)
			}

			switch p.Position {

			case Prefix:
				return fmt.Sprintf("%v%v*", fieldName, p.StringValue), nil

			case Suffix:
				return fmt.Sprintf("%v*%v", fieldName, p.StringValue), nil

			default:
				return fmt.Sprintf("%v*%v*", fieldName, p.StringValue), nil
			}
		}
	}
	return "N/A", nil
}

// ToBleveSortSlice returns a slice of strings to be used as a sort modifier by Bleve
func (p *Parameter) ToBleveSortSlice() ([]string, error) {
	// TODO(ea): Return string slice with sort modifier
	// TODO(ea): Example: []string{"age", "-_score", "_id"}

	return nil, nil
}
