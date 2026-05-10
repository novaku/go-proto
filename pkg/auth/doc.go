// Package auth provides JWT issuance/validation, gRPC authentication interceptors,
// and helpers to read claims from context. High-level code depends on TokenValidator
// for validation (Dependency Inversion), not on concrete JWT wiring details.
package auth
