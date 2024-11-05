package user

import (
	"errors"
	"fmt"
	"encoding/base64"
	"testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRep struct {
	mock.Mock
}

// FindByID
func (m *MockRep) FindByID(id uuid.UUID) (User, error) {
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}

// Create
func (m *MockRep) Create(user User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindAll.
func (m *MockRep) FindAll() ([]User, error) {
	args := m.Called()
	return args.Get(0).([]User), args.Error(1)
}

// FindByEmail
func (m *MockRep) FindByEmail(email string) (User, error) {
	args := m.Called(email)
	return args.Get(0).(User), args.Error(1)
}


// Get fetches a user by id.
func TestServiceGet(t *testing.T) {

	testID := uuid.New()
	mockRepo := new(MockRep)
	svc := NewService(mockRepo)

	expectUser := User{ID: testID, Name: "Test"}
	mockRepo.On("FindByID", testID).Return(expectUser, nil)

	user, err := svc.Get(testID)

	assert.NoError(t, err)
	assert.Equal(t, expectUser, user)
	mockRepo.AssertCalled(t, "FindByID", testID)

}

// Get fetches a user by id.
func TestService_Get_NotFound(t *testing.T) {
	testID := uuid.New()
	mockRepo := new(MockRep)
	svc := NewService(mockRepo)

	mockRepo.On("FindByID", testID).Return(User{}, errors.New("User not found"))

	user, err := svc.Get(testID)

	assert.Error(t, err)
	assert.Equal(t, User{}, user)

	mockRepo.AssertCalled(t, "FindByID", testID)
}

// List fetches all users.
func TestService_FindAll(t *testing.T) {

	mockRepo := new(MockRep)
	svc := NewService(mockRepo)
	myId := uuid.New()

	testUsers := []User{
		{ID: myId, Name: "Testing"},
		{ID: myId, Name: "Testing 2"},
	}
	mockRepo.On("FindAll").Return(testUsers, nil)
	users, err := svc.List()

	assert.NoError(t, err)
	assert.Equal(t, testUsers, users)

	mockRepo.AssertCalled(t, "FindAll")

}



func TestService_Signin_All(t *testing.T) {
	mockRepo := new(MockRep)
	svc := NewService(mockRepo)

	cases := []struct {
		testName          string
		email             string
		password          string
		expectedUser      User
		expectedError     error
		repoExpectedEmail string
		repoOutUser       User
		repoOutError      error
	}{
		{
			testName: "ok",
			email:    "test@test.com",
			password: "testpas123",
			expectedUser: User{
				Name:     "Stas",
				Email:    "test@test.com",
				Password: "testpas123",
			},
			expectedError:     nil,
			repoExpectedEmail: "test@test.com",
			repoOutUser: User{
				Name:     "Stas",
				Email:    "test@test.com",
				Password: "testpas123",
			},
			repoOutError: nil,
		}, 
		{
			testName:          "userNotFound",
			email:             "testUserNotFound@test.com",
			password:          "wrongpasword",
			expectedUser:      User{},
			expectedError:     ErrNotFound,
			repoExpectedEmail: "testUserNotFound@test.com",
			repoOutUser:       User{},
			repoOutError:      ErrNotFound,
		},

		{
			testName: "wrongPassword",
			email:    "testUserPaasword@test.com",
			password: "wrongpasword",
			expectedUser: User{},
			expectedError:     ErrNotFound,
			repoExpectedEmail: "testUserPaasword@test.com",
			repoOutUser: User{
				Name:     "Stas",
				Email:    "test@test.com",
				Password: "testpas123",
			},
			repoOutError: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {			
			mockRepo.On("FindByEmail", tc.repoExpectedEmail).Return(tc.repoOutUser, tc.repoOutError)
			user, err := svc.Signin(tc.email, tc.password)
			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedError, err)
		})
	}

}

