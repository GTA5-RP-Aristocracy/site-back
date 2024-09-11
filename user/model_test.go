package user

import(
	"testing"
	"time"

	
)

func TestUserField(t *testing.T){
	
	testuser := User{
		Name: "test",
		Email: "test@email.com",
		Password: "123pass",
		Created: time.Now(),
		Updated: time.Now(),

	}

	if testuser.Name =="" || testuser.Email == "" || testuser.Password ==""{
		t.Errorf("Пустая строка")
	}
	if testuser.Created.IsZero() || testuser.Updated.IsZero(){
		t.Errorf("Время не может быть нулем")
	}
}