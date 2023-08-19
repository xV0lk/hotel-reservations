package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDbUri  = "mongodb://localhost:27017"
	testDbName = "hotel-reservations-test"
)

type expected struct {
	status int
	body   any
}

type postTest[T any] struct {
	name  string
	ttype string
	input T
	expected
}

type failUserResponse struct {
	Error map[string]string `json:"error"`
}

type testdb struct {
	db.UserStore
}

type userTest types.NewUserParams

func (tdb *testdb) Drop(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbUri))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to test MongoDB!")
	return &testdb{db.NewMongoUserStore(client)}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.Drop(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	for _, tc := range userTests {
		t.Run(tc.name, func(t *testing.T) {
			tc.testUser(t, app)
		})
	}

}

func (tc postTest[userTest]) testUser(t *testing.T, app *fiber.App) {
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

var userTests = []postTest[userTest]{
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
