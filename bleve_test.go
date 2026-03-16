package querystringparser

import (
	"testing"
)

func TestToBleveQuery(t *testing.T) {

	parser := NewParser()

	searchStringParameter := NewParameter("q", SearchString)
	searchStringParameter.OutputName = ""
	searchStringParameter.MaxLength = 80
	searchStringParameter.OutputCondition = Must
	parser.AddParameter(searchStringParameter)

	activeParameter := NewParameter("active", Boolean)
	activeParameter.OutputCondition = Must
	parser.AddParameter(activeParameter)

	ageParameter := NewParameter("age", IntegerRange)
	ageParameter.MinValue = 0
	ageParameter.MaxValue = 99
	ageParameter.OutputCondition = Must
	parser.AddParameter(ageParameter)

	villageParameter := NewParameter("villages", Strings)
	villageParameter.OutputName = "profile.villages"
	villageParameter.OutputCondition = Must
	parser.AddParameter(villageParameter)

	interestParameter := NewParameter("interests", Strings)
	interestParameter.OutputName = "profile.interest"
	interestParameter.OutputCondition = Should
	parser.AddParameter(interestParameter)

	registrationDateParameter := NewParameter("reg", DateRange)
	registrationDateParameter.OutputName = "created_at"
	registrationDateParameter.OutputCondition = Must
	parser.AddParameter(registrationDateParameter)

	offsetParameter := NewParameter("offset", Integer)
	offsetParameter.DefaultIntValue = 0
	offsetParameter.MinValue = 0
	offsetParameter.MaxValue = 999
	offsetParameter.IncludeInOutput = false
	parser.AddParameter(offsetParameter)

	sizeParameter := NewParameter("size", Integer)
	sizeParameter.DefaultIntValue = 50
	sizeParameter.MinValue = 50
	sizeParameter.MaxValue = 500
	sizeParameter.IncludeInOutput = false
	parser.AddParameter(sizeParameter)

	sortParameter := NewParameter("sort", SortStrings)
	sortParameter.AllowedValues = []string{"age", "name", "last_online"}
	sortParameter.IncludeInOutput = false
	parser.AddParameter(sortParameter)

	queryString := "http://www.domain.com/search?q=*hello*&age=18-45&active=T&villages=alfa,beta&interests=gamma,delta&reg=20200101-20200304&sort=-age,name,-last_online&size=100&offset=99150"
	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 9 {
		t.Error("Invalid number of processed parameters")
	}

	query, err := parser.ToBleveQuery()
	if err != nil {
		t.Error(err)
	}

	bleveString := "+*hello* +active:true +age:>=18 +age:<=45 +profile.villages:alfa +profile.villages:beta +(profile.interest:gamma profile.interest:delta) +created_at:>=20200101 +created_at:<=20200304"
	if query != bleveString {
		t.Errorf("Expected '%v' got '%v'", bleveString, query)
	}

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

	if intOffset != 999 {
		t.Error("Offset should be 999")
	}
}

func TestToBleveQueryDateRangeImplicitMin(t *testing.T) {
	parser := NewParser()

	regParameter := NewParameter("reg", DateRange)
	regParameter.OutputName = "created_at"
	regParameter.OutputCondition = Must
	parser.AddParameter(regParameter)

	err := parser.Parse("reg=-20200304")
	if err != nil {
		t.Error(err)
	}

	query, err := parser.ToBleveQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "+created_at:<=20200304"
	if query != expected {
		t.Errorf("Expected '%v' got '%v'", expected, query)
	}
}

func TestToBleveQueryDateRangeImplicitMax(t *testing.T) {
	parser := NewParser()

	regParameter := NewParameter("reg", DateRange)
	regParameter.OutputName = "created_at"
	regParameter.OutputCondition = Must
	parser.AddParameter(regParameter)

	err := parser.Parse("reg=20200101-")
	if err != nil {
		t.Error(err)
	}

	query, err := parser.ToBleveQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "+created_at:>=20200101"
	if query != expected {
		t.Errorf("Expected '%v' got '%v'", expected, query)
	}
}

