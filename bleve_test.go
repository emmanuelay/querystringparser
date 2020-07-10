package querystringparser

import (
	"fmt"
	"testing"
)

func TestToBleveQuery(t *testing.T) {

	parser := NewParser()

	interestParameter := NewParameter("interests")
	interestParameter.Type = Strings
	interestParameter.OutputName = "profile.interest"
	parser.AddParameter(interestParameter)

	villageParameter := NewParameter("villages")
	villageParameter.Type = Strings
	villageParameter.OutputName = "profile.villages"
	parser.AddParameter(villageParameter)

	ageParameter := NewParameter("age")
	ageParameter.Type = IntegerRange
	ageParameter.OutputName = "age"
	ageParameter.MinValue = 0
	ageParameter.MaxValue = 99
	parser.AddParameter(ageParameter)

	//registrationDateParameter := NewParameter("reg")
	//registrationDateParameter.OutputName = "created_at"
	//registrationDateParameter.Type = type.DateRange
	//registrationDateParameter.DateFormat = "yyyyMMdd"
	//parser.AddParameter(registrationDateParameter)

	searchStringParameter := NewParameter("q")
	searchStringParameter.Type = SearchString
	searchStringParameter.OutputName = ""
	searchStringParameter.MaxLength = 80
	parser.AddParameter(searchStringParameter)

	offsetParameter := NewParameter("offset")
	offsetParameter.Type = Integer
	offsetParameter.DefaultIntValue = 0
	offsetParameter.MinValue = 0
	offsetParameter.MaxValue = 999
	offsetParameter.IncludeInOutput = false
	parser.AddParameter(offsetParameter)

	sizeParameter := NewParameter("size")
	sizeParameter.Type = Integer
	sizeParameter.DefaultIntValue = 50
	sizeParameter.MinValue = 50
	sizeParameter.MaxValue = 500
	sizeParameter.IncludeInOutput = false
	parser.AddParameter(sizeParameter)

	sortParameter := NewParameter("sort")
	sortParameter.Type = SortStrings
	sortParameter.AllowedValues = []string{"age", "name", "last_online"}
	sortParameter.IncludeInOutput = false
	parser.AddParameter(sortParameter)

	queryString := "http://www.domain.com/search?q=hello*&age=18-45&villages=alfa,beta&interests=gamma,delta&sort=age,name,-last_online&size=100&offset=50"
	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 7 {
		t.Error("Invalid number of processed parameters")
	}

	query, err := parser.ToBleveQuery()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("in :", queryString)
	fmt.Println("out:", query)

	intSize, err := parser.GetIntValue("size")
	if err != nil {
		t.Error(err)
	}

	if intSize != 100 {
		t.Error("Size should be 100")
	}

	intOffset, err := parser.GetIntValue("offset")
	if err != nil {
		t.Error(err)
	}

	if intOffset != 50 {
		t.Error("Offset should be 50")
	}
}
