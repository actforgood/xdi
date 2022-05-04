// Copyright 2022 Bogdan Constantinescu.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://github.com/actforgood/xdi/LICENSE.

package xdi_test

import (
	"strconv"
	"testing"

	"github.com/actforgood/xdi"
)

func TestDiManager_Get(t *testing.T) {
	t.Parallel()

	t.Run("shared dependency", testDiManagerGetShared)
	t.Run("not shared dependency", testDiManagerGetNotShared)
	t.Run("not found dependency", testDiManagerGetNotFound)
}

type dummy struct {
	age int
}

// testDiManagerGetShared tests same instance is returned for a shared object.
func testDiManagerGetShared(t *testing.T) {
	t.Parallel()

	// arrange
	var (
		subject = xdi.NewDiManager()
		testID  = "test.dummy"
	)
	subject.AddDefinition(xdi.DiManagerDef{
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

// testDiManagerGetNotShared tests different instances are returned for a non-shared object.
func testDiManagerGetNotShared(t *testing.T) {
	t.Parallel()

	// arrange
	var (
		subject = xdi.NewDiManager()
		testID  = "test.dummy"
	)
	subject.AddDefinition(xdi.DiManagerDef{
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

// testDiManagerGetNotFound tests that nil is returned for an unknown id.
func testDiManagerGetNotFound(t *testing.T) {
	t.Parallel()

	// arrange
	subject := xdi.NewDiManager()

	// act
	result := subject.Get("unknown")

	// assert
	if result != nil {
		t.Errorf("expected '%+v', but got '%+v'", nil, result)
	}
}

func BenchmarkDiManager_Get_shared(b *testing.B) {
	subject := xdi.NewDiManager()
	for i := 0; i < 50; i++ {
		age := i + 1
		subject.AddDefinition(xdi.DiManagerDef{
			ID: "dummy.age." + strconv.FormatInt(int64(age), 10),
			Initializer: func() interface{} {
				return &dummy{age: age}
			},
			Shared: true,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = subject.Get("dummy.age.35")
	}
}

func BenchmarkDiManager_Get_notShared(b *testing.B) {
	subject := xdi.NewDiManager()
	for i := 0; i < 50; i++ {
		age := i + 1
		subject.AddDefinition(xdi.DiManagerDef{
			ID: "dummy.age." + strconv.FormatInt(int64(age), 10),
			Initializer: func() interface{} {
				return &dummy{age: age}
			},
			Shared: false,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = subject.Get("dummy.age.35")
	}
}
