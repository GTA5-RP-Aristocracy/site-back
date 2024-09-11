package user

import (
	// "container/list"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)



   
type MockRep struct{
	mock.Mock
}
// Get fetches a user by id. надо обратиться id по int
func (m *MockRep) FindByID(id uuid.UUID) (User, error){
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}

// List fetches all users.
// func (m *MockRep) FindAll(List()) (User, error){
// 	args := m.Called(List)
// 	return args.Findall(List).(User), args.Error(1)
// }
//птом с денисом еще оьбьясним
type testingMok struct {
	repo *MockRep
}
// Get fetches a user by id  
func (mk *testingMok) Get(id uuid.UUID) (User, error){
	return mk.repo.FindByID(id)
}

// List fetches all users.
// func (mk *testingMok) FindAll(List) (User, error){
// 	return mk.repo.FindAll(List)
// }


// Get fetches a user by id.
func TestServiceGet(t *testing.T){
	
	testUUID := uuid.New()
	// new(MockRep) -> так уже не пишут, пишут MockRep{}
	mockRepo := new(MockRep)
	svc := &testingMok{repo: mockRepo} 

	expectUser := User{ID: testUUID, Name: "Test"}
	mockRepo.On("FindByID",testUUID).Return(expectUser, nil)
	 
	user, err := svc.Get(testUUID)

	assert.NoError(t,err)
	assert.Equal(t,expectUser,user)
	mockRepo.AssertCalled(t,"FindByID",testUUID)


}

// Get fetches a user by id.
func TestService_Get_NotFound(t *testing.T){
	testUUID := uuid.New()
	mockRepo := new(MockRep)
	svc := &testingMok{repo: mockRepo}

	mockRepo.On("FindByID",testUUID).Return(User{}, errors.New("User not found"))
	
	user, err := svc.Get(testUUID)

	assert.Error(t,err)
	assert.Equal(t, User{}, user)

	mockRepo.AssertCalled(t, "FindByID",testUUID)
}


// List fetches all users.