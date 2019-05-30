## Quick-Start Select
Follow the steps below to use SQM on your select commands.

### 1. Get a Query instance
```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```

### 2. Add Where conditions
Where func accepts `sqm.C` or `sqm.Group` commands, as much as you want.

### sqm.C


```go
    //Represents: WHERE c1 = "foo" AND c2 = "bar"
    groupCondition := sqm.Group(
        sqm.C("c1", sqm.Equal, "foo")
        sqm.C("c2", sqm.Equal, "bar")
    )
```

#### sqm.Group
Creates a Group Condition. Basicaly adds a parenthesis on Conditions. Examples:

```go
    //Represents: ( c1 = "foo" AND c2 = "bar" )
    groupCondition1 := sqm.Group(
        sqm.C("c1", sqm.Equal, "foo"),
        sqm.C("c2", sqm.Equal, "bar"),
    )
```


```go
    //Represents: ( c1 = "foo" OR c2 = "bar" )
    groupCondition2 := sqm.Group(
        sqm.C("c1", sqm.Equal, "foo"),
        sqm.Or,
        sqm.C("c2", sqm.Equal, "bar"),
    )
```

```go
    //Represents: (( c1 = "foo" OR c2 = "bar" ) AND ( c1 = "foo" AND c2 = "bar" ))
    groupCondition2 := sqm.Group(
        groupCondition1,
        groupCondition2,
    )
```


#### Supported Conditions

| SQM  |  SQL |
|---|---|
|  sqm.Equal |  = |
|  sqm.NotEqual |  != |
|  sqm.Like | LIKE  |
|  sqm.NotLike | NOT LIKE  |
|  sqm.In |  IN |
|  sqm.NotIn |  NOT IN |
|  sqm.IsNull |  IS NULL |
|  sqm.IsNotNull | IS NOT NULL  |
|  sqm.Between |  BETWEEN |
|  sqm.Greater |  > |
|  sqm.Lower |  < |
|  sqm.LowerEqual |  <= |
|  sqm.GreaterEqual | >=  |


