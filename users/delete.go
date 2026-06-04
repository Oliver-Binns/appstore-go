package users

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/oliver-binns/googleplay-go/networking"
)

func Delete(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "users", id)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, parsedURL.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)

	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		// if the user is not found, it might be an unaccepted user invitation:
		return revokeInvitation(c, ctx, rawURL, id)
	} else if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	} else {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
}

func revokeInvitation(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	// Parse the raw URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, "userInvitations", id)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, parsedURL.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	} else {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
}
