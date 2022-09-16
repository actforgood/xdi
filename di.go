// Copyright The ActForGood Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://github.com/actforgood/xdi/LICENSE.

package xdi

const (
	// FilterAllRegistered specifies that all registered object IDs should be returned by [Manager.ListIDs].
	FilterAllRegistered byte = 1
	// FilterInitializedShared specifies that only shared objects that have been already initialized
	// (there was at least one [Manager.Get] call for them) should be returned by [Manager.ListIDs].
	FilterInitializedShared byte = 2
)

type (
	// InitializerFunc is a simple function that returns an instance of a dependency,
	// whatever that dependency is (an object / a function / a config, etc.).
	InitializerFunc func() interface{}

	// Manager is a container for dependencies and their definitions.
	// Its APIs are not concurrent safe for use.
	Manager struct {
		// def contains initialization definitions indexed by a string identifier.
		def map[string]Definition
		// sharedRegistry contains shared instances.
		sharedRegistry map[string]interface{}
	}

	// Definition holds information needed to initialize a dependency.
	Definition struct {
		// ID is an identifier of a dependency.
		// Example: "app.logger" / "logger" / "config.timeout"
		ID string
		// Initializer is a factory/function that returns a dependency.
		Initializer InitializerFunc
		// Shared is a flag whether this dependency is shared or not.
		// If it is shared, multiple Get() calls for the same identifier will return "same" instance from a registry,
		// see also the note below.
		// If it is not shared, multiple Get() calls for the same identifier will return "different" instances
		// (Initializer function will be called on each Get() call).
		// Note: it makes sense for Initializer function to produce a pointer/reference type if
		// you really want it to be shared. On the same principle, if your Initializer function returns
		// a singleton for example, shared flag is useless, that instance will always be shared.
		Shared bool
	}
)

// NewManager instantiates a new dependency injection manager.
func NewManager() *Manager {
	return &Manager{
		def:            make(map[string]Definition),
		sharedRegistry: make(map[string]interface{}),
	}
}

// AddDefinition adds new definition for a dependency instantiation/retrieval.
// Multiple calls with same definition ID will overwrite previous definition.
func (diMngr *Manager) AddDefinition(def Definition) {
	diMngr.def[def.ID] = def
}

// Get returns a dependency if a definition for it was provided previously, or nil otherwise.
func (diMngr *Manager) Get(id string) interface{} {
	// look for definition.
	def, foundDef := diMngr.def[id]
	if !foundDef {
		// if it's a shared dependency already instantiated, return it.
		if dep, foundRegistry := diMngr.sharedRegistry[id]; foundRegistry {
			return dep
		}

		return nil
	}

	dep := def.Initializer()

	if def.Shared {
		// store instance in shared registry.
		diMngr.sharedRegistry[id] = dep
		// if it's a shared dependency we can dispose of its definition and free memory,
		// dependency will be returned from sharedRegistry on an eventual next Get() call for it.
		delete(diMngr.def, id)
	}

	return dep
}

// ListIDs returns the IDs list of all registered definitions,
// or of shared and initialized objects if filter is specified.
// Example of usage:
//
//	// appShutDown closes all resources at application shutdown.
//	func appShutDown(diManager *xdi.Manager) {
//	    for _, id := range diManager.ListIDs(xdi.FilterInitializedShared) {
//	        if closable, ok := diManager.Get(id).(io.Closer); ok {
//	            if err := closable.Close(); err != nil {
//	                fmt.Printf("could not close resource '%s': %+v\n", id, err)
//	            }
//	        }
//	    }
//	}
func (diMngr *Manager) ListIDs(filter ...byte) []string {
	var (
		fltr = FilterAllRegistered
		IDs  []string
		idx  = 0
	)
	if len(filter) > 0 && filter[0] == FilterInitializedShared {
		fltr = FilterInitializedShared
	}

	switch fltr {
	case FilterAllRegistered:
		IDs = make([]string, len(diMngr.sharedRegistry)+len(diMngr.def))
		for id := range diMngr.def {
			IDs[idx] = id
			idx++
		}
		for id := range diMngr.sharedRegistry {
			IDs[idx] = id
			idx++
		}
	case FilterInitializedShared:
		IDs = make([]string, len(diMngr.sharedRegistry))
		for id := range diMngr.sharedRegistry {
			IDs[idx] = id
			idx++
		}
	}

	return IDs
}
