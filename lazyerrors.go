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

// Package lazyerrors provides error wrapping with location information:
// file path, line number, and function/method name.
//
// [New], [Errorf], [Error], [Maybe]/[Maybe2]/[Maybe3], and [Join] functions create a new error
// with location captured as a single uintptr for Program Counter (PC).
//
// Only one location is captured for each error value, not a full call stack.
// If the "return stack" is needed, use the functions mentioned above
// with each return statement, channel operations, etc.
// This return stack is superior to the regular call stack,
// as it correctly shows line numbers of error propagation.
//
// Actual error formatting happens lazily in the `Error() string` method,
// and can be changed by setting [FileSegments], [FunctionSegments], and [Format] variables.
package lazyerrors

import (
	"errors"
	"fmt"
	"runtime"
)

// New returns an error created with [errors.New] wrapped with a single location.
func New(s string) error {
	return lazyerror{
		err: errors.New(s),
		pc:  pc(),
	}
}

// Errorf returns an error created with [fmt.Errorf] wrapped with a single location.
func Errorf(format string, a ...any) error {
	return lazyerror{
		err: fmt.Errorf(format, a...),
		pc:  pc(),
	}
}

// Error returns an error wrapped with a single location.
// It panics if err is nil; the caller is expected to check if err != nil before using this function.
// Alternatively, use [Maybe] instead.
func Error(err error) error {
	if err == nil {
		panic("err is nil")
	}

	return lazyerror{
		err: err,
		pc:  pc(),
	}
}

// Maybe returns an error wrapped with a single location,
// or nil, if err is nil.
func Maybe(err error) error {
	if err == nil {
		return nil
	}

	return lazyerror{
		err: err,
		pc:  pc(),
	}
}

// Maybe2 returns an error wrapped with a single location, along with v, if err is not nil.
// Otherwise, it returns (v, nil).
func Maybe2[T any](v T, err error) (T, error) {
	if err == nil {
		return v, nil
	}

	return v, lazyerror{
		err: err,
		pc:  pc(),
	}
}

// Maybe3 returns an error wrapped in a single location, along with v1 and v2, if err is not nil.
// Otherwise, it returns (v1, v2, nil).
func Maybe3[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2, error) {
	if err == nil {
		return v1, v2, nil
	}

	return v1, v2, lazyerror{
		err: err,
		pc:  pc(),
	}
}

// Join returns an error created with [errors.Join] wrapped with a single location.
//
// Any nil error values are discarded, and nil is returned if no values are left.
// But unlike [errors.Join], a non-nil error returned implements the `Unwrap() error` method,
// not `Unwrap() []error`.
func Join(errs ...error) error {
	err := errors.Join(errs...)
	if err == nil {
		return nil
	}

	return lazyerror{
		err: err,
		pc:  pc(),
	}
}

// pc returns a program counter of the caller's caller.
func pc() uintptr {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)

	return pc[0]
}
