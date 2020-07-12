# querystringparser

A package used to parse querystrings.
The package is designed to be used in a context where a querystring is passed as input to a search engine. 

In such a context the packages utility is to give you the possibility to set strict and granular rules for each parameter.

Example:
```
http://www.domain.com/search?q=alfa*&age=18-40&interests=alfa,beta,gamma&offset=10&size=15
```

# TODO

- [ ] Bleve support
	- [ ] Add support for MUST(+), NOT(-) and SHOULD(..) for each parameter
	- [ ] Ensure order of parameters in output corresponds to that of input
- [ ] Implement support for date ranges 
	- [ ] Validate date (check format and content)
	- [ ] Explicit range (?range=20200101-20200304)
	- [ ] Implicit range, ends with (?range=-20200304)
	- [ ] Implicit range, begins with (?range=20200101-)


MIT-licensed