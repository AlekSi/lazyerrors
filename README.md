# lazyerrors

[![Go Reference](https://pkg.go.dev/badge/github.com/AlekSi/lazyerrors.svg)](https://pkg.go.dev/github.com/AlekSi/lazyerrors)

You know that as a good Go developer, you have to annotate your errors with the necessary context:

```go
func (mux *ServeMux) registerErr(patstr string, handler Handler) error {
	// [...]
	if err != nil {
		return fmt.Errorf("parsing %q: %w", patstr, err)
	}
	// [...]
}
```

But sometimes you can't think of anything other than the function/method name,
or you're just (pragmatically) lazy:

```go
func (r *Reader) nextPart(rawPart bool, maxMIMEHeaderSize, maxMIMEHeaders int64) (*Part, error) {
	// [...]
	if err != nil {
		return nil, fmt.Errorf("multipart: NextPart: %w", err)
	}
	// [...]
}
```

Package lazyerrors provides error wrapping with location information
(file path, line number, and function/method name)
for when you are lazy:

```sh
go get github.com/AlekSi/lazyerrors@latest
```

```go
import "github.com/AlekSi/lazyerrors"

func Example() {
	func1 := func() error {
		// provide additional context
		return lazyerrors.Errorf("i'm not lazy: %w", io.EOF)
	}

	func2 := func() error {
		// or don't
		return lazyerrors.Error(func1())
	}

	err := func2()

	fmt.Println(err)
	fmt.Println(errors.Is(err, io.EOF))
}
```

The code above produces:

```
lazyerrors_test.go:288 (Example.func2): lazyerrors_test.go:283 (Example.func1): i'm not lazy: EOF
true
```

`New`, `Error`, `Errorf`, and `Join` functions create a new error
with location captured as a single uintptr for Program Counter (PC).

Only one location is captured for each error value, not a full call stack.
If the "return stack" is needed, use the functions mentioned above
with each return statement, channel operations, etc.
This return stack is superior to the regular call stack,
as it correctly shows line numbers of error propagation.

Actual error formatting happens lazily in the `Error() string` method,
and can be changed by setting `FileSegments`, `FunctionSegments`, and `Format` variables.

See [package documentation](https://pkg.go.dev/github.com/AlekSi/lazyerrors) for more details.