func TestService_Signup_All(t *testing.T){
	

	cases := [] struct{
		testName  string
		email string
		password string
		name string
		expectedUser User
		expectedError error
		repoExpectedEmail string
		repoOutUser User
		repoOutError error
	}{
		{
			testName: "findEmail",
			email: "test12345@test.com",
			password: "pass123",
			name: "Test",
			expectedUser: User{},
			expectedError: ErrEmailExists,
			repoExpectedEmail: "test12345@test.com",
			repoOutUser: User{},
			repoOutError: nil,
		},
		{
			testName: "createUser",
			email: "test12345@test.com",
			name: "Test",
			password: "pass123",
			expectedUser: User{
				Email: "test12345@test.com",
				Name: "Test",
				Password: "pass123",
			},
			expectedError: nil,
			repoExpectedEmail: "test12345@test.com",
			repoOutUser: User{},
			repoOutError: ErrNotFound,
				
		},
		{
			testName: "createUser",
			email: "test12345@test.com",
			name: "Test",
			password: "pass123",
			expectedUser: User{
				Email: "test12345@test.com",
				Name: "Test",
				Password: "pass123",
			},
			expectedError: nil,
			repoExpectedEmail: "test12345@test.com",
			repoOutUser: User{},
			repoOutError: ErrNotFound,
				
		},
		{
			testName: "errorEmail",
			email: "test123456@test.com",
			name: "Test123",
			password: "pa123",
			
			expectedUser: User{
				Email: "test123456@test.com",
				Name: "Test123",
				Password: "pass123",
			},
			expectedError: fmt.Errorf("error get email:%w", errors.New("random error")),
			repoExpectedEmail: "test123456@test.com",
			repoOutUser: User{},
			repoOutError: errors.New("random error"),
		},

		

	}
	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			mockRepo := new(MockRep)
			svc := NewService(mockRepo)
			mockRepo.On("FindByEmail", tc.repoExpectedEmail).Return(tc.repoOutUser, tc.repoOutError)
	
			if errors.Is(tc.repoOutError, ErrNotFound) {
				mockRepo.On("Create", mock.MatchedBy(func(u User) bool{
					return u.Email == tc.expectedUser.Email && 
					u.Name == tc.expectedUser.Name
				})).Return(tc.expectedError)
			}
	
			err := svc.Signup(tc.email, tc.name, tc.password)
			t.Logf("Expected error: %v, Actual error: %v", tc.expectedError, err)
	
			if tc.expectedError != nil {
				assert.ErrorContains(t, err, tc.expectedError.Error()) 
			} else {
				assert.NoError(t, err)
			}
		})
	}
}




func TestServicee_checkPasswordHash(t *testing.T){
	egz := &service{}

	//пароль верный и hash
	password := "password"
	encodedHash, err := egz.passHashed(password)
	require.NoError(t,err)
	


	// пароль совпадает с hash
	isVal, err := egz.checkPasswordHash(password,encodedHash)
	assert.NoError(t,err)
	assert.True(t, isVal)

	// неверный пароль
	wrongPas := "wrong"
	isVal,err = egz.checkPasswordHash(wrongPas,encodedHash)
	assert.NoError(t,err)
	assert.False(t,isVal)

	// неверный формат hash
	invalHash := "this_is_not_valid_base64!"

	_, err = egz.checkPasswordHash(password,invalHash)
	assert.Error(t,err)


	//  Валидный формат, но невалидная Base64 строка для соли
	invalidBase64Salt := "invalidbase64$" + base64.RawStdEncoding.EncodeToString([]byte("valid_hash"))
	_, err = egz.checkPasswordHash(password, invalidBase64Salt)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode salt")

	//  Валидная соль, но невалидная Base64 строка для хеша
	invalidBase64Hash := base64.RawStdEncoding.EncodeToString([]byte("valid_salt")) + "$invalidbase64"
	_, err = egz.checkPasswordHash(password, invalidBase64Hash)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode hash")
	
	
}
	

