//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllExpense(t *testing.T) {
	seedUser(t)
	var us []Expense

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&us)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(us), 0)
}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "iPhone 17 Pro Max 2TB",
		"amount": 76900,
		"note": "birthday gift from my love",
		"tags": ["gadget"]
	}`)
	var u Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&u)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, u.ID)
	assert.Equal(t, "iPhone 17 Pro Max 2TB", u.Title)
	assert.Equal(t, 76900, u.Amount)
}

func TestGetUserByID(t *testing.T) {
	c := seedUser(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
}

func seedUser(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "iPhone 17 Pro Max 2TB",
		"amount": 76900,
		"note": "birthday gift from my love",
		"tags": ["gadget"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return c
}

func TestUpdateUserByID(t *testing.T) {
	t.Skip("TODO: implement me")
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
