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

func TestStringsParameter(t *testing.T) {

	queryString := "http://www.domain.com/search?interest=alfa,beta,gamma,delta"

	parser := NewParser()
	interestParameter := NewParameter("interest")
	interestParameter.OutputName = "profile.interest"
	interestParameter.Type = Strings
	parser.AddParameter(interestParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	expectedStrings := []string{"alfa", "beta", "gamma", "delta"}
	if !testEq(parser.Parameters[0].StringsValue, expectedStrings) {
		t.Fail()
	}
}

func TestEmptyStringsParameter(t *testing.T) {

	queryString := "http://www.domain.com/search?interest="

	parser := NewParser()
	interestParameter := NewParameter("interest")
	interestParameter.OutputName = "profile.interest"
	interestParameter.Type = Strings
	parser.AddParameter(interestParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	expectedStrings := []string{}
	generatedStrings := parser.Parameters[0].StringsValue
	if !testEq(generatedStrings, expectedStrings) {
		t.Errorf("Expected '%v' got '%v'", expectedStrings, generatedStrings)
	}
}

func TestIntegerRangeParameter(t *testing.T) {

	queryString := "http://www.domain.com/search?age=15-35"

	parser := NewParser()
	ageParameter := NewParameter("age")
	ageParameter.Type = IntegerRange
	ageParameter.MinValue = 10
	ageParameter.MaxValue = 99
	parser.AddParameter(ageParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	p := parser.Parameters[0]
	if p.MinValue != 15 || p.MaxValue != 35 {
		t.Fail()
	}
}

func TestInvalidIntegerRangeParameter(t *testing.T) {

	queryString := "http://www.domain.com/search?age=45-35"

	parser := NewParser()
	ageParameter := NewParameter("age")
	ageParameter.Type = IntegerRange
	ageParameter.MinValue = 10
	ageParameter.MaxValue = 99
	parser.AddParameter(ageParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	p := parser.Parameters[0]
	if p.MinValue != 35 || p.MaxValue != 45 {
		t.Fail()
	}
}

func TestSearchStringParameter(t *testing.T) {
	queryString := "https://www.domain.com/search?q=alfa*"

	parser := NewParser()
	searchStringParameter := NewParameter("q")
	searchStringParameter.OutputNames = []string{"name", "lastname", "about"}
	searchStringParameter.Type = SearchString
	searchStringParameter.MaxLength = 80
	parser.AddParameter(searchStringParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	p := parser.Parameters[0]
	if p.StringValue != "alfa" || p.Position != Suffix {
		t.Fail()
	}
}

func TestSearchStringParameterSurround(t *testing.T) {
	queryString := "https://www.domain.com/search?q=*beta*"

	parser := NewParser()
	searchStringParameter := NewParameter("q")
	searchStringParameter.OutputNames = []string{"name", "lastname", "about"}
	searchStringParameter.Type = SearchString
	searchStringParameter.MaxLength = 4
	parser.AddParameter(searchStringParameter)

	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	p := parser.Parameters[0]
	if p.StringValue != "beta" || p.Position != Surrounded {
		t.Fail()
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
	offsetParameter.Output = false
	parser.AddParameter(offsetParameter)

	sizeParameter := NewParameter("size")
	sizeParameter.Type = Integer
	sizeParameter.DefaultIntValue = 50
	sizeParameter.MinValue = 50
	sizeParameter.MaxValue = 500
	sizeParameter.Output = false
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

// https://stackoverflow.com/a/15312097/254695
func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
