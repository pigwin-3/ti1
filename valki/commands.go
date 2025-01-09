package valki

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

func SetValkeyValue(ctx context.Context, client valkey.Client, key, value string) error {
	err := client.Do(ctx, client.B().Set().Key(key).Value(value).Build()).Error()
	if err != nil {
		return fmt.Errorf("failed to set value in Valkey: %v", err)
	}
	return nil
}

func GetValkeyValue(ctx context.Context, client valkey.Client, key string) (string, error) {
	value, err := client.Do(ctx, client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return "hehe", nil
		//return "", fmt.Errorf("failed to get value from Valkey: %v", err)
	}
	return value, nil
}
