// Copyright 2022 Bogdan Constantinescu.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://github.com/actforgood/xdi/LICENSE.

package xdi

type (
	// InitializerFunc is a simple function that returns an instance of a dependency,
	// whatever that dependency is (an object / a function / a config, etc.).
	InitializerFunc func() interface{}

	// DiManager is a container for dependencies and their definitions.
	// Its APIs are not concurrent safe for use.
	DiManager struct {
		// def contains initialization definitions indexed by a string identifier.
		def map[string]DiManagerDef
		// sharedRegistry contains shared instances.
		sharedRegistry map[string]interface{}
	}

	// DiManagerDef holds information needed to initialize a dependency.
	DiManagerDef struct {
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

// NewDiManager instantiates a new dependency injection manager.
func NewDiManager() *DiManager {
	return &DiManager{
		def:            make(map[string]DiManagerDef),
		sharedRegistry: make(map[string]interface{}),
	}
}

// AddDefinition adds new definition for a dependency instantiation/retrieval.
// Multiple calls with same definition ID will overwrite previous definition.
func (om *DiManager) AddDefinition(def DiManagerDef) {
	om.def[def.ID] = def
}

// Get returns a dependency if a definition for it was provided previously, or nil otherwise.
func (om *DiManager) Get(id string) interface{} {
	// look for definition.
	def, foundDef := om.def[id]
	if !foundDef {
		// if it's a shared dependency already instantiated, return it.
		if dep, foundRegistry := om.sharedRegistry[id]; foundRegistry {
			return dep
		}

		return nil
	}

	dep := def.Initializer()

	if def.Shared {
		// store instance in shared registry.
		om.sharedRegistry[id] = dep
		// if it's a shared dependency we can dispose of its definition and free memory,
		// dependency will be returned from sharedRegistry on an eventual next Get() call for it.
		delete(om.def, id)
	}

	return dep
}
