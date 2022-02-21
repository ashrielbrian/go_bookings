package forms

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("GET", "/some-url", nil)

	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("Expected a valid form")
	}

	form.Values = url.Values{}
	form.Add("a", "b")
	form.Add("c", "d")

	if !form.Valid() {
		t.Error("Expected a valid form")
	}

	form.Errors.Add("a", "test failed msg")

	if form.Errors.Get("a") == "" {
		t.Error("Should have an error, but did not get one.")
	}

	if form.Errors.Get("zz") != "" {
		t.Error("Should not have an error, but got one.")
	}

	if form.Valid() {
		t.Error("Expected form to be invalid; got valid instead")
	}
}

func TestForm_Required(t *testing.T) {

	// Test for negatively valid required form
	r := httptest.NewRequest("GET", "/some-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("Expected form to be invalid; got valid instead")
	}

	// Test for positively valid required form
	r = httptest.NewRequest("GET", "/some-url", nil)
	form = New(r.PostForm)

	form.Values = url.Values{}
	form.Add("a", "test")
	form.Add("b", "test")
	form.Add("c", "test")

	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("Expected form to be valid; got invalid instead")
	}
}

func TestForm_Has(t *testing.T) {
	key := "a"
	nonExistentKey := "b"

	postedData := url.Values{
		key: []string{"test"},
	}

	form := New(postedData)

	if !form.Has(key) {
		t.Error(fmt.Sprintf("Expected form has %s; got false instead", key))
	}
	if form.Has(nonExistentKey) {
		t.Error(fmt.Sprintf("Expected form to not have %s; got true instead", nonExistentKey))
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{
		"pass": []string{"meetexpectation"},
		"fail": []string{"short"},
	}

	form := New(postedData)

	if form.MinLength("fail", 10) {
		t.Error("Expected min length to fail; got success instead.")
	}
	if form.MinLength("nonexistent", 10) {
		t.Error("Form shows min length for non existent path")
	}
	if !form.MinLength("pass", 10) {
		t.Error("Expected min length to pass; got fail instead.")
	}
}

func TestForm_Email(t *testing.T) {
	postedData := url.Values{
		"pass_email":   []string{"a@b.com"},
		"fail_email":   []string{"a@b"},
		"fail_email_2": []string{"a"},
	}

	form := New(postedData)

	if !form.IsEmail("pass_email") {
		t.Error("Expected valid email passing; failed instead.")
	}
	if form.IsEmail("fail_email") {
		t.Error("Expected invalid email failing; passed instead.")
	}
	if form.IsEmail("fail_email_2") {
		t.Error("Expected invalid email failing; passed instead.")
	}
	if form.IsEmail("nonexistent") {
		t.Error("Form shows non existent field is invalid email.")
	}

}
