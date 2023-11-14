package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("connection error")
	ErrNotFound        = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid server response")
	ErrInvalid         = errors.New("invalid data")
	ErrNotNumber       = errors.New("not a number")
)

const timeFormat = "Jan/02 @15:04"

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type response struct {
	Results      []item `json:"results"`
	Date         int    `json:"date"`
	TotalResults int    `json:"total_results"`
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}

func getItems(url string) ([]item, error) {
	r, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)

		if err != nil {
			return nil, fmt.Errorf("cannot read body: %w", err)
		}

		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}

		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	var resp response

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("%w: no results found", ErrNotFound)
	}

	return resp.Results, nil
}

func getAll(apiRoot string) ([]item, error) {
	u := fmt.Sprintf("%s/todo", apiRoot)

	return getItems(u)
}

func getOne(apiRoot string, id int) (item, error) {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)

	items, err := getItems(u)

	if err != nil {
		return item{}, err
	}

	if len(items) != 1 {
		return item{}, fmt.Errorf("%w: invalid results", ErrInvalid)
	}

	return items[0], nil
}

func sendRequest(url, method, contentType string, expStatus int, body io.Reader) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	r, err := newClient().Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if r.StatusCode != expStatus {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s", err, msg)
	}

	return nil
}

func addItem(apiRoot, task string) error {
	u := fmt.Sprintf("%s/todo", apiRoot)

	item := struct {
		Task string `json:"task"`
	}{
		Task: task,
	}

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(item); err != nil {
		return err
	}

	return sendRequest(u, http.MethodPost, "application/json", http.StatusCreated, &body)
}

func completeItem(apiRoot string, id int) error {
	u := fmt.Sprintf("%s/todo/%d?complete", apiRoot, id)

	return sendRequest(u, http.MethodPatch, "", http.StatusNoContent, nil)
}

func deleteItem(apiRoot string, id int) error {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)

	return sendRequest(u, http.MethodDelete, "", http.StatusNoContent, nil)
}
