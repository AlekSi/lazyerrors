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
	"fmt"
	"runtime"
	"strings"
)

// Those variables are intentionally global and not synchronized.
// It is expected that they are set in the init() or main() function once.
var (
	// FileSegments specifies how many directories should be used for the file's path shortened form
	// in the Error() method.
	// Negative value uses the full path, 0 omits the whole path, 1 uses the file name only,
	// 2 uses one directory name and file name, etc.
	FileSegments = 1

	// FunctionSegments specifies how many import path segments should be used in the function name shortened form
	// in the Error() method.
	// Negative value uses the full path, 0 omits the whole function name, 1 uses the function name only,
	// 2 uses one import path segment and function name, etc.
	FunctionSegments = 1

	// Format specifies Error() method output format.
	// Passed arguments are:
	//   1. File's path shortened form (see FileSegments above).
	//   2. Line number.
	//   3. Function name shortened form (see FunctionSegments above).
	//   4. Original error.
	// Explicit argument indexes could be used; for example, to completely remove file and line, use "%[3]s: %[4]s".
	Format = "%s:%d (%s): %s"
)

// lazyerror adds a single program counter to the wrapped error.
//
// TODO https://github.com/AlekSi/lazyerrors/issues/1
type lazyerror struct {
	err error
	pc  uintptr
}

// loc returns file, line and function name for the stored program counter.
//
// Should it return original or shortened paths?
// TODO https://github.com/AlekSi/lazyerrors/issues/1
func (le lazyerror) loc() (file string, line int, function string) {
	if le.pc == 0 {
		return
	}

	frame, _ := runtime.CallersFrames([]uintptr{le.pc}).Next()
	return frame.File, frame.Line, frame.Function
}

// Error implements the [error] interface.
//
// It returns the wrapped error's message with location information.
func (le lazyerror) Error() string {
	file, line, function := le.loc()
	if file == "" && function == "" {
		return le.err.Error()
	}

	return fmt.Sprintf(
		Format,
		shorten(file, FileSegments),
		line,
		shorten(function, FunctionSegments),
		le.err.Error(),
	)
}

// GoString implements the [fmt.GoStringer] interface.
//
// It exists so `%#v` fmt verb could correctly print wrapped errors.
func (le lazyerror) GoString() string {
	return fmt.Sprintf("lazyerror{%q}", le.Error())
}

// Unwrap returns the wrapped error.
func (le lazyerror) Unwrap() error {
	return le.err
}

// shorten returns the shortened form of the given path.
//
// TODO https://github.com/AlekSi/lazyerrors/issues/1
func shorten(path string, segments int) string {
	switch {
	case segments == 0:
		return ""

	case segments == 1:
		if i := strings.LastIndexByte(path, '/'); i != -1 {
			return path[i+1:]
		}

	case segments > 1:
		p := strings.SplitAfter(path, "/")
		if i := len(p) - segments; i > 0 {
			return strings.Join(p[i:], "")
		}
	}

	return path
}

// check interfaces
var (
	_ error                       = &lazyerror{}
	_ fmt.GoStringer              = &lazyerror{}
	_ interface{ Unwrap() error } = &lazyerror{}

	// Should the receiver be a value?
	// TODO https://github.com/AlekSi/lazyerrors/issues/1
	_ error                       = lazyerror{}
	_ fmt.GoStringer              = lazyerror{}
	_ interface{ Unwrap() error } = lazyerror{}
)
