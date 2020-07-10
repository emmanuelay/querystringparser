package querystringparser

import (
	"testing"
)

func TestStrings(t *testing.T) {
	stringsParameter := NewParameter("interest")
	stringsParameter.Type = Strings
	err := stringsParameter.Parse("interest", "alfa,beta,gamma,delta")
	if err != nil {
		t.Error(err)
	}

	if testEqString(stringsParameter.StringsValue, []string{"alfa", "beta", "gamma", "delta"}) != true {
		t.Error("Invalid StringsValue")
	}
}

func TestEmptyStringsParameter(t *testing.T) {

	interestParameter := NewParameter("interest")
	interestParameter.Type = Strings
	err := interestParameter.Parse("interest", "")
	if err != nil {
		t.Error(err)
	}

	expectedStrings := []string{}
	if testEqString(interestParameter.StringsValue, expectedStrings) != true {
		t.Error("Invalid StringsValue")
	}
}
func TestIntegerRange(t *testing.T) {
	intRangeParameter := NewParameter("age")
	intRangeParameter.Type = IntegerRange
	err := intRangeParameter.Parse("age", "118-35")
	if err != nil {
		t.Error(err)
	}

	if intRangeParameter.MinValue != 35 {
		t.Error("Invalid MinValue")
	}

	if intRangeParameter.MaxValue != 118 {
		t.Error("Invalid MaxValue")
	}
}

func TestIntegerRangeMinDefault(t *testing.T) {
	intRangeParameter := NewParameter("age")
	intRangeParameter.Type = IntegerRange
	intRangeParameter.MinValue = 10
	err := intRangeParameter.Parse("age", "-35")
	if err != nil {
		t.Error(err)
	}

	if intRangeParameter.MinValue != 10 {
		t.Error("Invalid MinValue")
	}

	if intRangeParameter.MaxValue != 35 {
		t.Error("Invalid MaxValue")
	}
}

func TestIntegerRangeMaxDefault(t *testing.T) {
	intRangeParameter := NewParameter("age")
	intRangeParameter.Type = IntegerRange
	intRangeParameter.MaxValue = 80
	err := intRangeParameter.Parse("age", "18-")
	if err != nil {
		t.Error(err)
	}

	if intRangeParameter.MinValue != 18 {
		t.Error("Invalid MinValue")
	}

	if intRangeParameter.MaxValue != 80 {
		t.Error("Invalid MaxValue")
	}
}

func TestInteger(t *testing.T) {
	integerParameter := NewParameter("count")
	integerParameter.Type = Integer
	err := integerParameter.Parse("count", "18")
	if err != nil {
		t.Error(err)
	}

	if integerParameter.IntValue != 18 {
		t.Error("Invalid IntValue")
	}
}

func TestIntegerRoof(t *testing.T) {
	integerParameter := NewParameter("offset")
	integerParameter.Type = Integer
	integerParameter.MaxValue = 100
	err := integerParameter.Parse("offset", "118")
	if err != nil {
		t.Error(err)
	}

	if integerParameter.IntValue != 100 {
		t.Error("Invalid IntValue")
	}
}

func TestIntegerFloor(t *testing.T) {
	integerParameter := NewParameter("offset")
	integerParameter.Type = Integer
	integerParameter.MaxValue = 100
	integerParameter.MinValue = 10
	err := integerParameter.Parse("offset", "8")
	if err != nil {
		t.Error(err)
	}

	if integerParameter.IntValue != 10 {
		t.Error("Invalid IntValue")
	}
}

func TestSearchStringSuffix(t *testing.T) {
	searchStringParameter := NewParameter("q")
	searchStringParameter.Type = SearchString
	err := searchStringParameter.Parse("q", "alfa*")
	if err != nil {
		t.Fail()
	}

	if searchStringParameter.StringValue != "alfa" {
		t.Error("Invalid string")
	}

	if searchStringParameter.Position != Suffix {
		t.Error("Invalid position")
	}
}

func TestSearchStringPrefix(t *testing.T) {
	searchStringParameter := NewParameter("q")
	searchStringParameter.Type = SearchString
	err := searchStringParameter.Parse("q", "*beta")
	if err != nil {
		t.Fail()
	}

	if searchStringParameter.StringValue != "beta" {
		t.Error("Invalid string")
	}

	if searchStringParameter.Position != Prefix {
		t.Error("Invalid position")
	}
}

func TestSearchStringSurrounded(t *testing.T) {
	searchStringParameter := NewParameter("q")
	searchStringParameter.Type = SearchString
	err := searchStringParameter.Parse("q", "*gamma*")
	if err != nil {
		t.Fail()
	}

	if searchStringParameter.StringValue != "gamma" {
		t.Error("Invalid string")
	}

	if searchStringParameter.Position != Surrounded {
		t.Error("Invalid position")
	}
}

func TestSearchStringSurroundedMaxLength(t *testing.T) {
	searchStringParameter := NewParameter("q")
	searchStringParameter.Type = SearchString
	searchStringParameter.MaxLength = 3
	err := searchStringParameter.Parse("q", "*gamma*")
	if err.Error() != "Invalid length (5) for parameter 'q' (max 3)" {
		t.Error("Expected error")
	}

	if searchStringParameter.StringValue != "gam" {
		t.Error("Invalid resulting string")
	}

	if searchStringParameter.Position != Surrounded {
		t.Error("Invalid position")
	}
}

func TestSearchStringSuffixedMinLength(t *testing.T) {
	searchStringParameter := NewParameter("q")
	searchStringParameter.Type = SearchString
	searchStringParameter.MinLength = 3
	err := searchStringParameter.Parse("q", "ga*")
	if err.Error() != "Invalid length (2) for parameter 'q' (min 3)" {
		t.Error("Expected error")
	}

	if searchStringParameter.StringValue != "ga" {
		t.Error("Invalid resulting string", searchStringParameter.StringValue)
	}

	if searchStringParameter.Position != Suffix {
		t.Error("Invalid position")
	}
}

func TestSortModifierParameter(t *testing.T) {
	sortModifierParameter := NewParameter("sort")
	sortModifierParameter.Type = SortStrings
	err := sortModifierParameter.Parse("sort", "name,-age,height")
	if err != nil {
		t.Error(err)
	}

	if testEqString(sortModifierParameter.StringsValue, []string{"name", "age", "height"}) != true {
		t.Error("Invalid StringsValue")
	}

	if testEqBool(sortModifierParameter.SortDirections, []bool{true, false, true}) != true {
		t.Error("Invalid SortDirections")
	}

}

// https://stackoverflow.com/a/15312097/254695
func testEqString(a, b []string) bool {

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

// https://stackoverflow.com/a/15312097/254695
func testEqBool(a, b []bool) bool {

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
