// Copyright 2021 FerretDB Inc.
// Copyright 2025 Alexey Palazhchenko.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lazyerrors

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// assertEqual fails the test if expected and actual are not equal.
func assertEqual[T any](t testing.TB, expected, actual T) {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		return
	}

	t.Errorf("Not equal, but should be:\nexpected: %#v\nactual  : %#v", expected, actual)
}

// assertNotEqual fails the test if expected and actual are equal.
func assertNotEqual[T any](t testing.TB, expected, actual T) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		return
	}

	t.Errorf("Equal, but should not be:\nexpected: %#v\nactual  : %#v", expected, actual)
}

// unwrap [errors.Unwrap] err n times.
func unwrap(err error, n int) error {
	for range n {
		err = errors.Unwrap(err)
	}
	return err
}

func TestStdErrors(t *testing.T) {
	t.Parallel()

	err := errors.New("err")
	err1 := fmt.Errorf("err1: %w", err)
	err2 := fmt.Errorf("err2: %w", err1)
	err3 := fmt.Errorf("err3: %w", err2)

	assertEqual(t, "err", err.Error())
	assertEqual(t, "err1: err", err1.Error())
	assertEqual(t, "err2: err1: err", err2.Error())
	assertEqual(t, "err3: err2: err1: err", err3.Error())

	assertEqual(t, `&errors.errorString{s:"err"}`, fmt.Sprintf("%#v", err))
	assertEqual(t, true, strings.Contains(fmt.Sprintf("%#v", err1), `&fmt.wrapError{msg:"err1: err", err:(*errors.errorString)(0x`))

	assertEqual(t, err, unwrap(err1, 1))
	assertEqual(t, nil, unwrap(err1, 2))

	assertEqual(t, err1, unwrap(err2, 1))
	assertEqual(t, err, unwrap(err2, 2))
	assertEqual(t, nil, unwrap(err2, 3))

	assertEqual(t, err2, unwrap(err3, 1))
	assertEqual(t, err1, unwrap(err3, 2))
	assertEqual(t, err, unwrap(err3, 3))
	assertEqual(t, nil, unwrap(err3, 4))

	assertEqual(t, true, errors.Is(err3, err3))
	assertEqual(t, true, errors.Is(err3, err2))
	assertEqual(t, true, errors.Is(err3, err1))
	assertEqual(t, true, errors.Is(err3, err))
}

func TestErrors(t *testing.T) {
	t.Parallel()

	err := New("err")
	err1 := Errorf("err1: %w", err)
	err2 := Errorf("err2: %w", err1)
	err3 := Errorf("err3: %w", err2)

	expected := "lazyerrors_test.go:95 (lazyerrors.TestErrors): err"
	assertEqual(t, expected, err.Error())
	expected = "lazyerrors_test.go:96 (lazyerrors.TestErrors): err1: " +
		"lazyerrors_test.go:95 (lazyerrors.TestErrors): err"
	assertEqual(t, expected, err1.Error())
	expected = "lazyerrors_test.go:97 (lazyerrors.TestErrors): err2: " +
		"lazyerrors_test.go:96 (lazyerrors.TestErrors): err1: " +
		"lazyerrors_test.go:95 (lazyerrors.TestErrors): err"
	assertEqual(t, expected, err2.Error())
	expected = "lazyerrors_test.go:98 (lazyerrors.TestErrors): err3: " +
		"lazyerrors_test.go:97 (lazyerrors.TestErrors): err2: " +
		"lazyerrors_test.go:96 (lazyerrors.TestErrors): err1: " +
		"lazyerrors_test.go:95 (lazyerrors.TestErrors): err"
	assertEqual(t, expected, err3.Error())

	expected = `lazyerror{"lazyerrors_test.go:95 (lazyerrors.TestErrors): err"}`
	assertEqual(t, expected, fmt.Sprintf("%#v", err))
	expected = `lazyerror{"lazyerrors_test.go:96 (lazyerrors.TestErrors): err1: ` +
		`lazyerrors_test.go:95 (lazyerrors.TestErrors): err"}`
	assertEqual(t, expected, fmt.Sprintf("%#v", err1))
	expected = `lazyerror{"lazyerrors_test.go:97 (lazyerrors.TestErrors): err2: ` +
		`lazyerrors_test.go:96 (lazyerrors.TestErrors): err1: ` +
		`lazyerrors_test.go:95 (lazyerrors.TestErrors): err"}`
	assertEqual(t, expected, fmt.Sprintf("%#v", err2))
	expected = `lazyerror{"lazyerrors_test.go:98 (lazyerrors.TestErrors): err3: ` +
		`lazyerrors_test.go:97 (lazyerrors.TestErrors): err2: ` +
		`lazyerrors_test.go:96 (lazyerrors.TestErrors): err1: ` +
		`lazyerrors_test.go:95 (lazyerrors.TestErrors): err"}`
	assertEqual(t, expected, fmt.Sprintf("%#v", err3))

	assertNotEqual(t, err, unwrap(err1, 1))
	assertEqual(t, err, unwrap(err1, 2))
	assertNotEqual(t, nil, unwrap(err1, 3))
	assertEqual(t, nil, unwrap(err1, 4))

	assertNotEqual(t, err1, unwrap(err2, 1))
	assertEqual(t, err1, unwrap(err2, 2))
	assertNotEqual(t, err, unwrap(err2, 3))
	assertEqual(t, err, unwrap(err2, 4))
	assertNotEqual(t, nil, unwrap(err2, 5))
	assertEqual(t, nil, unwrap(err2, 6))

	assertNotEqual(t, err2, unwrap(err3, 1))
	assertEqual(t, err2, unwrap(err3, 2))
	assertNotEqual(t, err1, unwrap(err3, 3))
	assertEqual(t, err1, unwrap(err3, 4))
	assertNotEqual(t, err, unwrap(err3, 5))
	assertEqual(t, err, unwrap(err3, 6))
	assertNotEqual(t, nil, unwrap(err3, 7))
	assertEqual(t, nil, unwrap(err3, 8))

	assertEqual(t, true, errors.Is(err3, err3))
	assertEqual(t, true, errors.Is(err3, err2))
	assertEqual(t, true, errors.Is(err3, err1))
	assertEqual(t, true, errors.Is(err3, err))
}

