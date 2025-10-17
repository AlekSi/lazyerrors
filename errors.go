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

import "errors"

// ErrUnsupported is the same as [errors.ErrUnsupported].
// It is provided for convenience.
var ErrUnsupported = errors.ErrUnsupported

// As is the same as [errors.As].
// It is provided for convenience.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is is the same as [errors.Is].
// It is provided for convenience.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Unwrap is the same as [errors.Unwrap].
// It is provided for convenience.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
