// Copyright The ActForGood Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://github.com/actforgood/xdi/blob/main/LICENSE

package xdi_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/actforgood/xdi"
)

type dummy struct {
	age int64
}

func TestManager_Get(t *testing.T) {
	t.Parallel()

	t.Run("shared dependency", testManagerGetShared)
	t.Run("not shared dependency", testManagerGetNotShared)
	t.Run("not found dependency", testManagerGetNotFound)
}

// testManagerGetShared tests same instance is returned for a shared object.
func testManagerGetShared(t *testing.T) {
	t.Parallel()

	// arrange
	var (
		subject = xdi.NewManager()
		testID  = "test.dummy"
	)
	subject.AddDefinition(xdi.Definition{
		ID: testID,
		Initializer: func() interface{} {
			return &dummy{age: 35}
		},
		Shared: true,
	})

	// act & assert
	result := subject.Get(testID)
	if result == nil {
		t.Fatalf("expected '%+v', but got '%+v'\n", "*dummy", nil)
	}
	for i := 0; i < 10; i++ {
		resultAgain := subject.Get(testID)
		if result != resultAgain {
			t.Errorf("expected '%+v', but got '%+v'\n", result, resultAgain)
		}
	}
}

// testManagerGetNotShared tests different instances are returned for a non-shared object.
func testManagerGetNotShared(t *testing.T) {
	t.Parallel()

	// arrange
	var (
		subject = xdi.NewManager()
		testID  = "test.dummy"
	)
	subject.AddDefinition(xdi.Definition{
		ID: testID,
		Initializer: func() interface{} {
			return &dummy{age: 35}
		},
		Shared: false,
	})

	// act & assert
	result := subject.Get(testID)
	if result == nil {
		t.Fatalf("expected '%+v', but got '%+v'\n", "*dummy", nil)
	}
	for i := 0; i < 10; i++ {
		resultAgain := subject.Get(testID)
		if result == resultAgain {
			t.Errorf("expected not to be the same, trial = %+v\n", i+1)
		}
	}
}

// testManagerGetNotFound tests that nil is returned for an unknown id.
func testManagerGetNotFound(t *testing.T) {
	t.Parallel()

	// arrange
	subject := xdi.NewManager()

	// act
	result := subject.Get("unknown")

	// assert
	if result != nil {
		t.Errorf("expected '%+v', but got '%+v'", nil, result)
	}
}

func TestManager_ListIDs(t *testing.T) {
	t.Parallel()

	// arrange
	var (
		subject        = setUpDiManager()
		expectedShared = []string{"dummy.age.shared.5", "dummy.age.shared.6"}
		expectedAll    = make([]string, 0, 20)
		idTpl          = "dummy.age.%s.%d"
		id             string
	)
	for i := 1; i <= 50; i++ {
		if i > 30 {
			id = fmt.Sprintf(idTpl, "notshared", i)
		} else {
			id = fmt.Sprintf(idTpl, "shared", i)
		}
		expectedAll = append(expectedAll, id)
	}
	_ = subject.Get("dummy.age.shared.5")
	_ = subject.Get("dummy.age.shared.6")
	_ = subject.Get("dummy.age.notshared.15")

	// act 1
	result1 := subject.ListIDs(xdi.FilterInitializedShared)

	// assert 1
	if len(result1) != len(expectedShared) {
		t.Errorf("expected '%+v', but got '%+v'", len(expectedShared), len(result1))
	}
	for _, id := range expectedShared {
		if !isInSlice(id, result1) {
			t.Errorf("expected id '%+v' to be found", id)
		}
	}

	// act 2
	result2 := subject.ListIDs()

	// assert 2
	if len(result2) != len(expectedAll) {
		t.Errorf("expected '%+v', but got '%+v'", len(expectedAll), len(result2))
	}
	for _, id := range expectedAll {
		if !isInSlice(id, result2) {
			t.Errorf("expected id '%+v' to be found", id)
		}
	}
}

func TestManagerInstance(t *testing.T) {
	t.Parallel()

	// act
	subject := xdi.ManagerInstance()

	// assert
	if subject == nil {
		t.Fatalf("expected '%+v', but got '%+v'\n", "*Manager", nil)
	}

	for i := 0; i < 10; i++ {
		// act
		sameObject := xdi.ManagerInstance()

		// assert
		if sameObject != subject {
			t.Fatalf("expected '%p', but got '%p'\n", subject, sameObject)
		}
	}
}

func BenchmarkManager_Get_shared(b *testing.B) {
	diManager := setUpDiManager()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = diManager.Get("dummy.age.shared.20")
	}
}

func BenchmarkManager_Get_notShared(b *testing.B) {
	diManager := setUpDiManager()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = diManager.Get("dummy.age.notshared.40")
	}
}

func BenchmarkManager_ListIDs(b *testing.B) {
	diManager := setUpDiManager()
	_ = diManager.Get("dummy.age.shared.20")
	_ = diManager.Get("dummy.age.shared.21")

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = diManager.ListIDs()
	}
}

// setUpDiManager makes a Manager with 50 registered objects,
// 30 - shared, 20 - not shared.
func setUpDiManager() *xdi.Manager {
	var (
		diManager = xdi.NewManager()
		i         int64
	)
	for i = 1; i <= 50; i++ {
		id := "dummy.age."
		shared := true
		if i > 30 {
			shared = false
			id += "notshared."
		} else {
			id += "shared."
		}
		age := i // capture range variable
		id += strconv.FormatInt(age, 10)
		diManager.AddDefinition(xdi.Definition{
			ID: id,
			Initializer: func() interface{} {
				return &dummy{age: age}
			},
			Shared: shared,
		})
	}

	return diManager
}

// isInSlice searches for needle in the haystack and returns true if it is found.
func isInSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}
