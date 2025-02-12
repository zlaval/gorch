package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BodyParserFunc[T any] func(response *http.Response) (T, error)

func SuccessResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&data)
}

func Post[T any](url string, req any, fn BodyParserFunc[T]) (T, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return *new(T), fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return *new(T), fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()
	return fn(resp)
}

func ExtractBody[T any](resp *http.Response) (T, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(T), fmt.Errorf("read response: %w", err)
	}
	var r T
	if err := json.Unmarshal(b, &r); err != nil {
		return *new(T), fmt.Errorf("unmarshal response: %w", err)
	}
	return r, nil
}

func OmitBody(_ *http.Response) (any, error) {
	return nil, nil
}
