package user

import (
	"errors"
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


type MockService struct{
	 funcSignup func(email,name,password string)(error)
	 funcList func()([]User,error)
	 funcSignin func(email,password string)(User,error)
	 funcGet func(id uuid.UUID)(User,error)

}



// List
func (m *MockService) List()([]User,error){
	if m.funcList  !=nil{
		return m.funcList()	
	}
	return []User{}, nil	
}

// Get
func(m *MockService) Get(id uuid.UUID)(User,error){
	return m.funcGet(id)
}

// Signin
func(m *MockService) Signin(email,password string)(User,error){
	return m.funcSignin(email,password)
}

// Signup
func(m *MockService) Signup(email,name,password string)(error){
	if m.funcSignup != nil{
		return m.funcSignup(email,name,password)
	}
	return nil
}




// Signup 
func TestHandlerSignup(t *testing.T){
	cases := []struct{
		testName string
		requestBody string
		funcSignup func(email,name,password string) error
		expectedStatus int
	}{
		{
			testName: "Successful signup",
			requestBody: "email=test@test.com&name=testName&password=test123",
			funcSignup: func(email,name,password string) error{
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			testName: "Signup error",
			requestBody: "email=test1@test.com&name=TestName1&password=password1234",
			funcSignup: func(email,name,password string)error{
				return errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		
		
	}

	for _, tc := range cases{
		t.Run(tc.testName, func(t *testing.T) {
			req, err:= http.NewRequest("POST","/signup", strings.NewReader(tc.requestBody))
			if err != nil{
				t.Fatal(err)
			}
			req.Header.Set("content-type","appication/json")

			
			rr := httptest.NewRecorder()

			
			mockService := &MockService{
				funcSignup: tc.funcSignup,
			}

			
			handler := &Handler{service: mockService}

			handler.Signup(rr,req)

			if rr.Code != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatus)
            }


		})
	}
}


// List
func TestList(t *testing.T){
	cases :=[]struct{
		testName string
		funcList func()([]User,error)
		expectedStatus  int
	}{
		{
			testName: "Succesfull list",
			funcList: func()([]User,error){
				      
				return []User{},nil
				
			},
				
			expectedStatus: http.StatusOK,
		},
		{
			testName: "List error",
			funcList: func()([]User,error){
					   
				return nil,errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc :=range(cases){
		t.Run(tc.testName,func(t *testing.T) {
			req, err:= http.NewRequest("GET","/list",strings.NewReader(""))

			if err != nil{
				t.Fatal(err)
			}
			req.Header.Set("content-type","application/json")

			rr := httptest.NewRecorder()

			mockService := MockService{
				funcList: tc.funcList,
			}
											
			handler := &Handler{service: &mockService}
			handler.List(rr,req)

			if rr.Code != tc.expectedStatus{
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			})
	}
}

// Signin
func TestSignin(t *testing.T){
	cases :=[]struct{
		testName string
		requestBody string
		funcSignin func(email,password string)(User,error)
		expectedStatus int
	}{
		{
			testName: "Successfull Signin",
			requestBody: "email=teststas@ex.com&password=123test",
			funcSignin: func(email,password string) (User,error){
						
				return User{Email: "teststas@ex.com", Name: "test"}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "Unauthorized Signin",
			requestBody: "email=123@email.com&password=1231231231",
			funcSignin: func(email,password string)(User,error){
						
				return User{},ErrInvalidCredentials
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			testName: "Internal Server error Signin",
			requestBody: "email=123@email.com&password=1231231231",
			funcSignin: func(email,password string)(User,error){
				return User{},errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			testName: "Email and password are required Signin",
			requestBody: "email=&password=",
			funcSignin: func (email string, password string) (User, error){
				return User{}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range(cases){
		t.Run(tc.testName,func(t *testing.T) {
			req,err := http.NewRequest("POST","/signin",strings.NewReader(tc.requestBody))

			if err !=nil{
				t.Fatal(err)
			}
											
			req.Header.Set("content-type","application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			mockService := MockService{
				funcSignin: tc.funcSignin,
			}

			handler := &Handler{service: &mockService}

			handler.Signin(rr,req)
			

			if rr.Code != tc.expectedStatus{
				t.Errorf("expected status %d, got %d",tc.expectedStatus,rr.Code)
			}

		})

	}
}

// Get
func TestGet(t *testing.T){
	cases := []struct{
		nameTest string
		requestUUID string
		funcGet func(id uuid.UUID)(User, error)
		expectedStatus int
	}{
		{
			nameTest: "Success Get",
						
			requestUUID: "123e4567-e89b-12d3-a456-426614174000",
			funcGet: func(id uuid.UUID)(User,error){
						
				return User{}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			nameTest: "User UUID is required Get",
			requestUUID: "",
			funcGet: func(id uuid.UUID) (User, error){
				return User{}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			nameTest: "Invalid UUID format Get",
			requestUUID: "inval uuid form",
			funcGet: func(id uuid.UUID)(User, error){
				return User{}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			nameTest: "User not found Get",
			requestUUID: "123e4567-e89b-12d3-a456-426614174000",
			funcGet: func(id uuid.UUID)(User,error){
				return User{}, ErrUserNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			nameTest: "Internal server error Get",
			requestUUID: "123e4567-e89b-12d3-a456-426614174000",
			funcGet: func(id uuid.UUID)(User,error){
				return User{}, errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		
	}

	for _, tc := range(cases){
		t.Run(tc.nameTest,func(t *testing.T) {
										
			req, err := http.NewRequest("GET","/get?uuid="+tc.requestUUID,nil)

			if err != nil{
				t.Fatal(err)
			}
											
			req.Header.Set("content-type","application/json")

			rr := httptest.NewRecorder()

			mockService := MockService{funcGet: tc.funcGet}

			handler := &Handler{service: &mockService}

			handler.Get(rr,req)

			if rr.Code != tc.expectedStatus{
				t.Errorf("expected status %d, got %d", tc.expectedStatus,rr.Code)
			}
		})

	} 
}

// RegisterUserRouter

func TestRegisterUserRout(t *testing.T){
	
	
	cases :=[]struct{
		nameTest string
		nameRouts string
		nameMethod string
		expectedStatus int

	}{
		{
			nameTest: "Signup Handler",
			nameMethod: http.MethodPost,
			nameRouts: "/user/signup",
			expectedStatus: http.StatusCreated,
		},
		{
			nameTest: "Path Handler",
			nameMethod: http.MethodGet,
			nameRouts:  "/user/list",
			expectedStatus: http.StatusOK,
		},
	}
	
	for _,tc:=range(cases){
		t.Run(tc.nameTest,func(t *testing.T) {
			
			router := chi.NewRouter()
							
			mockService := &MockService{}

			handler :=&Handler{
						service: mockService,
					}

			handler.RegisterUserRouter(router)
			req, err :=http.NewRequest(tc.nameMethod,tc.nameRouts,nil)
			if err !=nil{
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			
			router.ServeHTTP(rr,req)
			
			if rr.Code != tc.expectedStatus {
                t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
            }

		})
	}
}