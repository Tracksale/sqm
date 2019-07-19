## Quick-Start Delete

Follow the steps below to use SQM on your Delete commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```

### 2. Add Some Conditions
You can add conditions as you want, for more details see [select session](SelectQuickStart.md)

### 3. Call Delete Function
Delete function exec a DELETE query. See example below:

```go
    //Represents: DELETE FROM table
    //              WHERE uuid='my_uuid'
    query.Where(
		sqm.C("uuid", sqm.Equal, uuid),
	).Delete()
```