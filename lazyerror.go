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
type lazyerror struct {
	err error
	pc  uintptr
}

// Error implements [error] interface.
func (le lazyerror) Error() string {
	if le.pc == 0 {
		return le.err.Error()
	}

	frame, _ := runtime.CallersFrames([]uintptr{le.pc}).Next()
	if frame.File == "" || frame.Function == "" {
		return le.err.Error()
	}

	return fmt.Sprintf(
		Format,
		shortPath(frame.File, FileSegments), frame.Line,
		shortPath(frame.Function, FunctionSegments),
		le.err.Error(),
	)
}

// GoString implements [fmt.GoStringer] interface.
//
// It exists so %#v fmt verb could correctly print wrapped errors.
func (le lazyerror) GoString() string {
	return fmt.Sprintf("lazyerror{%q}", le.Error())
}

// Unwrap returns the wrapped error.
func (le lazyerror) Unwrap() error {
	return le.err
}

// shortPath returns shorter path for the given path.
func shortPath(path string, parts int) string {
	switch {
	case parts == 0:
		return ""

	case parts == 1:
		if i := strings.LastIndexByte(path, '/'); i != -1 {
			return path[i+1:]
		}

	case parts > 1:
		p := strings.SplitAfter(path, "/")
		if i := len(p) - parts; i > 0 {
			return strings.Join(p[i:], "")
		}
	}

	return path
}

// check interfaces
var (
	_ error                       = lazyerror{}
	_ fmt.GoStringer              = lazyerror{}
	_ interface{ Unwrap() error } = lazyerror{}
)
