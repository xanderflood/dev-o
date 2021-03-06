// Code generated by counterfeiter. DO NOT EDIT.
package devofakes

import (
	"sync"

	devo "github.com/xanderflood/dev-o"
)

type FakeBinner struct {
	BuildStub        func() error
	buildMutex       sync.RWMutex
	buildArgsForCall []struct{}
	buildReturns     struct {
		result1 error
	}
	buildReturnsOnCall map[int]struct {
		result1 error
	}
	ExecStub         func()
	execMutex        sync.RWMutex
	execArgsForCall  []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBinner) Build() error {
	fake.buildMutex.Lock()
	ret, specificReturn := fake.buildReturnsOnCall[len(fake.buildArgsForCall)]
	fake.buildArgsForCall = append(fake.buildArgsForCall, struct{}{})
	fake.recordInvocation("Build", []interface{}{})
	fake.buildMutex.Unlock()
	if fake.BuildStub != nil {
		return fake.BuildStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildReturns.result1
}

func (fake *FakeBinner) BuildCallCount() int {
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	return len(fake.buildArgsForCall)
}

func (fake *FakeBinner) BuildReturns(result1 error) {
	fake.BuildStub = nil
	fake.buildReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBinner) BuildReturnsOnCall(i int, result1 error) {
	fake.BuildStub = nil
	if fake.buildReturnsOnCall == nil {
		fake.buildReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.buildReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBinner) Exec() {
	fake.execMutex.Lock()
	fake.execArgsForCall = append(fake.execArgsForCall, struct{}{})
	fake.recordInvocation("Exec", []interface{}{})
	fake.execMutex.Unlock()
	if fake.ExecStub != nil {
		fake.ExecStub()
	}
}

func (fake *FakeBinner) ExecCallCount() int {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return len(fake.execArgsForCall)
}

func (fake *FakeBinner) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildMutex.RLock()
	defer fake.buildMutex.RUnlock()
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBinner) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ devo.Binner = new(FakeBinner)
