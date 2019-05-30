<p align="center">
	<img src="https://sqm.dev/imgs/logo-sqm.png" height=350>
</center>


**In Development but with stable public API**


This lib is intended to be used as a one-to-one mapping of internal structure and a relational datasource, that is to say you're expected to use either raw sql or other tools for database administration queries.

While we do care about performance, our main focus is a readable and composable API. There will be some further study into techniques like memmoization and query compilation at startup but at the current state all the reflection and dynamic query building has no signinifcant overhead relative to the network operations a query does.

## Supported SQL Flavors
- MySQL
- Postgres
- Aurora
- SQLite


## Quick-start

Link with example Folder:

* [Select Quick-start](/docs/SelectQuickStart.md)
* [Delete Quick-start](/docs/DeleteQuickStart.md)
* [Insert Quick-start](/docs/InsertQuickStart.md)
* [Update Quick-start](/docs/UpdateQuickStart.md)
* [Count Quick-start](/docs/CountQuickStart.md)
* [Full Example](/docs/FullExample.md)