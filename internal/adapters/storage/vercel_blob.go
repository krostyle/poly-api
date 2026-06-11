package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const blobAPIBase = "https://blob.vercel-storage.com"

type VercelBlobStorage struct {
	token string
}

func NewVercelBlobStorage() *VercelBlobStorage {
	return &VercelBlobStorage{token: os.Getenv("BLOB_READ_WRITE_TOKEN")}
}

type blobPutResponse struct {
	URL string `json:"url"`
}

func (s *VercelBlobStorage) Upload(ctx context.Context, name string, content []byte, contentType string) (string, error) {
	safeName := sanitizeName(name)
	url := fmt.Sprintf("%s/%s", blobAPIBase, safeName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(content))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-api-version", "7")
	req.Header.Set("x-add-random-suffix", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("vercel blob upload failed (%d): %s", resp.StatusCode, string(body))
	}

	var result blobPutResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.URL, nil
}

func (s *VercelBlobStorage) Delete(ctx context.Context, blobURL string) error {
	body, _ := json.Marshal(map[string][]string{"urls": {blobURL}})

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, blobAPIBase, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-version", "7")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("vercel blob delete failed (%d): %s", resp.StatusCode, string(b))
	}
	return nil
}

func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
