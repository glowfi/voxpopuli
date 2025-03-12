// Code generated by counterfeiter. DO NOT EDIT.
package postfakes

import (
	"context"
	"sync"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/glowfi/voxpopuli/backend/pkg/transport/post"
)

type FakePostService struct {
	PostsPaginatedStub        func(context.Context, int, int) ([]models.Post, error)
	postsPaginatedMutex       sync.RWMutex
	postsPaginatedArgsForCall []struct {
		arg1 context.Context
		arg2 int
		arg3 int
	}
	postsPaginatedReturns struct {
		result1 []models.Post
		result2 error
	}
	postsPaginatedReturnsOnCall map[int]struct {
		result1 []models.Post
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePostService) PostsPaginated(arg1 context.Context, arg2 int, arg3 int) ([]models.Post, error) {
	fake.postsPaginatedMutex.Lock()
	ret, specificReturn := fake.postsPaginatedReturnsOnCall[len(fake.postsPaginatedArgsForCall)]
	fake.postsPaginatedArgsForCall = append(fake.postsPaginatedArgsForCall, struct {
		arg1 context.Context
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.PostsPaginatedStub
	fakeReturns := fake.postsPaginatedReturns
	fake.recordInvocation("PostsPaginated", []interface{}{arg1, arg2, arg3})
	fake.postsPaginatedMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePostService) PostsPaginatedCallCount() int {
	fake.postsPaginatedMutex.RLock()
	defer fake.postsPaginatedMutex.RUnlock()
	return len(fake.postsPaginatedArgsForCall)
}

func (fake *FakePostService) PostsPaginatedCalls(stub func(context.Context, int, int) ([]models.Post, error)) {
	fake.postsPaginatedMutex.Lock()
	defer fake.postsPaginatedMutex.Unlock()
	fake.PostsPaginatedStub = stub
}

func (fake *FakePostService) PostsPaginatedArgsForCall(i int) (context.Context, int, int) {
	fake.postsPaginatedMutex.RLock()
	defer fake.postsPaginatedMutex.RUnlock()
	argsForCall := fake.postsPaginatedArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePostService) PostsPaginatedReturns(result1 []models.Post, result2 error) {
	fake.postsPaginatedMutex.Lock()
	defer fake.postsPaginatedMutex.Unlock()
	fake.PostsPaginatedStub = nil
	fake.postsPaginatedReturns = struct {
		result1 []models.Post
		result2 error
	}{result1, result2}
}

func (fake *FakePostService) PostsPaginatedReturnsOnCall(i int, result1 []models.Post, result2 error) {
	fake.postsPaginatedMutex.Lock()
	defer fake.postsPaginatedMutex.Unlock()
	fake.PostsPaginatedStub = nil
	if fake.postsPaginatedReturnsOnCall == nil {
		fake.postsPaginatedReturnsOnCall = make(map[int]struct {
			result1 []models.Post
			result2 error
		})
	}
	fake.postsPaginatedReturnsOnCall[i] = struct {
		result1 []models.Post
		result2 error
	}{result1, result2}
}

func (fake *FakePostService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.postsPaginatedMutex.RLock()
	defer fake.postsPaginatedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePostService) recordInvocation(key string, args []interface{}) {
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

var _ post.PostService = new(FakePostService)
