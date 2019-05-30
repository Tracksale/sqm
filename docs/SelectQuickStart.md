## Quick-Start Select

Follow the steps below to use SQM on your select commands.

### 1. Get a Query instance

```go
    // db its a native *sql.DB instance
    query := sqm.Using(db, "{{table_name}}")
```

### 2. Add Where conditions

Where func accepts `sqm.C` or `sqm.Group` commands, as much as you want.

#### sqm.C


```go
    //Represents: WHERE c1 = "foo" AND c2 = "bar"
    groupCondition := sqm.Group(
        sqm.C("c1", sqm.Equal, "foo")
        sqm.C("c2", sqm.Equal, "bar")
    )
```

#### sqm.Group

Creates a Group Condition. Basically adds a parenthesis on Conditions. Examples:

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

Some conditions require more than one value. Example:

```go
    //Represents: IN ('A', 'B')
    var interfaceValues []interface{}
    interfaceValues = append(interfaceValues, "A")
    interfaceValues = append(interfaceValues, "B")

    sqm.C(field, sqm.In, interfaceValues...)
```

If you pass fewer parameters that are required for this query an error will occur.

All supported conditions to use:

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


### 3. Opitinional Values

#### OrderBy

Add OrderBy operation on this query. Example:

```go
    //Represents: ORDER BY field DESC
    query.OrderBy(field, sqm.Desc)
```

```go
    //Represents: ORDER BY field ASC
    query.OrderBy(field, sqm.Asc)
```


#### Limit

Add Limit operation on this query, the first value of this function is `limit`, the second is `offset`. Example:

```go
    //Represents: LIMIT 2 OFFSET 1
    query.Limit(2, 1)
```

### 4. Select Command

Run produced query and stores the result on the variable passed on args. 
If the passed variable on args doesn't support the result of query an error will occur.


The columns that will be selected are those on struct passed on Select func

**SQM ONLY SUPPORTS STRUCTS ON SELECT**



Examples:

```go
    //Represents: SELECT ... FROM customers 
    //                WHERE deleted_at IS NULL AND
    //                  (phone = 'phone_value' OR email = 'email_value' )
	phoneOrEmailCond := sqm.Group(
		sqm.C("phone", sqm.Equal, phone),
		sqm.Or,
		sqm.C("email", sqm.Equal, email),
	)

	customer := Customer{}
	err = sqm.Using(db, "customers").
		Where(
			sqm.C("deleted_at", sqm.IsNull),
			phoneOrEmailCond,
		).Limit(1, 0).Select(&customer)

```

