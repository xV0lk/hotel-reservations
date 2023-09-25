package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/db/fixtures"
	"github.com/xV0lk/hotel-reservations/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.Drop(t)

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	userHandler := NewUserHandler(tdb.Store.User)
	api := app.Group("/", JWTAuth(tdb.Store.User))
	api.Post("/", userHandler.HandlePostUser)

	for _, tc := range userTests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testUser(t, app, tdb.Store)
		})
	}

}

func (tc testCase[userTest]) testUser(t *testing.T, app *fiber.App, store *db.Store) {
	user := fixtures.AddUser(store, "test", "user", true)
	userToken, _ := CreateUserToken(user)
	b, _ := json.Marshal(tc.input)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", userToken)
	res, _ := app.Test(req)
	if res.StatusCode != tc.status {
		t.Fatalf("expected status %d, got %d", tc.status, res.StatusCode)
	}
	if tc.ttype == TESTPASS {
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
	} else if tc.ttype == TESTFAIL {
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
		ttype: TESTPASS,
		input: userTest{
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser@test.com",
			Password:  "Testp@ssword123",
		},
		expected: expected{
			status: http.StatusCreated,
			body: map[string]string{
				"firstName": "Test",
				"lastName":  "User",
				"email":     "testuser@test.com",
			},
		},
	},
	{
		name:  "invalid email and password",
		ttype: TESTFAIL,
		input: userTest{
			FirstName: "Test",
			LastName:  "User",
			Email:     "testuser",
			Password:  "testpassword123",
		},
		expected: expected{
			status: http.StatusBadRequest,
			body: map[string]string{
				"email":    "invalid email address",
				"password": "password must contain at least one uppercase letter, password must contain at least one special character",
			},
		},
	},
}
