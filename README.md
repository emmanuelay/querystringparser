# querystringparser

A package used to parse querystrings.
The package is designed to be used in a context where a querystring is passed as input to a search engine.

In such a context the packages utility is to give you the possibility to set strict and granular rules for each parameter.

Example:
```
http://www.domain.com/search?q=alfa*&age=18-40&interests=alfa,beta,gamma&offset=10&size=15&active=true
```

## Parameter Types

| Type | Description | Example |
|------|-------------|---------|
| `Strings` | Delimited array of strings | `interests=alfa,beta,gamma` |
| `SearchString` | Wildcard pre/suffixed string for text search | `q=alfa*` |
| `SortStrings` | Delimited array with directional modifiers | `sort=name,-age,height` |
| `IntegerRange` | Integer range with hyphen separator | `age=18-40` |
| `Integer` | Single integer with min/max restrictions | `offset=10` |
| `Boolean` | String converted to boolean | `active=true` |
| `DateRange` | Date range with hyphen separator (YYYYMMDD) | `reg=20200101-20200304` |

### Boolean

Accepts the following values (case-insensitive): `true`, `t`, `false`, `f`.

### Strings with AllowedValues

The `AllowedValues` field on a parameter acts as a whitelist filter. When populated, only values present in the `AllowedValues` list are included in the parsed result.

### DateRange

Supports explicit and implicit date ranges using the `YYYYMMDD` format:
- Explicit: `reg=20200101-20200304` (both min and max)
- Implicit min: `reg=-20200304` (only max date set)
- Implicit max: `reg=20200101-` (only min date set)

If dates are swapped (min > max), they are automatically corrected. The `DateFormat` field on the parameter can be customized (defaults to `YYYYMMDD`).

### SortStrings

Supports directional modifiers where a `-` prefix indicates descending order. For example, `sort=name,-age` means sort by name ascending, then by age descending.

## Key Validation

Parameter keys are validated during parsing. Only lowercase alphanumeric characters, underscores, and dots are allowed. Parsing fails with `ErrInvalidKeyName` if a key contains unsanitized characters.

## Bleve Support

The package includes built-in support for generating [Bleve](https://github.com/blevesearch/bleve) search queries.

- `ToBleveQuery()` generates a Bleve query string with support for `Must` (+), `Not` (-), and `Should` conditions per parameter
- `ToBleveSortSlice()` converts `SortStrings` parameters into a Bleve-compatible sort slice

# TODO

- [x] Bleve support
	- [x] Add support for MUST(+), NOT(-) and SHOULD(..) for each parameter
	- [x] Ensure order of parameters in output corresponds to that of input
- [x] Implement support for date ranges
	- [x] Validate date (check format and content)
	- [x] Explicit range (?range=20200101-20200304)
	- [x] Implicit range, ends with (?range=-20200304)
	- [x] Implicit range, begins with (?range=20200101-)


MIT-licensed
