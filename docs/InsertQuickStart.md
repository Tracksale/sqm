## Quick-Start Insert

Follow the steps below to use SQM on your Insert commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```

### 2. Insert new Object on DB

Just pass a struct that represents table entity. See the example below:

```go
    //Represents: INSERT INTO table (fields)
    //              VALUES (...)
    result, err := sqm.Using(db, "table").Insert(newObject)
```