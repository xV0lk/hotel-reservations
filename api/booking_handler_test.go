package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xV0lk/hotel-reservations/db"
	"github.com/xV0lk/hotel-reservations/db/fixtures"
	"github.com/xV0lk/hotel-reservations/types"
)

func TestBookings(t *testing.T) {
	db := setup(t)
	defer db.Drop(t)

	user := fixtures.AddUser(db.Store, "test", "user", false)
	hotel := fixtures.AddHotel(db.Store, "test hotel", "test address", 4, 2)
	room := fixtures.AddRoom(db.Store, types.Single, 100, hotel.ID)
	fixtures.AddBooking(db.Store, user.ID, room, time.Now(), time.Now().AddDate(0, 0, 5), 1)

	app := fiber.New()
	bookingHandler := NewBookingHandler(db.Store)
	api := app.Group("/", JWTAuth(db.Store.User))
	admin := api.Group("/bookings", AdminAuth)
	api.Get("/booking", bookingHandler.HandleGetBooking)
	admin.Get("/", bookingHandler.HandleGetBookings)

	for _, tc := range getBookingTests {
		t.Run(tc.name, func(t *testing.T) {
			tc.getBookings(t, app, db.Store)
		})
	}
}

func (tc bookingCase) getBookings(t *testing.T, app *fiber.App, store *db.Store) {
	user := fixtures.AddUser(store, tc.input.name, tc.input.lname, tc.input.admin)
	userToken, _ := CreateUserToken(user)
	req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
	req.Header.Add("Authorization", userToken)
	res, _ := app.Test(req)
	if res.StatusCode != tc.status {
		t.Fatalf("expected status %d, got %d", tc.status, res.StatusCode)
	}
	if tc.ttype == TESTPASS {
		var resBody []types.Booking
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		if len(resBody) < 1 {
			t.Fatalf("expected at least 1 booking, got %d", len(resBody))
		}
	}
	if tc.ttype == TESTFAIL {
		var resBody failAuthResponse
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		if resBody.Error != tc.body.(string) {
			t.Errorf("expected %s, got %s", tc.body.(string), resBody.Error)
		}
	}
}

type getBookingTest struct {
	name  string
	lname string
	admin bool
}

type bookingCase struct {
	name  string
	ttype string
	input getBookingTest
	expected
}

var getBookingTests = []bookingCase{
	{
		name:  "Admin get all bookings",
		ttype: TESTPASS,
		input: getBookingTest{
			name:  "test",
			lname: "admin",
			admin: true,
		},
		expected: expected{
			status: http.StatusOK,
			body:   nil,
		},
	},
	{
		name:  "Not Admin get all bookings",
		ttype: TESTFAIL,
		input: getBookingTest{
			name:  "test",
			lname: "notadmin",
			admin: false,
		},
		expected: expected{
			status: http.StatusForbidden,
			body:   "Unauthorized",
		},
	},
}
