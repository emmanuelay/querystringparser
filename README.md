# querystringparser

A package used to parse querystrings.
The package is designed to be used in a context where a querystring is passed as input to a search engine. 

In such a context the packages utility is to give you the possibility to set strict and granular rules for each parameter.

Example:
```
http://www.domain.com/search?q=alfa*&age=18-40&interests=alfa,beta,gamma&offset=10&size=15
```

MIT-licensed