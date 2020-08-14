// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.

package mocks

import (
	serialization "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/serialization"
	types "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/types"
	"sync"
)

// MockSerializer is a mock implementation of the Serializer interface (from
// the package
// github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/persistence/serialization)
// used for unit testing.
type MockSerializer struct {
	// MarshalDocumentDataFunc is an instance of a mock function object
	// controlling the behavior of the method MarshalDocumentData.
	MarshalDocumentDataFunc *SerializerMarshalDocumentDataFunc
	// MarshalLocationsFunc is an instance of a mock function object
	// controlling the behavior of the method MarshalLocations.
	MarshalLocationsFunc *SerializerMarshalLocationsFunc
	// MarshalResultChunkDataFunc is an instance of a mock function object
	// controlling the behavior of the method MarshalResultChunkData.
	MarshalResultChunkDataFunc *SerializerMarshalResultChunkDataFunc
	// UnmarshalDocumentDataFunc is an instance of a mock function object
	// controlling the behavior of the method UnmarshalDocumentData.
	UnmarshalDocumentDataFunc *SerializerUnmarshalDocumentDataFunc
	// UnmarshalLocationsFunc is an instance of a mock function object
	// controlling the behavior of the method UnmarshalLocations.
	UnmarshalLocationsFunc *SerializerUnmarshalLocationsFunc
	// UnmarshalResultChunkDataFunc is an instance of a mock function object
	// controlling the behavior of the method UnmarshalResultChunkData.
	UnmarshalResultChunkDataFunc *SerializerUnmarshalResultChunkDataFunc
}

// NewMockSerializer creates a new mock of the Serializer interface. All
// methods return zero values for all results, unless overwritten.
func NewMockSerializer() *MockSerializer {
	return &MockSerializer{
		MarshalDocumentDataFunc: &SerializerMarshalDocumentDataFunc{
			defaultHook: func(types.DocumentData) ([]byte, error) {
				return nil, nil
			},
		},
		MarshalLocationsFunc: &SerializerMarshalLocationsFunc{
			defaultHook: func([]types.Location) ([]byte, error) {
				return nil, nil
			},
		},
		MarshalResultChunkDataFunc: &SerializerMarshalResultChunkDataFunc{
			defaultHook: func(types.ResultChunkData) ([]byte, error) {
				return nil, nil
			},
		},
		UnmarshalDocumentDataFunc: &SerializerUnmarshalDocumentDataFunc{
			defaultHook: func([]byte) (types.DocumentData, error) {
				return types.DocumentData{}, nil
			},
		},
		UnmarshalLocationsFunc: &SerializerUnmarshalLocationsFunc{
			defaultHook: func([]byte) ([]types.Location, error) {
				return nil, nil
			},
		},
		UnmarshalResultChunkDataFunc: &SerializerUnmarshalResultChunkDataFunc{
			defaultHook: func([]byte) (types.ResultChunkData, error) {
				return types.ResultChunkData{}, nil
			},
		},
	}
}

// NewMockSerializerFrom creates a new mock of the MockSerializer interface.
// All methods delegate to the given implementation, unless overwritten.
func NewMockSerializerFrom(i serialization.Serializer) *MockSerializer {
	return &MockSerializer{
		MarshalDocumentDataFunc: &SerializerMarshalDocumentDataFunc{
			defaultHook: i.MarshalDocumentData,
		},
		MarshalLocationsFunc: &SerializerMarshalLocationsFunc{
			defaultHook: i.MarshalLocations,
		},
		MarshalResultChunkDataFunc: &SerializerMarshalResultChunkDataFunc{
			defaultHook: i.MarshalResultChunkData,
		},
		UnmarshalDocumentDataFunc: &SerializerUnmarshalDocumentDataFunc{
			defaultHook: i.UnmarshalDocumentData,
		},
		UnmarshalLocationsFunc: &SerializerUnmarshalLocationsFunc{
			defaultHook: i.UnmarshalLocations,
		},
		UnmarshalResultChunkDataFunc: &SerializerUnmarshalResultChunkDataFunc{
			defaultHook: i.UnmarshalResultChunkData,
		},
	}
}

// SerializerMarshalDocumentDataFunc describes the behavior when the
// MarshalDocumentData method of the parent MockSerializer instance is
// invoked.
type SerializerMarshalDocumentDataFunc struct {
	defaultHook func(types.DocumentData) ([]byte, error)
	hooks       []func(types.DocumentData) ([]byte, error)
	history     []SerializerMarshalDocumentDataFuncCall
	mutex       sync.Mutex
}

