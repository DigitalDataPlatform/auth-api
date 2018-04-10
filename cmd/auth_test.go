package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	ddpAuth "gitlab.adeo.com/ddp-auth/cmd"
)

func Test_GetLogin_EmptyUsername(t *testing.T) {
	req, err := http.NewRequest("GET", "/login?password=password", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ddpAuth.LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusPreconditionFailed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPreconditionFailed)
	}

}

func Test_GetLogin_EmptyPassword(t *testing.T) {
	req, err := http.NewRequest("GET", "/login?username=username", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ddpAuth.LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusPreconditionFailed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPreconditionFailed)
	}

}

func Test_GetLogin_ReturnValidToken(t *testing.T) {

	req, err := http.NewRequest("GET", "/login?username=test&password=password", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ddpAuth.LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	token := rr.Body.String()

	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("monsupersecretvraimentsupersecret"), nil
	})

	if err != nil || !tk.Valid {
		t.Fatal(err)
	}

	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		if claims["iss"] != "digital-data-platform" {
			t.Errorf("Claims %v should be %v", "iss", "digital-data-platform")
		}
		if claims["sub"] != "digital-data-platform" {
			t.Errorf("Claims %v should be %v", "sub", "digital-data-platform")
		}
		if claims["aud"] != "digital-data-platform" {
			t.Errorf("Claims %v should be %v", "aud", "digital-data-platform")
		}
		if claims["user"] != "20009473" {
			t.Errorf("Claims %v should be %v", "user", "20009473")
		}
		if claims["username"] != "fclaeys" {
			t.Errorf("Claims %v should be %v", "username", "fclaeys")
		}
	}
}
