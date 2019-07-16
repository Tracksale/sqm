## Quick-Start Update

Follow the steps below to use SQM on your Update commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```

### 2. Add Where conditions

'Where' function can be built using groups and conditions, see [select topic](SelectQuickStart.md)

### 3. Call Update Function

Update function must be called using a struct that represents the new form of object on DB

```go
    //Represents: UPDATE table
    //              SET ...
    //              WHERE uuid='my_uuid'
    query.Where(
		sqm.C("uuid", sqm.Equal, uuid),
	).Update(object)
```