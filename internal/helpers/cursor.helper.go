package helpers

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ValidateCursor(cursor string) (lastMsgId *uuid.UUID, lastMsgDate *time.Time, err error) {
	if cursor == "" {
		return nil, nil, nil
	}

	// cursor is base64 version of `{{id}}&{{date}}`
	// eg: decoded cursor -> "6d35b78d-a9ef-44dc-bd6e-d0b60c39a0e0&2024-04-18T18:38:53.739Z"
	// encoded cursor -> "NmQzNWI3OGQtYTllZi00NGRjLWJkNmUtZDBiNjBjMzlhMGUwJjIwMjQtMDQtMThUMTg6Mzg6NTMuNzM5Wg=="

	decoded, err := Base64Decode(cursor)

	if err != nil {
		return nil, nil, errors.New("Invalid cursor")
	}

	split := strings.Split(decoded, "&")

	if len(split) != 2 {
		return nil, nil, errors.New("Malformed cursor")
	}

	uuidParsed, err := uuid.Parse(split[0])

	if err != nil {
		return nil, nil, errors.New("Invalid message id")

	} else {
		lastMsgId = &uuidParsed
	}

	t, err := time.Parse(time.RFC3339, split[1])

	if err != nil {
		return nil, nil, errors.New("Invalid date")
	}

	lastMsgDate = &t

	return lastMsgId, lastMsgDate, nil
}
