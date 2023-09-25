package db

import (
	"regexp"
	"strings"
)

const (
	DBNAME = "hotel-reservations"
	DBURI  = "mongodb://localhost:27017"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func FormatMongoE(e error) string {
	re := regexp.MustCompile(`\{([^}]*)\}`)
	matches := re.FindStringSubmatch(e.Error())
	if len(matches) > 1 {
		return strings.ReplaceAll(strings.TrimSpace(matches[1]), "\"", "'")
	}
	return ""
}
