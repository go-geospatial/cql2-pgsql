# cql2-pgsql

This is a port of the CQL2 transpiler developed for [pg_featureserv](https://github.com/CrunchyData/pg_featureserv)

## CQL Overview

CQL is a simple language to describe logical expressions. A CQL expression applies to values provided by feature properties and constants including numbers, booleans and text values. Values can be combined using the arithmetic operators +,-,*, / and % (modulo). Conditions on values are expressed using simple comparisons (<,>,<=,>=,=,<>). Other predicates include:

```sql
prop IN (val1, val2, ...)
prop BETWEEN val1 AND val2
prop IS [NOT] NULL
prop LIKE | ILIKE pattern
Conditions can be combined with the boolean operators AND,OR and NOT.
```

You will notice that this is very similar to SQL (probably not a coincidence!). That makes it straightforward to implement, and easy to use for us database people.

CQL also defines syntax for spatial and temporal filters.

## Filtering with CQL

A CQL expression can be used in a pg_featureserv request in the filter parameter. This is converted to SQL and included in the WHERE clause of the underlying database query. This allows the database to use its query planner and any defined indexes to execute the query efficiently.

Here's an example. We'll query the Natural Earth admin boundaries dataset, which we've loaded into PostGIS as a spatial table. (See this post for details on how to do this.) We're going to retrieve information about European countries where the population is 5,000,000 or less. The CQL expression for this is continent = 'Europe' AND pop_est <= 5000000.

Here's the query to get this result:

```sql
continent = "Europe" AND pop_est >= 5000000
```