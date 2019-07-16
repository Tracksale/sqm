## Quick-Start Count

Follow the steps below to use SQM on your Count commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```

### 2. Add Some Conditions
You can add conditions as you want, for more details see [select session](SelectQuickStart.md)


### 3. Call Count Function
Count function expects an `int` param where the result will be store

```go
    //Represents: SELECT COUNT(*)
    //              FROM table
    var counter int
    err := sqmQuery.Count(&counter)
```