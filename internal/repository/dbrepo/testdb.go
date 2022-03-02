package dbrepo

import (
	"errors"
	"time"

	"github.com/ashrielbrian/go_bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID > 2 {
		return 0, errors.New("failed to insert reservation")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ts, _ := time.Parse("2006-01-02", "1970-01-01")
	if r.StartDate.Equal(ts) {
		return errors.New("failed to insert room restriction")
	}

	return nil

}

// SearchAvailabilityByDatesByRoomID returns true is there is room availability for room ID; returns false otherwise
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 1 {
		return true, nil
	}
	return false, nil

}

// SearchAvailabiltyForAllRooms returns a slice of all available rooms, if any, for a given date range
func (m *testDBRepo) SearchAvailabiltyForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomByID gets the room details by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("no such room ID")
	}
	return room, nil
}
