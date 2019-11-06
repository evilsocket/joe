# Joe

**WORK IN PROGRESS, DO NOT USE IN PRODUCTION**

Scenario: your company happens to have big amounts of historical data in one or more relational databases and a
 relatively complex backend wrapping the business logic around this data. But it doesn't matter how complex this
  backend gets, your users will always find those few corner cases that:

* Are not worth implementing as a whole new feature for the backend, being extremely specific to that instance / request.
* Are still something that must be done every few days or weeks.

So you either implement stored procedures and views from the database side of things, but then you need to give access 
to the database to people that are not DBA and don't know SQL.

Or you duplicate all of your relational data to an Elastic search or whatever other stack that allows you to use tools 
such as Kibana. But this means not only duplicating data, but also integrating a bunch of new technologies just to have 
basic charts, csv files and answer your C*O's question of the week.

Or you keep .sql/.txt files all over your computer with all the queries your coworkers ask you periodically to run.

So I thought about Joe.

![joe](https://i.imgur.com/EtzqrVe.png)

Joe is a very simple microservice written in Go that allows you to:

* Version and annotate SQL queries (for several types of relational databases all handled transparently by the service).
* Those queries, used as statements, therefore parametrized, will be wrapped in an automatically generated REST API that
 can be just called from CURL (with JWT or http based authentication).
* The output, processed by the meta engine (that will analyze the format of the output and apply transformation and caching policies), will be available as simple, paginated and normalized JSON output.
* This JSON output, that can be already used directly from a researcher or whatever, will be optionally connected to the presentation engine, that will expose the data in different ways than JSON, such as charts created automatically if temporal data is detected, searchable datagrids with paging, simple transformations like JSON -> CSV and so on.

All of the above and more, by just storing the query ... instead of using those txt files you keep around your
 computer :D It's basically a weird mix of a smart phpmyadmin, stored procedures, git (for annotation and versioning), 
 kibana (for charting) and  ... things :D
 
## License

`joe` is made with ♥  by [@evilsocket](https://twitter.com/evilsocket) and it is released under the GPL3 license.