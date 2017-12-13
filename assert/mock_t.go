package assert

import (
	"fmt"
	"strings"
	"testing"
)

// MatchType is used in MockT's CheckPredicates to control whether
// recorded testing.TB operations match predicates.
type MatchType int

const (
	// Success indicates a predicate matches in CheckPredicates
	Success MatchType = iota

	// Failure indicates a predicate does not match in CheckPredicates
	Failure
)

// MockTestOperation represents an invocation of a method on
// testing.TB and its arguments.
type MockTestOperation struct {
	// Name represents the name of an operation on TB. One of: Error,
	// Error, Fatal, Fatalf, Fail, FailNow, Log, Logf, Skip, SkipNow,
	// Skipf.
	Name string

	// Args are the formatted arguments of the operation. For the
	// formatted variants (e.g., Errorf), this is the result of invoke
	// fmt.Sprintf with the format string and arguments. For the
	// non-formatted variables, this is the result of the arguments,
	// converted to strings as in fmt.Sprint, with spaces between each
	// one argument. Finally for the no-arg operations this is an
	// empty string.
	Args string
}

// String provides a string representation of MockTestOperation.
func (op *MockTestOperation) String() string {
	return fmt.Sprintf("%s(%q)", op.Name, op.Args)
}

// Predicate represents a predicate testing a
// MockTestOperation. Typically Predicates are composed from a
// NamePredicate and an ArgPredicate via Match.
type Predicate func(MockTestOperation) MatchType

// NamePredicate compares the given operation name (the name of a
// testing.TB method) to determine if it matches the predicate.
type NamePredicate func(string) MatchType

// OneOfOp returns a NamePredicate requires its op to equal one of the
// given operations.
func OneOfOp(expectedOps ...string) NamePredicate {
	return func(op string) MatchType {
		for _, eop := range expectedOps {
			if op == eop {
				return Success
			}
		}
		return Failure
	}
}

// ExactOp is equivalent to OneOfOp(op).
func ExactOp(op string) NamePredicate { return OneOfOp(op) }

// ErrorOp matches Error or Errorf operations.
func ErrorOp() NamePredicate { return OneOfOp("Error", "Errorf") }

// FatalOp matches Fatal or Fatalf operations.
func FatalOp() NamePredicate { return OneOfOp("Fatal", "Fatalf") }

// Fails matches any testing.TB operation that causes a test failure.
func Fails() NamePredicate {
	return OneOfOp("Error", "Errorf", "Fatal", "Fatalf", "Fail", "FailNow")
}

// ArgsPredicate tests the given operation arguments to determine if
// they match the predicate.
type ArgsPredicate func(string) MatchType

// PrefixedArgs requires the operation arguments to match the given
// prefix.
func PrefixedArgs(prefix string) ArgsPredicate {
	return func(args string) MatchType {
		if strings.HasPrefix(args, prefix) {
			return Success
		}
		return Failure
	}
}

// ArgsContain requires the operation arguments to contain the given
// expression.
func ArgsContain(expr string) ArgsPredicate {
	return func(args string) MatchType {
		if strings.Contains(args, expr) {
			return Success
		}
		return Failure
	}
}

// Match creates a Predicate from the given NamePredicate and
// ArgPredicate. The NamePredicate is tested first. If it fails, the
// entire Predicate fails. Otherwise, the ArgPredicate is tested it's
// result is returned.
func Match(namePred NamePredicate, argPred ArgsPredicate) Predicate {
	return func(op MockTestOperation) MatchType {
		if m := namePred(op.Name); m != Success {
			return m
		}

		return argPred(op.Args)
	}
}

// MockT is an implementation of testing.TB that allows testing of
// code that makes test assertions. Calls to Fatal, Fatalf, FailNow
// and SkipNow are recorded, but do not terminate execution. If the
// code under test depends on termination of execution, unexpected
// results may occur.
type MockT struct {
	testing.TB

	Operations []MockTestOperation

	failed  bool
	skipped bool
	helper  bool
}

func (t *MockT) record(op string, args string) {
	t.Operations = append(t.Operations, MockTestOperation{op, args})
}

