package samples

import (
	"testing"

	"github.com/le-yams/gotestingmock"
)

func Test_UtilityAssertionPassTestWhenConditionIsMet(t *testing.T) {
	// Arrange
	mockedT := testingmock.New(t)
	util := NewTestingUtility(mockedT)

	// Act
	util.conditionMet = true
	util.AssertStuff()

	mockedT.AssertDidNotFailed()
}

func Test_UtilityAssertionFailsTestWhenConditionIsNotMet(t *testing.T) {
	mockedT := testingmock.New(t)

	util := NewTestingUtility(mockedT)

	util.conditionMet = false

	util.AssertStuff()

	mockedT.AssertFailedWithErrorMessage("Condition not met!")
}