func TestToBleveQueryShouldDisjunction(t *testing.T) {
	parser := NewParser()

	gender := NewParameter("gender", Strings)
	gender.AllowedValues = []string{"true", "false"}
	gender.OutputCondition = Must
	parser.AddParameter(gender)

	visibility := NewParameter("visibility", Strings)
	visibility.AllowedValues = []string{"pending", "suspended", "banned", "inactive"}
	visibility.OutputCondition = Should
	parser.AddParameter(visibility)

	err := parser.Parse("gender=true&visibility=pending,suspended")
	if err != nil {
		t.Error(err)
	}

	query, err := parser.ToBleveQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "+gender:true +(visibility:pending visibility:suspended)"
	if query != expected {
		t.Errorf("Expected '%v' got '%v'", expected, query)
	}
}

func TestToBleveQueryShouldSingleValue(t *testing.T) {
	parser := NewParser()

	gender := NewParameter("gender", Strings)
	gender.AllowedValues = []string{"true", "false"}
	gender.OutputCondition = Must
	parser.AddParameter(gender)

	visibility := NewParameter("visibility", Strings)
	visibility.AllowedValues = []string{"pending", "suspended", "banned", "inactive"}
	visibility.OutputCondition = Should
	parser.AddParameter(visibility)

	err := parser.Parse("gender=true&visibility=pending")
	if err != nil {
		t.Error(err)
	}

	query, err := parser.ToBleveQuery()
	if err != nil {
		t.Error(err)
	}

	expected := "+gender:true +(visibility:pending)"
	if query != expected {
		t.Errorf("Expected '%v' got '%v'", expected, query)
	}
}

func TestToBleveQuerySortSlice(t *testing.T) {

	parser := NewParser()

	sortParameter := NewParameter("sort", SortStrings)
	sortParameter.AllowedValues = []string{"age", "name", "last_online"}
	sortParameter.IncludeInOutput = false
	parser.AddParameter(sortParameter)

	queryString := "http://www.domain.com/search?sort=-age,name,-last_online"
	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 1 {
		t.Error("Invalid number of processed parameters")
	}

	expectedSortSlice := []string{"-age", "name", "-last_online"}
	sortSlice, err := parser.ToBleveSortSlice("sort")
	if err != nil {
		t.Error(err)
	}

	if !testEqString(expectedSortSlice, sortSlice) {
		t.Errorf("Expected '%v' got '%v'", expectedSortSlice, sortSlice)
	}

}

func TestToBleveQueryEmptyUnrecognizedSortSliceValues(t *testing.T) {

	parser := NewParser()

	sortParameter := NewParameter("sort", SortStrings)
	sortParameter.AllowedValues = []string{"age", "name", "last_online"}
	sortParameter.IncludeInOutput = false
	parser.AddParameter(sortParameter)

	queryString := "http://www.domain.com/search?sort=nothing,-that,it,-can,recognize"
	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 1 {
		t.Error("Invalid number of processed parameters")
	}

	// Get Bleve SortSlice (array of strings with '-' character as desc prefix, ex: []strings{"-age","created_at"})
	sortSlice, err := parser.ToBleveSortSlice("sort")
	if err != nil {
		t.Error(err)
	}

	// Expects a nil value if sort is empty -or- contains unrecognized values
	if sortSlice != nil {
		t.Error(err)
	}
}

func TestToBleveQueryNoSortParameter(t *testing.T) {

	parser := NewParser()

	sortParameter := NewParameter("sort", SortStrings)
	sortParameter.AllowedValues = []string{"age", "name", "last_online"}
	sortParameter.IncludeInOutput = false
	parser.AddParameter(sortParameter)

	activeParameter := NewParameter("active", Boolean)
	activeParameter.OutputCondition = Must
	parser.AddParameter(activeParameter)

	queryString := "http://www.domain.com/search?active=t"
	err := parser.Parse(queryString)
	if err != nil {
		t.Error(err)
	}

	if parser.ParsedParameterCount() != 1 {
		t.Error("Invalid number of processed parameters")
	}

	// Get Bleve SortSlice (array of strings with '-' character as desc prefix, ex: []strings{"-age","created_at"})
	sortSlice, err := parser.ToBleveSortSlice("sort")
	if err != nil {
		t.Error(err)
	}

	// Expects a nil value if sort is empty -or- contains unrecognized values
	if sortSlice != nil {
		t.Error(err)
	}
}
