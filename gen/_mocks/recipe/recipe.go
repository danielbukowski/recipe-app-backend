// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/recipe/handlers.go
//
// Generated by this command:
//
//	mockgen -source=./internal/recipe/handlers.go -destination=./gen/_mocks/recipe/recipe.go -mock_names=cacheStorage=MockCacheStorage,recipeService=MockRecipeService
//

// Package mock_recipe is a generated GoMock package.
package mock_recipe

import (
	context "context"
	reflect "reflect"
	time "time"

	recipe "github.com/danielbukowski/recipe-app-backend/internal/recipe"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRecipeService is a mock of recipeService interface.
type MockRecipeService struct {
	ctrl     *gomock.Controller
	recorder *MockRecipeServiceMockRecorder
	isgomock struct{}
}

// MockRecipeServiceMockRecorder is the mock recorder for MockRecipeService.
type MockRecipeServiceMockRecorder struct {
	mock *MockRecipeService
}

// NewMockRecipeService creates a new mock instance.
func NewMockRecipeService(ctrl *gomock.Controller) *MockRecipeService {
	mock := &MockRecipeService{ctrl: ctrl}
	mock.recorder = &MockRecipeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecipeService) EXPECT() *MockRecipeServiceMockRecorder {
	return m.recorder
}

// CreateNewRecipe mocks base method.
func (m *MockRecipeService) CreateNewRecipe(arg0 context.Context, arg1 recipe.NewRecipeRequest) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewRecipe", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewRecipe indicates an expected call of CreateNewRecipe.
func (mr *MockRecipeServiceMockRecorder) CreateNewRecipe(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewRecipe", reflect.TypeOf((*MockRecipeService)(nil).CreateNewRecipe), arg0, arg1)
}

// DeleteRecipeById mocks base method.
func (m *MockRecipeService) DeleteRecipeById(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRecipeById", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRecipeById indicates an expected call of DeleteRecipeById.
func (mr *MockRecipeServiceMockRecorder) DeleteRecipeById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRecipeById", reflect.TypeOf((*MockRecipeService)(nil).DeleteRecipeById), arg0, arg1)
}

// GetRecipeById mocks base method.
func (m *MockRecipeService) GetRecipeById(arg0 context.Context, arg1 uuid.UUID) (recipe.RecipeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecipeById", arg0, arg1)
	ret0, _ := ret[0].(recipe.RecipeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecipeById indicates an expected call of GetRecipeById.
func (mr *MockRecipeServiceMockRecorder) GetRecipeById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecipeById", reflect.TypeOf((*MockRecipeService)(nil).GetRecipeById), arg0, arg1)
}

// UpdateRecipeById mocks base method.
func (m *MockRecipeService) UpdateRecipeById(arg0 context.Context, arg1 uuid.UUID, arg2 time.Time, arg3 recipe.UpdateRecipeRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecipeById", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecipeById indicates an expected call of UpdateRecipeById.
func (mr *MockRecipeServiceMockRecorder) UpdateRecipeById(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecipeById", reflect.TypeOf((*MockRecipeService)(nil).UpdateRecipeById), arg0, arg1, arg2, arg3)
}

// MockCacheStorage is a mock of cacheStorage interface.
type MockCacheStorage struct {
	ctrl     *gomock.Controller
	recorder *MockCacheStorageMockRecorder
	isgomock struct{}
}

// MockCacheStorageMockRecorder is the mock recorder for MockCacheStorage.
type MockCacheStorageMockRecorder struct {
	mock *MockCacheStorage
}

// NewMockCacheStorage creates a new mock instance.
func NewMockCacheStorage(ctrl *gomock.Controller) *MockCacheStorage {
	mock := &MockCacheStorage{ctrl: ctrl}
	mock.recorder = &MockCacheStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheStorage) EXPECT() *MockCacheStorageMockRecorder {
	return m.recorder
}

// DeleteItem mocks base method.
func (m *MockCacheStorage) DeleteItem(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockCacheStorageMockRecorder) DeleteItem(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockCacheStorage)(nil).DeleteItem), key)
}

// GetItem mocks base method.
func (m *MockCacheStorage) GetItem(key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockCacheStorageMockRecorder) GetItem(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockCacheStorage)(nil).GetItem), key)
}

// InsertItem mocks base method.
func (m *MockCacheStorage) InsertItem(key string, value []byte, expiration int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertItem", key, value, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertItem indicates an expected call of InsertItem.
func (mr *MockCacheStorageMockRecorder) InsertItem(key, value, expiration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertItem", reflect.TypeOf((*MockCacheStorage)(nil).InsertItem), key, value, expiration)
}
