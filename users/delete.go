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

	switch resp.StatusCode {
	case http.StatusNoContent:
		return resp.Body.Close()
	case http.StatusNotFound:
		if err := resp.Body.Close(); err != nil {
			return err
		}
		// if the user is not found, it might be an unaccepted user invitation:
		return revokeInvitation(c, ctx, rawURL, id)
	default:
		_ = resp.Body.Close()
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
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

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return resp.Body.Close()
	}
	_ = resp.Body.Close()
	return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
}