func TestPC(t *testing.T) {
	t.Parallel()

	runtime.LockOSThread()

	ch := make(chan error, 1)

	go func() {
		runtime.LockOSThread()
		ch <- New("err")
	}()

	err := <-ch
	runtime.Gosched()
	assertEqual(t, "lazyerrors_test.go:166 (lazyerrors.TestPC.func1): err", err.Error())
}

// errPackage is a package-level error to test init function call location.
var errPackage = New("err package")

func TestFormat(t *testing.T) {
	var fis, fns int
	var f string

	fis, FileSegments = FileSegments, 0
	fns, FunctionSegments = FunctionSegments, -1
	f, Format = Format, "%[3]s: %[4]s"

	t.Cleanup(func() {
		FileSegments = fis
		FunctionSegments = fns
		Format = f
	})

	assertEqual(t, "github.com/AlekSi/lazyerrors.init: err package", errPackage.Error())
}

func TestShortPath(t *testing.T) {
	t.Parallel()

	assertEqual(t, "/absolute/path/lazyerrors.go", shortPath("/absolute/path/lazyerrors.go", -1))
	assertEqual(t, "", shortPath("/absolute/path/lazyerrors.go", 0))
	assertEqual(t, "lazyerrors.go", shortPath("/absolute/path/lazyerrors.go", 1))
	assertEqual(t, "path/lazyerrors.go", shortPath("/absolute/path/lazyerrors.go", 2))
	assertEqual(t, "absolute/path/lazyerrors.go", shortPath("/absolute/path/lazyerrors.go", 3))
	assertEqual(t, "/absolute/path/lazyerrors.go", shortPath("/absolute/path/lazyerrors.go", 4))
	assertEqual(t, "/absolute/path/lazyerrors.go", shortPath("/absolute/path/lazyerrors.go", 5))

	assertEqual(t, "/lazyerrors.go", shortPath("/lazyerrors.go", -1))
	assertEqual(t, "", shortPath("/lazyerrors.go", 0))
	assertEqual(t, "lazyerrors.go", shortPath("/lazyerrors.go", 1))
	assertEqual(t, "/lazyerrors.go", shortPath("/lazyerrors.go", 2))

	assertEqual(t, "relative/path/lazyerrors.go", shortPath("relative/path/lazyerrors.go", -1))
	assertEqual(t, "", shortPath("relative/path/lazyerrors.go", 0))
	assertEqual(t, "lazyerrors.go", shortPath("relative/path/lazyerrors.go", 1))
	assertEqual(t, "path/lazyerrors.go", shortPath("relative/path/lazyerrors.go", 2))
	assertEqual(t, "relative/path/lazyerrors.go", shortPath("relative/path/lazyerrors.go", 3))
	assertEqual(t, "relative/path/lazyerrors.go", shortPath("relative/path/lazyerrors.go", 4))
	assertEqual(t, "relative/path/lazyerrors.go", shortPath("relative/path/lazyerrors.go", 5))

	assertEqual(t, "lazyerrors.go", shortPath("lazyerrors.go", -1))
	assertEqual(t, "", shortPath("lazyerrors.go", 0))
	assertEqual(t, "lazyerrors.go", shortPath("lazyerrors.go", 1))
	assertEqual(t, "lazyerrors.go", shortPath("lazyerrors.go", 2))
}

var drain error

func BenchmarkStdNew(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		drain = errors.New("lazyerrors_test.go:179 (lazyerrors.BenchmarkStdNew): err")
	}

	b.StopTimer()

	assertNotEqual(b, nil, drain)
	assertEqual(b, "lazyerrors_test.go:179 (lazyerrors.BenchmarkStdNew): err", drain.Error())
}

func BenchmarkNStdNew(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		drain = errors.New("lazyerrors_test.go:192 (lazyerrors.BenchmarkNStdNew): err")
	}

	b.StopTimer()

	assertNotEqual(b, nil, drain)
	assertEqual(b, "lazyerrors_test.go:192 (lazyerrors.BenchmarkNStdNew): err", drain.Error())
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		drain = New("err")
	}

	b.StopTimer()

	assertNotEqual(b, nil, drain)
	assertEqual(b, "lazyerrors_test.go:235 (lazyerrors.BenchmarkNew): err", drain.Error())
}

func BenchmarkNNew(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		drain = New("err")
	}

	b.StopTimer()

	assertNotEqual(b, nil, drain)
	assertEqual(b, "lazyerrors_test.go:248 (lazyerrors.BenchmarkNNew): err", drain.Error())
}

func Example() {
	func1 := func() error {
		// provide additional context
		return Errorf("i'm not lazy: %w", io.EOF)
	}

	func2 := func() error {
		// or don't
		return Error(func1())
	}

	err := func2()

	FunctionSegments = 0
	Format = "%[1]s:%[2]d: %[4]s"

	fmt.Println(err)
	fmt.Println(errors.Is(err, io.EOF))

	// Output:
	// lazyerrors_test.go:286: lazyerrors_test.go:281: i'm not lazy: EOF
	// true
}