// MarshalDocumentData delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockSerializer) MarshalDocumentData(v0 types.DocumentData) ([]byte, error) {
	r0, r1 := m.MarshalDocumentDataFunc.nextHook()(v0)
	m.MarshalDocumentDataFunc.appendCall(SerializerMarshalDocumentDataFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the MarshalDocumentData
// method of the parent MockSerializer instance is invoked and the hook
// queue is empty.
func (f *SerializerMarshalDocumentDataFunc) SetDefaultHook(hook func(types.DocumentData) ([]byte, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// MarshalDocumentData method of the parent MockSerializer instance inovkes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *SerializerMarshalDocumentDataFunc) PushHook(hook func(types.DocumentData) ([]byte, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SerializerMarshalDocumentDataFunc) SetDefaultReturn(r0 []byte, r1 error) {
	f.SetDefaultHook(func(types.DocumentData) ([]byte, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SerializerMarshalDocumentDataFunc) PushReturn(r0 []byte, r1 error) {
	f.PushHook(func(types.DocumentData) ([]byte, error) {
		return r0, r1
	})
}

func (f *SerializerMarshalDocumentDataFunc) nextHook() func(types.DocumentData) ([]byte, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SerializerMarshalDocumentDataFunc) appendCall(r0 SerializerMarshalDocumentDataFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SerializerMarshalDocumentDataFuncCall
// objects describing the invocations of this function.
func (f *SerializerMarshalDocumentDataFunc) History() []SerializerMarshalDocumentDataFuncCall {
	f.mutex.Lock()
	history := make([]SerializerMarshalDocumentDataFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SerializerMarshalDocumentDataFuncCall is an object that describes an
// invocation of method MarshalDocumentData on an instance of
// MockSerializer.
type SerializerMarshalDocumentDataFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 types.DocumentData
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []byte
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SerializerMarshalDocumentDataFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SerializerMarshalDocumentDataFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// SerializerMarshalLocationsFunc describes the behavior when the
// MarshalLocations method of the parent MockSerializer instance is invoked.
type SerializerMarshalLocationsFunc struct {
	defaultHook func([]types.Location) ([]byte, error)
	hooks       []func([]types.Location) ([]byte, error)
	history     []SerializerMarshalLocationsFuncCall
	mutex       sync.Mutex
}

// MarshalLocations delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockSerializer) MarshalLocations(v0 []types.Location) ([]byte, error) {
	r0, r1 := m.MarshalLocationsFunc.nextHook()(v0)
	m.MarshalLocationsFunc.appendCall(SerializerMarshalLocationsFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the MarshalLocations
// method of the parent MockSerializer instance is invoked and the hook
// queue is empty.
func (f *SerializerMarshalLocationsFunc) SetDefaultHook(hook func([]types.Location) ([]byte, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// MarshalLocations method of the parent MockSerializer instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *SerializerMarshalLocationsFunc) PushHook(hook func([]types.Location) ([]byte, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SerializerMarshalLocationsFunc) SetDefaultReturn(r0 []byte, r1 error) {
	f.SetDefaultHook(func([]types.Location) ([]byte, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SerializerMarshalLocationsFunc) PushReturn(r0 []byte, r1 error) {
	f.PushHook(func([]types.Location) ([]byte, error) {
		return r0, r1
	})
}

func (f *SerializerMarshalLocationsFunc) nextHook() func([]types.Location) ([]byte, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SerializerMarshalLocationsFunc) appendCall(r0 SerializerMarshalLocationsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SerializerMarshalLocationsFuncCall objects
// describing the invocations of this function.
func (f *SerializerMarshalLocationsFunc) History() []SerializerMarshalLocationsFuncCall {
	f.mutex.Lock()
	history := make([]SerializerMarshalLocationsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SerializerMarshalLocationsFuncCall is an object that describes an
// invocation of method MarshalLocations on an instance of MockSerializer.
type SerializerMarshalLocationsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 []types.Location
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []byte
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SerializerMarshalLocationsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SerializerMarshalLocationsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// SerializerMarshalResultChunkDataFunc describes the behavior when the
// MarshalResultChunkData method of the parent MockSerializer instance is
// invoked.
type SerializerMarshalResultChunkDataFunc struct {
	defaultHook func(types.ResultChunkData) ([]byte, error)
	hooks       []func(types.ResultChunkData) ([]byte, error)
	history     []SerializerMarshalResultChunkDataFuncCall
	mutex       sync.Mutex
}

// MarshalResultChunkData delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockSerializer) MarshalResultChunkData(v0 types.ResultChunkData) ([]byte, error) {
	r0, r1 := m.MarshalResultChunkDataFunc.nextHook()(v0)
	m.MarshalResultChunkDataFunc.appendCall(SerializerMarshalResultChunkDataFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// MarshalResultChunkData method of the parent MockSerializer instance is
// invoked and the hook queue is empty.
func (f *SerializerMarshalResultChunkDataFunc) SetDefaultHook(hook func(types.ResultChunkData) ([]byte, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// MarshalResultChunkData method of the parent MockSerializer instance
// inovkes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *SerializerMarshalResultChunkDataFunc) PushHook(hook func(types.ResultChunkData) ([]byte, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SerializerMarshalResultChunkDataFunc) SetDefaultReturn(r0 []byte, r1 error) {
	f.SetDefaultHook(func(types.ResultChunkData) ([]byte, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SerializerMarshalResultChunkDataFunc) PushReturn(r0 []byte, r1 error) {
	f.PushHook(func(types.ResultChunkData) ([]byte, error) {
		return r0, r1
	})
}

func (f *SerializerMarshalResultChunkDataFunc) nextHook() func(types.ResultChunkData) ([]byte, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SerializerMarshalResultChunkDataFunc) appendCall(r0 SerializerMarshalResultChunkDataFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SerializerMarshalResultChunkDataFuncCall
// objects describing the invocations of this function.
func (f *SerializerMarshalResultChunkDataFunc) History() []SerializerMarshalResultChunkDataFuncCall {
	f.mutex.Lock()
	history := make([]SerializerMarshalResultChunkDataFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SerializerMarshalResultChunkDataFuncCall is an object that describes an
// invocation of method MarshalResultChunkData on an instance of
// MockSerializer.
type SerializerMarshalResultChunkDataFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 types.ResultChunkData
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []byte
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SerializerMarshalResultChunkDataFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SerializerMarshalResultChunkDataFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// SerializerUnmarshalDocumentDataFunc describes the behavior when the
// UnmarshalDocumentData method of the parent MockSerializer instance is
// invoked.
type SerializerUnmarshalDocumentDataFunc struct {
	defaultHook func([]byte) (types.DocumentData, error)
	hooks       []func([]byte) (types.DocumentData, error)
	history     []SerializerUnmarshalDocumentDataFuncCall
	mutex       sync.Mutex
}

// UnmarshalDocumentData delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockSerializer) UnmarshalDocumentData(v0 []byte) (types.DocumentData, error) {
	r0, r1 := m.UnmarshalDocumentDataFunc.nextHook()(v0)
	m.UnmarshalDocumentDataFunc.appendCall(SerializerUnmarshalDocumentDataFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// UnmarshalDocumentData method of the parent MockSerializer instance is
// invoked and the hook queue is empty.
func (f *SerializerUnmarshalDocumentDataFunc) SetDefaultHook(hook func([]byte) (types.DocumentData, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UnmarshalDocumentData method of the parent MockSerializer instance
// inovkes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *SerializerUnmarshalDocumentDataFunc) PushHook(hook func([]byte) (types.DocumentData, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SerializerUnmarshalDocumentDataFunc) SetDefaultReturn(r0 types.DocumentData, r1 error) {
	f.SetDefaultHook(func([]byte) (types.DocumentData, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SerializerUnmarshalDocumentDataFunc) PushReturn(r0 types.DocumentData, r1 error) {
	f.PushHook(func([]byte) (types.DocumentData, error) {
		return r0, r1
	})
}

func (f *SerializerUnmarshalDocumentDataFunc) nextHook() func([]byte) (types.DocumentData, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SerializerUnmarshalDocumentDataFunc) appendCall(r0 SerializerUnmarshalDocumentDataFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SerializerUnmarshalDocumentDataFuncCall
// objects describing the invocations of this function.
func (f *SerializerUnmarshalDocumentDataFunc) History() []SerializerUnmarshalDocumentDataFuncCall {
	f.mutex.Lock()
	history := make([]SerializerUnmarshalDocumentDataFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SerializerUnmarshalDocumentDataFuncCall is an object that describes an
// invocation of method UnmarshalDocumentData on an instance of
// MockSerializer.
type SerializerUnmarshalDocumentDataFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 []byte
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 types.DocumentData
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SerializerUnmarshalDocumentDataFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SerializerUnmarshalDocumentDataFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// SerializerUnmarshalLocationsFunc describes the behavior when the
// UnmarshalLocations method of the parent MockSerializer instance is
// invoked.
type SerializerUnmarshalLocationsFunc struct {
	defaultHook func([]byte) ([]types.Location, error)
	hooks       []func([]byte) ([]types.Location, error)
	history     []SerializerUnmarshalLocationsFuncCall
	mutex       sync.Mutex
}

// UnmarshalLocations delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockSerializer) UnmarshalLocations(v0 []byte) ([]types.Location, error) {
	r0, r1 := m.UnmarshalLocationsFunc.nextHook()(v0)
	m.UnmarshalLocationsFunc.appendCall(SerializerUnmarshalLocationsFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the UnmarshalLocations
// method of the parent MockSerializer instance is invoked and the hook
// queue is empty.
func (f *SerializerUnmarshalLocationsFunc) SetDefaultHook(hook func([]byte) ([]types.Location, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UnmarshalLocations method of the parent MockSerializer instance inovkes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *SerializerUnmarshalLocationsFunc) PushHook(hook func([]byte) ([]types.Location, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SerializerUnmarshalLocationsFunc) SetDefaultReturn(r0 []types.Location, r1 error) {
	f.SetDefaultHook(func([]byte) ([]types.Location, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SerializerUnmarshalLocationsFunc) PushReturn(r0 []types.Location, r1 error) {
	f.PushHook(func([]byte) ([]types.Location, error) {
		return r0, r1
	})
}

func (f *SerializerUnmarshalLocationsFunc) nextHook() func([]byte) ([]types.Location, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SerializerUnmarshalLocationsFunc) appendCall(r0 SerializerUnmarshalLocationsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SerializerUnmarshalLocationsFuncCall
// objects describing the invocations of this function.
func (f *SerializerUnmarshalLocationsFunc) History() []SerializerUnmarshalLocationsFuncCall {
	f.mutex.Lock()
	history := make([]SerializerUnmarshalLocationsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SerializerUnmarshalLocationsFuncCall is an object that describes an
// invocation of method UnmarshalLocations on an instance of MockSerializer.
type SerializerUnmarshalLocationsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 []byte
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []types.Location
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SerializerUnmarshalLocationsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SerializerUnmarshalLocationsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// SerializerUnmarshalResultChunkDataFunc describes the behavior when the
// UnmarshalResultChunkData method of the parent MockSerializer instance is
// invoked.
type SerializerUnmarshalResultChunkDataFunc struct {
	defaultHook func([]byte) (types.ResultChunkData, error)
	hooks       []func([]byte) (types.ResultChunkData, error)
	history     []SerializerUnmarshalResultChunkDataFuncCall
	mutex       sync.Mutex
}

// UnmarshalResultChunkData delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockSerializer) UnmarshalResultChunkData(v0 []byte) (types.ResultChunkData, error) {
	r0, r1 := m.UnmarshalResultChunkDataFunc.nextHook()(v0)
	m.UnmarshalResultChunkDataFunc.appendCall(SerializerUnmarshalResultChunkDataFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// UnmarshalResultChunkData method of the parent MockSerializer instance is
// invoked and the hook queue is empty.
func (f *SerializerUnmarshalResultChunkDataFunc) SetDefaultHook(hook func([]byte) (types.ResultChunkData, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UnmarshalResultChunkData method of the parent MockSerializer instance
// inovkes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *SerializerUnmarshalResultChunkDataFunc) PushHook(hook func([]byte) (types.ResultChunkData, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SerializerUnmarshalResultChunkDataFunc) SetDefaultReturn(r0 types.ResultChunkData, r1 error) {
	f.SetDefaultHook(func([]byte) (types.ResultChunkData, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SerializerUnmarshalResultChunkDataFunc) PushReturn(r0 types.ResultChunkData, r1 error) {
	f.PushHook(func([]byte) (types.ResultChunkData, error) {
		return r0, r1
	})
}

func (f *SerializerUnmarshalResultChunkDataFunc) nextHook() func([]byte) (types.ResultChunkData, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SerializerUnmarshalResultChunkDataFunc) appendCall(r0 SerializerUnmarshalResultChunkDataFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SerializerUnmarshalResultChunkDataFuncCall
// objects describing the invocations of this function.
func (f *SerializerUnmarshalResultChunkDataFunc) History() []SerializerUnmarshalResultChunkDataFuncCall {
	f.mutex.Lock()
	history := make([]SerializerUnmarshalResultChunkDataFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SerializerUnmarshalResultChunkDataFuncCall is an object that describes an
// invocation of method UnmarshalResultChunkData on an instance of
// MockSerializer.
type SerializerUnmarshalResultChunkDataFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 []byte
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 types.ResultChunkData
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SerializerUnmarshalResultChunkDataFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SerializerUnmarshalResultChunkDataFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
