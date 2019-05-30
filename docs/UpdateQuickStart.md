## Quick-Start Update

Follow the steps below to use SQM on your Update commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```
