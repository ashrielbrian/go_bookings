package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ashrielbrian/go_bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	layout := "2006-01-02"

	sd, _ := time.Parse(layout, "2022-01-01")
	ed, _ := time.Parse(layout, "2022-01-05")
	var res = models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
		StartDate: sd,
		EndDate:   ed,
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", res)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// test case when session is empty
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test case when room ID does not exist
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	res.RoomID = 100
	session.Put(ctx, "reservation", res)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	layout := "2006-01-02"

	sd, _ := time.Parse(layout, "2022-01-01")
	ed, _ := time.Parse(layout, "2022-01-05")

	var res = models.Reservation{

		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
		StartDate: sd,
		EndDate:   ed,
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "j@smith.com")
	postedData.Add("phone", "012384131")

	// reqBody := "first_name=John"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=012393123i1")
	// req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", res)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation expected status code %d, got %d", http.StatusSeeOther, rr.Code)
	}

	// test case when form failed to parse by setting an empty body
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", res)
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation expected status code %d due to missing body, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test case when form is not valid by removing first_name
	postedData.Del("first_name")
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", res)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation expected status code %d due to invalid form, got %d", http.StatusOK, rr.Code)
	}

	// test case with failed reservation by setting a large roomID
	postedData.Add("first_name", "John")
	res.RoomID = 1000
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", res)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation expected status code %d due to failed reservation, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test case with a failed room restriction insertion by setting the StartDate to be an incorrect value
	res.RoomID = 1
	res.StartDate, _ = time.Parse("2006-01-02", "1970-01-01")

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", res)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation expected status code %d due to failed restriction, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	postedData := url.Values{}

	postedData.Add("start", "2022-01-01")
	postedData.Add("end", "2022-01-09")
	postedData.Add("room_id", "1")

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	var body jsonResponse
	json.Unmarshal(rr.Body.Bytes(), &body)

	if rr.Code != http.StatusOK && !body.OK {
		t.Errorf("Expected status code %d, instead got %d", http.StatusOK, rr.Code)
	}

	// test invalid start date
	postedData.Set("start", "invalid")
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)

	json.Unmarshal(rr.Body.Bytes(), &body)
	expectedStr := "Error parsing start date."
	if body.Message != "Error parsing start date." {
		t.Errorf("Expected message `%s`, instead got `%s`", expectedStr, body.Message)
	}

	// test invalid end date
	postedData.Set("start", "2022-01-01")
	postedData.Set("end", "invalid")
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)

	json.Unmarshal(rr.Body.Bytes(), &body)
	expectedStr = "Error parsing end date."
	if body.Message != expectedStr {
		t.Errorf("Expected message `%s`, instead got `%s`", expectedStr, body.Message)
	}

	// test invalid form setting empty body

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler.ServeHTTP(rr, req)

	json.Unmarshal(rr.Body.Bytes(), &body)
	expectedStr = "Error parsing form."
	if body.Message != expectedStr {
		t.Errorf("Expected message `%s`, instead got `%s`", expectedStr, body.Message)
	}

}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
		return nil
	}

	return ctx
}
