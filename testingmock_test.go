package testingmock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TestingMock_Should(t *testing.T) {
	t.Parallel()

	t.Run("not be failed by default", func(t *testing.T) {
		t.Parallel()
		// Arrange
		tMock := New(t)

		// Assert
		assert.False(t, tMock.Failed())
	})

	t.Run("be failed when Error() has been called", func(t *testing.T) {
		t.Parallel()
		// Arrange
		tMock := New(t)
		// Act
		tMock.Error("an error occurred")
		// Assert
		assert.True(t, tMock.Failed())
	})

	t.Run("be failed when Fatal() has been called", func(t *testing.T) {
		t.Parallel()
		// Arrange
		tMock := New(t)
		// Act
		tMock.Fatal("a fatal error occurred")
		// Assert
		assert.True(t, tMock.Failed())
	})

	t.Run("be failed when FailNow() has been called", func(t *testing.T) {
		t.Parallel()
		// Arrange
		tMock := New(t)
		// Act
		tMock.FailNow()
		// Assert
		assert.True(t, tMock.Failed())
	})

	t.Run("provide assertions", func(t *testing.T) {
		t.Parallel()

		t.Run("AssertDidNotFailed()", func(t *testing.T) {
			t.Parallel()
			t.Run("to not fail the encapsulated testing.IT when actually not failed", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act & Assert
				tMock.AssertDidNotFailed()
				assert.False(t, encapsulatedT.Failed())
			})

			var testFailureCases = []struct {
				name string
				fail func(tMock *MockedT)
			}{
				{
					name: "Error()",
					fail: func(tMock *MockedT) {
						tMock.Error("an error occurred")
					},
				},
				{
					name: "Errorf()",
					fail: func(tMock *MockedT) {
						tMock.Errorf("an error occurred")
						tMock.Errorf("an %s occurred", "error")
					},
				},
				{
					name: "Fatal()",
					fail: func(tMock *MockedT) {
						tMock.Fatal("a fatal error occurred")
					},
				},
				{
					name: "Fatalf()",
					fail: func(tMock *MockedT) {
						tMock.Fatalf("a fatal error occurred")
						tMock.Fatalf("a %s error occurred", "fatal")
					},
				},
				{
					name: "FailNow()",
					fail: func(tMock *MockedT) {
						tMock.FailNow()
					},
				},
			}

			for _, tc := range testFailureCases {
				t.Run("to fail the encapsulated testing.IT when actually failed with "+tc.name, func(t *testing.T) {
					t.Parallel()
					// Arrange
					encapsulatedT := &testing.T{}
					tMock := New(encapsulatedT)

					// Act
					tc.fail(tMock)
					tMock.AssertDidNotFailed()

					// Assert
					assert.True(t, encapsulatedT.Failed())
				})
			}
		})

		t.Run("AssertFailedWithError()", func(t *testing.T) {
			t.Parallel()
			t.Run("to not fail the encapsulated testing.IT when failed with Error()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Error("an error occurred")
				tMock.AssertFailedWithError()

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to not fail the encapsulated testing.IT when failed with Errorf()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Errorf("an error occurred")
				tMock.AssertFailedWithError()

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when no failing method was called", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.AssertFailedWithError()

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when failing for another reason than Error() or Errorf()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)
				tMock.Fatal()
				tMock.FailNow()

				// Act
				tMock.AssertFailedWithError()

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
		})

		t.Run("AssertFailedWithErrorMessage()", func(t *testing.T) {
			t.Parallel()
			t.Run("to not fail the encapsulated testing.IT when failed with Error() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Error("an error occurred")
				tMock.AssertFailedWithErrorMessage("an error occurred")

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to not fail the encapsulated testing.IT when failed with Errorf() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Errorf("an %s error occurred", "unexpected")
				tMock.AssertFailedWithErrorMessage("an unexpected error occurred")

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when it did not fail with Error() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Error("an error occurred")
				tMock.AssertFailedWithErrorMessage("another error occurred")

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when it did not fail with Errorf() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock2 := New(encapsulatedT)

				// Act
				tMock2.Errorf("an %s error occurred", "unexpected")
				tMock2.AssertFailedWithErrorMessage("another unexpected error occurred")

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
		})

		t.Run("AssertFailedWithFatalMessage()", func(t *testing.T) {
			t.Parallel()
			t.Run("to not fail the encapsulated testing.IT when failed with Fatal() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Fatal("an error occurred")
				tMock.AssertFailedWithFatalMessage("an error occurred")

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to not fail the encapsulated testing.IT when failed with Fatalf() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Fatalf("an %s error occurred", "unexpected")
				tMock.AssertFailedWithFatalMessage("an unexpected error occurred")

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when it did not fail with Fatal() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Fatal("an error occurred")
				tMock.AssertFailedWithFatalMessage("another error occurred")

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when it did not fail with Fatalf() and the expected message", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock2 := New(encapsulatedT)

				// Act
				tMock2.Fatalf("an %s error occurred", "unexpected")
				tMock2.AssertFailedWithFatalMessage("another unexpected error occurred")

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
		})

		t.Run("AssertFailedWithFatal()", func(t *testing.T) {
			t.Parallel()
			t.Run("to not fail the encapsulated testing.IT when failed with Fatal()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Fatal("a fatal error occurred")
				tMock.AssertFailedWithFatal()

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to not fail the encapsulated testing.IT when failed with Fatalf()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.Fatalf("a fatal error occurred")
				tMock.AssertFailedWithFatal()

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when no failing method was called", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.AssertFailedWithFatal()

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when failing for another reason than Fatal() or Fatalf()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)
				tMock.Error("an error occurred")
				tMock.FailNow()

				// Act
				tMock.AssertFailedWithFatal()

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
		})

		t.Run("AssertFailNowHasBeenCalled()", func(t *testing.T) {
			t.Parallel()
			t.Run("to not fail the encapsulated testing.IT when FailNow() actually has been called", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.FailNow()
				tMock.AssertFailNowHasBeenCalled()

				// Assert
				assert.False(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when no failing method was called", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)

				// Act
				tMock.AssertFailNowHasBeenCalled()

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
			t.Run("to fail the encapsulated testing.IT when failing for another reason than FailNow()", func(t *testing.T) {
				t.Parallel()
				// Arrange
				encapsulatedT := &testing.T{}
				tMock := New(encapsulatedT)
				tMock.Error("an error occurred")
				tMock.Fatal("an error occurred")

				// Act
				tMock.AssertFailNowHasBeenCalled()

				// Assert
				assert.True(t, encapsulatedT.Failed())
			})
		})
	})

	t.Run("provide GetCleanups() to return registered cleanups", func(t *testing.T) {
		t.Parallel()

		// Arrange
		encapsulatedT := &testing.T{}
		tMock := New(encapsulatedT)
		i := 0

		// Act
		tMock.Cleanup(func() { i += 4 })
		tMock.Cleanup(func() { i -= 2 })
		cleanups := tMock.GetCleanups()
		for _, cleanup := range cleanups {
			cleanup()
		}

		// Assert
		assert.Equal(t, 2, i)
	})
}
