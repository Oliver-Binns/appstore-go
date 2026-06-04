package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oliver-binns/appstore-go/networking"
	"github.com/oliver-binns/appstore-go/openapi"
)

func Delete(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.UsersDeleteInstanceWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode() == http.StatusNotFound {
		return revokeInvitation(c, ctx, rawURL, id)
	}
	return fmt.Errorf("unexpected response code: %d", resp.StatusCode())
}

func revokeInvitation(c networking.HTTPClient, ctx context.Context, rawURL string, id string) error {
	client, err := openapi.NewClientWithResponses(rawURL, openapi.WithHTTPClient(c))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.UserInvitationsDeleteInstanceWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if resp.StatusCode() == http.StatusNoContent || resp.StatusCode() == http.StatusNotFound {
		return nil
	}
	return fmt.Errorf("unexpected response code: %d", resp.StatusCode())
}
