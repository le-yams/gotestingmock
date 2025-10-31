package samples

import "github.com/le-yams/gotestingmock"

type T testingmock.IT

type TestingUtility struct {
	t            T
	conditionMet bool
}

func NewTestingUtility(t T) *TestingUtility {
	u := &TestingUtility{t: t}
	t.Cleanup(u.cleanup)
	return u
}

func (m *TestingUtility) AssertStuff() {
	if !m.conditionMet {
		m.t.Error("Condition not met!")
	}
}

func (m *TestingUtility) cleanup() {
	// Any necessary cleanup actions
}
