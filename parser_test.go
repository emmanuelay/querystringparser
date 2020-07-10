package querystringparser

import (
	"testing"
)

func TestNoQuery(t *testing.T) {

	parser := NewParser()
	interestParameter := NewParameter("interest")
	interestParameter.OutputName = "profile.interest"
	interestParameter.Type = Strings
	parser.AddParameter(interestParameter)

	queryStringNoQuery := "http://www.domain.com/search"
	err := parser.Parse(queryStringNoQuery)
	if err != ErrNoQueryString {
		t.Fail()
	}

	queryStringEmptyQuery := "http://www.domain.com/search?"
	err = parser.Parse(queryStringEmptyQuery)
	if err != ErrNoQueryString {
		t.Fail()
	}
}

func TestNoParameter(t *testing.T) {
	queryString := "http://www.domain.com/search?interest=alfa,beta,gamma,delta"

	parser := NewParser()
	err := parser.Parse(queryString)
	if err == nil {
		t.Fail()
	}
}

func TestUnregisteredParameter(t *testing.T) {
	queryString := "http://www.domain.com/search?age=18-35&other=notreally"

	parser := NewParser()
	ageRangeParameter := NewParameter("age")
	ageRangeParameter.Type = IntegerRange
	ageRangeParameter.MinValue = 18
	ageRangeParameter.MaxValue = 80
	parser.AddParameter(ageRangeParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 1 {
		t.Errorf("Invalid number of parsed parameters (%v, expected 1)", parser.ParsedParameterCount())
	}
}

func TestIntegerParameter(t *testing.T) {

	queryString := "https://www.domain.com/search?q=alfa*&offset=1&size=50"

	parser := NewParser()

	searchStringParameter := NewParameter("q")
	searchStringParameter.OutputNames = []string{"name", "lastname", "about"}
	searchStringParameter.Type = SearchString
	searchStringParameter.MaxLength = 4
	parser.AddParameter(searchStringParameter)

	offsetParameter := NewParameter("offset")
	offsetParameter.Type = Integer
	offsetParameter.DefaultIntValue = 0
	offsetParameter.MinValue = 0
	offsetParameter.MaxValue = 1000
	offsetParameter.IncludeInOutput = false
	parser.AddParameter(offsetParameter)

	sizeParameter := NewParameter("size")
	sizeParameter.Type = Integer
	sizeParameter.DefaultIntValue = 50
	sizeParameter.MinValue = 50
	sizeParameter.MaxValue = 500
	sizeParameter.IncludeInOutput = false
	parser.AddParameter(sizeParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	val, err := parser.GetIntValue("size")
	if err != nil {
		t.Fail()
	}

	if val != 50 {
		t.Fail()
	}
}
