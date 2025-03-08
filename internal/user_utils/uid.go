package user_utils

import (
	"fmt"
)

func UIdToInt(uid string) (int, error) {
	var uidInt int
	if _, err := fmt.Sscanf(uid, "%d", &uidInt); err != nil {
		return -1, fmt.Errorf("failed to convert uid to int: %w", err)
	}
	return uidInt, nil
}
