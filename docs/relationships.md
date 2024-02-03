# Relationships

You may set up eager loaded relationships (`Preloads`).
You can define them from resources or at runtime.

## Preloads resource

```go
r := &ClientResource{
    tables.AbstractResource{
        DB:              db,
        Request:         req,
        HasGlobalSearch: true,
        Preloads:        []*tables.Preload{{Name: "Owner"}},
    },
}
```

## Preloads runtime 

```go
resource.Preloads = []*tables.Preload{{
    Name: "Owner",
    Extra: func(db *gorm.DB) *gorm.DB {
        return db.Select("id", "email")
    }},
}
```