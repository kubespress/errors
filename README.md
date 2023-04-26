# errors

Provides error enrichment utilities to store and retrieve additional context within errors as well as other general error handling utilities.

```
go get github.com/kubespress/errors
```

## Use cases

### Getting the error cause

The library mirrors the `errors.Is` and `errors.As` methods presented in the standard library, this allows for the original cause of an error to be determined and decisions to be taken based on this. All methods that mutate the errors in this package preserve the error chain, keeping this information.

```go
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("file missing")
    }
```

### Adding context to errors

For example to add a "user facing message" to errors you can use the `errors.Set` method:

```go
    return errors.Enrich(err,
        errors.Set[UserFacingMessage]("Invalid details entered")
    )
```

This context can be retrieved later:
```go
    fmt.Println(
        errors.Get[UserFacingMessage](err, "Internal error" /* default value*/)
    )
```

By using generics, types can be defined to hold specific context. Some examples include a `ErrorCode` or `UserFacingMessage` that can be presented to a user, while logging the original message. Context that is used internally can also be added, for example:

```go
    return errors.Enrich(err,
        errors.Set[IsTemporary](true)
    )
```

```go
    for {
        // Perform an action
        err := someFn()

        // If the error is temporary retry the action
        if errors.Check[IsTemporary](err) {
            continue
        }

        // Return the error (or nil)
        return err
    }
```

### Adding a call stack to an error

If you wish errors to include a call stack then this can be added using the `errors.Enrich` method. For example:

```go
    return errors.Enrich(err, errors.WithStack())
```

This call stack will be printed when the error is formatted with `%+v`.

### Filtering errors

There are some cases where depending on the error itself, you may want to ignore it. For example if you are writing a Kubernetes operator you often see this bit of code:

```go
    if err := r.Get(ctx, req.NamespacedName, &object); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }
```

This library allows the error to be enriched and for unwanted errors to be dropped in a single call, for example:

```go
    if err := r.Get(ctx, req.NamespacedName, &object); err != nil {
        return ctrl.Result{}, errors.Enrich(err, 
            client.IgnoreNotFound
            errors.Wrapf("failed to get %s", req.NamespacedName),
            errors.WithStack(),
        )
    }
```

### Aggregating errors

This allows multiple errors to be aggregated into a single error, while not breaking `errors.Is`, `errors.As`, ``errors.Check` and `errors.Get`, for example:

```go
func Example(fns []func()error) error {
    var errs errors.ErrorList

    for _, fn := range fns {
        if err := fn(); err != nil {
            errs = append(errs, err)
        }
    }

    return errs.Error()
}
```