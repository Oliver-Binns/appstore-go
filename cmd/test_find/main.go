package main

import (
	"context"
	"fmt"
	"os"

	appstore "github.com/oliver-binns/appstore-go"
)

func main() {
	privateKey, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read key: %v\n", err)
		os.Exit(1)
	}

	client := appstore.AppStoreClient(
		os.Args[2],
		os.Args[3],
		string(privateKey),
	)

	email := os.Args[4]
	fmt.Printf("Looking up user by email: %s\n", email)

	user, err := client.FindUserByEmail(context.Background(), email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if user == nil {
		fmt.Println("User not found")
		return
	}
	fmt.Printf("Found user: ID=%s, Name=%s %s, HasAcceptedInvite=%v\n",
		user.ID, user.FirstName, user.LastName, user.HasAcceptedInvite)
}