func (t *MockT) operationsToString() string {
	s := make([]string, 0, len(t.Operations))
	for _, e := range t.Operations {
		s = append(s, e.String())
	}
	return strings.Join(s, "\n")
}

// Reset resets the MockT's log of operations.
func (t *MockT) Reset() {
	t.Operations = nil
	t.failed = false
	t.skipped = false
}

// CheckHelper fails the test if the Helper method was not invoked on
// this MockT.
func (t *MockT) CheckHelper(realT testing.TB) {
	if !t.helper {
		Tracing(realT).Error("expected Helper invocation, but did not see it")
	}
}

// CheckPredicates evaluates the given predicates against the recorded
// testing.TB operations. Predicates must match operations in order.
func (t *MockT) CheckPredicates(realT testing.TB, predicates ...Predicate) {
	nops := len(t.Operations)
	npreds := len(predicates)

	n := nops
	if npreds < nops {
		n = npreds
	}

	errors := []string{}
	for i := 0; i < n; i++ {
		op := t.Operations[i]
		if predicates[i](op) != Success {
			errors = append(
				errors,
				fmt.Sprintf("operation %d, %s: did not match predicate", i+1, op.String()),
			)
		}
	}

	if nops != npreds {
		Tracing(realT).Errorf(
			"Expected %d operations, but got %d:\n%s",
			npreds,
			nops,
			t.operationsToString(),
		)
	}

	if len(errors) > 0 {
		Tracing(realT).Errorf("Predicate failure(s):\n%s", strings.Join(errors, "\n"))
	}
}

// CheckSuccess fails the test if any operations were logged.
func (t *MockT) CheckSuccess(realT testing.TB) {
	if len(t.Operations) > 0 {
		Tracing(realT).Errorf(
			"expected no testing.TB operations, but saw:\n%s",
			t.operationsToString(),
		)
	}
}

func fmtArgs(args ...interface{}) string {
	// Sprintln (vs Sprint) always adds spaces, but we don't want the
	// trailing newline.
	s := fmt.Sprintln(args...)
	return s[:len(s)-1]
}

// Name returns "MockTest", always.
func (t *MockT) Name() string { return "MockTest" }

// Error records an Error operation.
func (t *MockT) Error(args ...interface{}) {
	t.record("Error", fmtArgs(args...))
}

// Errorf records an Errorf operation.
func (t *MockT) Errorf(format string, args ...interface{}) {
	t.record("Errorf", fmt.Sprintf(format, args...))
}

// Fail records a Fail operation.
func (t *MockT) Fail() {
	t.record("Fail", "")
	t.failed = true
}

// FailNow records a FailNow operation.
func (t *MockT) FailNow() {
	t.record("FailNow", "")
	t.failed = true
}

// Failed reports whether Fail or FailNow was previously invoked.
func (t *MockT) Failed() bool {
	return t.failed
}

// Fatal records a Fatal operation.
func (t *MockT) Fatal(args ...interface{}) {
	t.record("Fatal", fmtArgs(args...))
}

// Fatalf records a Fatalf operation.
func (t *MockT) Fatalf(format string, args ...interface{}) {
	t.record("Fatalf", fmt.Sprintf(format, args...))
}

// Log records a Log operation.
func (t *MockT) Log(args ...interface{}) {
	t.record("Log", fmtArgs(args...))
}

// Logf records a Logf operation.
func (t *MockT) Logf(format string, args ...interface{}) {
	t.record("Logf", fmt.Sprintf(format, args...))
}

// Skip records a Skip operation.
func (t *MockT) Skip(args ...interface{}) {
	t.record("Skip", fmtArgs(args...))
	t.skipped = true
}

// Skipf records a Skipf operation.
func (t *MockT) Skipf(format string, args ...interface{}) {
	t.record("Skipf", fmt.Sprintf(format, args...))
	t.skipped = true
}

// SkipNow records a SkipNow operation.
func (t *MockT) SkipNow() {
	t.record("SkipNow", "")
	t.skipped = true
}

// Skipped reports whether Skip, Skipf, or SkipNow was previously
// invoked.
func (t *MockT) Skipped() bool {
	return t.skipped
}

// Helper records a Helper operation.
func (t *MockT) Helper() {
	t.helper = true
}
