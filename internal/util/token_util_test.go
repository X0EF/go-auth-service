package util

import (
	"testing"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/helper"
)

func TestVerifyAndExtract(t *testing.T) {
	secret := "secret"
	user := domain.User{Username: "onion", ID: 1}
	tokenString, err := CreateAccessToken(&user, secret, 1)
	if err != nil {
		t.Errorf("error creating token: %v", err.Error())
		return
	}

	claims, err := VerifyAndExtract(tokenString, secret)
	if err != nil {
		t.Errorf("VerifyAndExtract returned an error: %v", err.Error())
		return
	}

	// Check if the claims were extracted successfully
	if claims == nil {
		t.Error("VerifyAndExtract returned nil claims")
		return
	}
}

func TestVerifyAndExtractExpired(t *testing.T) {
	secret := "secret"
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFiMjIzIiwiaWQiOjMsImV4cCI6MTcwNzgwMjMxOH0.Nz1zRUGVlfE4Ac_k_SyYTYkxn63TfI5ClF_cW5XNaJo"

	_, err := VerifyAndExtract(tokenString, secret)
	if err.Error() != helper.ErrExpiredAuthToken {
		t.Error(err)
		t.Error("VerifyAndExtract did not flag expired token")
		return
	}
}
