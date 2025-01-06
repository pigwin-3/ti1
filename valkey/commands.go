package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

func SetValkeyValue(client valkey.Client, key, value string) error {
	ctx := context.Background()
	err := client.Do(ctx, client.B().Set().Key(key).Value(value).Build()).Error()
	if err != nil {
		return fmt.Errorf("failed to set value in Valkey: %v", err)
	}
	return nil
}

func GetValkeyValue(client valkey.Client, key string) (string, error) {
	ctx := context.Background()
	value, err := client.Do(ctx, client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return "", fmt.Errorf("failed to get value from Valkey: %v", err)
	}
	return value, nil
}
