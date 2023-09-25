package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db/fixtures"
	"github.com/xV0lk/hotel-reservations/types"
)

type failAuthResponse struct {
	Error string `json:"error"`
}

func TestAuth(t *testing.T) {
	tdb := setup(t)
	defer tdb.Drop(t)
	insertedUser := fixtures.AddUser(tdb.Store, "test", "user", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	for _, tc := range authTests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testAuth(t, app, insertedUser)
		})
	}
}

func (tc testCase[AuthParams]) testAuth(t *testing.T, app *fiber.App, expectedUser *types.User) {
	b, _ := json.Marshal(tc.input)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, _ := app.Test(req)
	if res.StatusCode != tc.status {
		t.Fatalf("expected status %d, got %s", tc.status, res.Status)
	}
	if tc.ttype == "pass" {
		var authResp AuthResponse
		if err := json.NewDecoder(res.Body).Decode(&authResp); err != nil {
			t.Fatal(err)
		}
		if authResp.Token == "" {
			t.Error("expected token to be non-empty")
		}
		expectedUser.Password = ""
		if !reflect.DeepEqual(authResp.User, expectedUser) {
			t.Errorf("expected %v, got %v", expectedUser, authResp.User)
		}
	} else if tc.ttype == "fail" {
		var resBody failAuthResponse
		err := json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			t.Fatal(err)
		}
		if resBody.Error != tc.body.(string) {
			t.Errorf("expected %s, got %s", tc.body.(string), resBody.Error)
		}
	}
}

var authTests = []testCase[AuthParams]{
	{
		name:  "valid user and password",
		ttype: "pass",
		input: AuthParams{
			Email:    "test@user.com",
			Password: "test_user_P4$$",
		},
		expected: expected{
			status: fiber.StatusOK,
			body:   nil,
		},
	},
	{
		name:  "invalid password",
		ttype: "fail",
		input: AuthParams{
			Email:    "test@testmail.com",
			Password: "t3tPassword",
		},
		expected: expected{
			status: fiber.StatusBadRequest,
			body:   "Invalid Password",
		},
	},
	{
		name:  "invalid email",
		ttype: "fail",
		input: AuthParams{
			Email:    "test@testmail.comm",
			Password: "t3tPassword",
		},
		expected: expected{
			status: fiber.StatusNotFound,
			body:   "Not Found",
		},
	},
}
