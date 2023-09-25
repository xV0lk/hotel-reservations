package api

import (
	"bytes"
	"encoding/json"
	"net/http"
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
	db := setup(t)
	defer db.Drop(t)
	insertedUser := fixtures.AddUser(db.Store, "test", "user", false)

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	authHandler := NewAuthHandler(db.Store.User)
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
	if tc.ttype == TESTPASS {
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
	} else if tc.ttype == TESTFAIL {
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
		ttype: TESTPASS,
		input: AuthParams{
			Email:    "test@user.com",
			Password: "test_user_P4$$",
		},
		expected: expected{
			status: http.StatusOK,
			body:   nil,
		},
	},
	{
		name:  "invalid password",
		ttype: TESTFAIL,
		input: AuthParams{
			Email:    "test@user.com",
			Password: "t3tPassword",
		},
		expected: expected{
			status: http.StatusBadRequest,
			body:   "Invalid password",
		},
	},
	{
		name:  "invalid email",
		ttype: TESTFAIL,
		input: AuthParams{
			Email:    "test@testmail.comm",
			Password: "t3tPassword",
		},
		expected: expected{
			status: http.StatusNotFound,
			body:   "The id you provided is invalid",
		},
	},
}
