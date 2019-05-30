## Quick-Start Delete

Follow the steps below to use SQM on your Delete commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```