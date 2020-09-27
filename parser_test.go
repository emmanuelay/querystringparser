package querystringparser

import (
	"testing"
)

func TestNoQuery(t *testing.T) {

	parser := NewParser()
	interestParameter := NewParameter("interest", Strings)
	interestParameter.OutputName = "profile.interest"
	parser.AddParameter(interestParameter)

	queryStringNoQuery := "http://www.domain.com/search"
	err := parser.Parse(queryStringNoQuery)
	if err != nil {
		t.Error("Expected Error")
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
	ageRangeParameter := NewParameter("age", IntegerRange)
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

func TestOnlyParameters(t *testing.T) {
	queryString := "age=18-35&other=notreally"

	parser := NewParser()
	ageRangeParameter := NewParameter("age", IntegerRange)
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

	searchStringParameter := NewParameter("q", SearchString)
	searchStringParameter.OutputNames = []string{"name", "lastname", "about"}
	searchStringParameter.MaxLength = 4
	parser.AddParameter(searchStringParameter)

	offsetParameter := NewParameter("offset", Integer)
	offsetParameter.DefaultIntValue = 0
	offsetParameter.MinValue = 0
	offsetParameter.MaxValue = 1000
	offsetParameter.IncludeInOutput = false
	parser.AddParameter(offsetParameter)

	sizeParameter := NewParameter("size", Integer)
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

func TestBooleanParameter(t *testing.T) {
	queryString := "active=true&other=notsomuch"

	parser := NewParser()
	activeParameter := NewParameter("active", Boolean)
	parser.AddParameter(activeParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 1 {
		t.Errorf("Invalid number of parsed parameters (%v, expected 1)", parser.ParsedParameterCount())
	}
}

func TestDefaultIntegerParameter(t *testing.T) {

	queryString := "https://www.domain.com/search?offset=1"

	parser := NewParser()

	offsetParameter := NewParameter("offset", Integer)
	offsetParameter.DefaultIntValue = 0
	offsetParameter.MinValue = 0
	offsetParameter.MaxValue = 1000
	offsetParameter.IncludeInOutput = false
	parser.AddParameter(offsetParameter)

	sizeParameter := NewParameter("size", Integer)
	sizeParameter.DefaultIntValue = 60
	sizeParameter.MinValue = 50
	sizeParameter.MaxValue = 500
	sizeParameter.IncludeInOutput = false
	parser.AddParameter(sizeParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	size, err := parser.GetIntValue("size")
	if err != nil {
		t.Fail()
	}

	if size != 60 {
		t.Fail()
	}

	offset, err := parser.GetIntValue("offset")
	if err != nil {
		t.Fail()
	}

	if offset != 1 {
		t.Fail()
	}
}

func TestSanitizeKey(t *testing.T) {
	name := "alFaBetic_tHInG åÅäÖöÄ !$\"€%/'--&//()#)€#!_123_ '--DROP TABLE"
	output := sanitizeKey(name)
	expectedOutput := "alfabetic_thing_123_droptable"
	if output != expectedOutput {
		t.Errorf("Expected '%v' got '%v'", expectedOutput, output)
	}
}

func TestUnsanitizedParameter(t *testing.T) {
	queryString := "https://www.domain.com/search?offSET=10"

	parser := NewParser()

	offsetParameter := NewParameter("offset", Integer)
	offsetParameter.DefaultIntValue = 0
	offsetParameter.MinValue = 0
	offsetParameter.MaxValue = 1000
	offsetParameter.IncludeInOutput = false
	parser.AddParameter(offsetParameter)

	err := parser.Parse(queryString)
	if err != ErrInvalidKeyName {
		t.Error("Should fail due to unsanitized key name")
	}

}
