package testingmock

import (
	"fmt"
	"slices"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// IT is an interface wrapper around *testing.T. It also implements the interfaces required by testify's assert,
// require and mock packages. This is useful for creating mocks of testing.T in order to test tests fixtures and utilities.
type IT interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	FailNow()
	Log(args ...any)
	Logf(format string, args ...any)
	Failed() bool
	Cleanup(f func())
}

var _ require.TestingT = (IT)(nil)
var _ assert.TestingT = (IT)(nil)
var _ mock.TestingT = (IT)(nil)
var _ IT = (*MockedT)(nil)

// MockedT is a mock implementation of IT that can be used to test fixtures and utilities.
type MockedT struct {
	t        IT
	m        mock.Mock
	assert   *assert.Assertions
	cleanups []func()
}

// New creates a new instance of MockedT.
func New(t IT) *MockedT {
	mockedT := &MockedT{
		t:      t,
		assert: assert.New(t),
	}

	mockedT.m.On("Error", mock.Anything).Return()
	mockedT.m.On("Errorf", mock.Anything, mock.Anything).Return()
	mockedT.m.On("Fatal", mock.Anything).Return()
	mockedT.m.On("Fatalf", mock.Anything, mock.Anything).Return()
	mockedT.m.On("FailNow").Return()
	mockedT.m.On("Log", mock.Anything).Return()
	mockedT.m.On("Logf", mock.Anything, mock.Anything).Return()

	return mockedT
}

// Error mocked implementation.
func (testState *MockedT) Error(args ...any) {
	_ = testState.m.Called(args...)
	_ = args
}

// Errorf mocked implementation.
func (testState *MockedT) Errorf(format string, args ...any) {
	allArgs := make([]any, 0, 1+len(args))
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	_ = testState.m.Called(allArgs...)
}

// Fatal mocked implementation.
func (testState *MockedT) Fatal(args ...any) {
	_ = testState.m.Called(args...)
	_ = args
}

// Fatalf mocked implementation.
func (testState *MockedT) Fatalf(format string, args ...any) {
	allArgs := make([]any, 0, 1+len(args))
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	_ = testState.m.Called(allArgs...)
}

// FailNow mocked implementation.
func (testState *MockedT) FailNow() {
	_ = testState.m.Called()
}

// Log mocked implementation.
func (testState *MockedT) Log(args ...any) {
	_ = testState.m.Called(args...)
}

// Logf mocked implementation.
func (testState *MockedT) Logf(format string, args ...any) {
	allArgs := make([]any, 0, 1+len(args))
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	_ = testState.m.Called(allArgs...)
}

// Failed mocked implementation.
func (testState *MockedT) Failed() bool {
	for _, call := range testState.m.Calls {
		if slices.Contains([]string{
			"Error",
			"Errorf",
			"Fatal",
			"Fatalf",
			"FailNow",
		}, call.Method) {
			return true
		}
	}
	return false
}

// Cleanup mocked implementation.
func (testState *MockedT) Cleanup(f func()) {
	testState.cleanups = append(testState.cleanups, f)
}

// AssertDidNotFailed asserts that no failure methods were called.
func (testState *MockedT) AssertDidNotFailed() {
	testState.m.AssertNotCalled(testState.t, "Error", mock.Anything)
	testState.m.AssertNotCalled(testState.t, "Errorf", mock.Anything, mock.Anything)
	testState.m.AssertNotCalled(testState.t, "Fatal", mock.Anything)
	testState.m.AssertNotCalled(testState.t, "Fatalf", mock.Anything, mock.Anything)
	testState.m.AssertNotCalled(testState.t, "FailNow")
}

// AssertFailedWithError asserts that an error method (Error, Errorf) was called.
func (testState *MockedT) AssertFailedWithError() {
	for _, call := range testState.m.Calls {
		if call.Method == "Error" || call.Method == "Errorf" {
			return
		}
	}
	testState.t.Error("an error was expected to occur but did not")
}

// AssertFailedWithErrorMessage asserts that an error method (Error, Errorf) was called with the expected message.
func (testState *MockedT) AssertFailedWithErrorMessage(expectedMessage string) {
	call := testState.findErrorCallWithMessage(levelError, expectedMessage)
	if call == nil {
		testState.t.Errorf("an error with message '%s' was expected to occur but did not", expectedMessage)
	}
}

// AssertFailedWithFatal asserts that a fatal method (Fatal, Fatalf) was called.
func (testState *MockedT) AssertFailedWithFatal() {
	for _, call := range testState.m.Calls {
		if call.Method == "Fatal" || call.Method == "Fatalf" {
			return
		}
	}
	testState.t.Error("a fatal was expected to occur but did not")
}

// AssertFailedWithFatalMessage asserts that a fatal method (Fatal, Fatalf) was called with the expected message.
func (testState *MockedT) AssertFailedWithFatalMessage(expectedMessage string) {
	call := testState.findErrorCallWithMessage(levelFatal, expectedMessage)
	if call == nil {
		testState.t.Errorf("a fatal error with message '%s' was expected to occur but did not", expectedMessage)
	}
}

// AssertFailNowHasBeenCalled asserts that FailNow was called.
func (testState *MockedT) AssertFailNowHasBeenCalled() {
	testState.m.AssertCalled(testState.t, "FailNow")
}

func (testState *MockedT) GetCleanups() []func() {
	return append([]func(){}, testState.cleanups...)
}

type errorLevel int

const (
	levelError errorLevel = 0
	levelFatal errorLevel = 1
)

func (testState *MockedT) findErrorCallWithMessage(level errorLevel, expectedMessage string) *mock.Call {
	m1 := "Error"
	m2 := "Errorf"
	if level == levelFatal {
		m1 = "Fatal"
		m2 = "Fatalf"
	}
	for _, call := range testState.m.Calls {
		if call.Method == m1 {
			args := make([]string, len(call.Arguments))
			for i, arg := range call.Arguments {
				args[i] = fmt.Sprintf("%v", arg)
			}
			if strings.Join(args, " ") == expectedMessage {
				return &call
			}
		}
		if call.Method == m2 {
			f := call.Arguments.Get(0).(string)
			args := make([]any, 0, len(call.Arguments)-1)
			for i := 1; i < len(call.Arguments); i++ {
				args = append(args, call.Arguments.Get(i))
			}
			msg := fmt.Sprintf(f, args...)
			if msg == expectedMessage {
				return &call
			}
		}
	}
	return nil
}
