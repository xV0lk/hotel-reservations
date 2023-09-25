package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.Drop(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	for _, tc := range userTests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testUser(t, app)
		})
	}

}

func (tc testCase[userTest]) testUser(t *testing.T, app *fiber.App) {
	b, _ := json.Marshal(tc.input)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != tc.status {
		t.Fatalf("expected status %d, got %d", tc.status, res.StatusCode)
	}
	if tc.ttype == "pass" {
		var resBody any
		err := json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			t.Fatal(err)
		}
		if resBody.(map[string]interface{})["password"] != nil {
			t.Error("expected password to be nil")
		}
		for k, v := range tc.body.(map[string]string) {
			if resBody.(map[string]interface{})[k] != v {
				t.Errorf("expected %s to be %s, got %s", k, v, resBody.(map[string]interface{})[k])
			}
		}
	} else if tc.ttype == "fail" {
		var resBody failUserResponse
		err := json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range tc.body.(map[string]string) {
			if resBody.Error[k] != v {
				t.Errorf("expected %s to be %s, got %s", k, v, resBody.Error[k])
			}
		}
	}
}

type userTest types.NewUserParams

var userTests = []testCase[userTest]{
	{
		name:  "valid user",
		ttype: "pass",
		input: userTest{
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser@test.com",
			Password:  "Testp@ssword123",
		},
		expected: expected{
			status: fiber.StatusCreated,
			body: map[string]string{
				"firstName": "Test",
				"lastName":  "User",
				"email":     "testuser@test.com",
			},
		},
	},
	{
		name:  "invalid email and password",
		ttype: "fail",
		input: userTest{
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser",
			Password:  "testpassword123",
		},
		expected: expected{
			status: fiber.StatusBadRequest,
			body: map[string]string{
				"email":    "invalid email address",
				"password": "password must contain at least one uppercase letter, password must contain at least one special character",
			},
		},
	},
}
