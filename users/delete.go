package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/openapi"
	"github.com/oliver-binns/googleplay-go/networking"
)

func Delete(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	resp, err := apiClient.UsersDeleteInstance(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		// if the user is not found, it might be an unaccepted user invitation:
		return revokeInvitation(c, ctx, rawURL, id)
	}
	return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
}

func revokeInvitation(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	apiClient, err := openapi.NewClient(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	resp, err := apiClient.UserInvitationsDeleteInstance(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil
	}
	return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
}
