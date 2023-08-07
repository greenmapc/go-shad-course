//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

var AuthError *kvapi.AuthError
var OutdatedVersion *kvapi.ConflictError

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	response, err := c.Get(&kvapi.GetRequest{
		Key: key,
	})
	var oldValue *string
	var oldVersion uuid.UUID
	newVersion := uuid.Must(uuid.NewV4())

	switch true {
	case err == nil:
		oldValue = &response.Value
		oldVersion = response.Version
	case errors.Is(err, kvapi.ErrKeyNotFound):
		oldValue = nil
	case errors.As(err, &AuthError):
		return err
	default:
		return UpdateValue(c, key, updateFn)
	}

	newValue, err := updateFn(oldValue)
	if err != nil {
		return err
	}

	for {
		_, setError := c.Set(&kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: oldVersion,
			NewVersion: newVersion,
		})

		switch true {
		case setError == nil:
			return nil
		case errors.As(setError, &AuthError):
			return setError
		case errors.Is(setError, kvapi.ErrKeyNotFound):
			newValue, err = updateFn(nil)
			if err != nil {
				return err
			}
			oldVersion = uuid.UUID{}
		case errors.As(setError, &OutdatedVersion):
			if OutdatedVersion.ExpectedVersion == newVersion {
				return nil
			}
			return UpdateValue(c, key, updateFn)
		default:
			continue
		}
	}
}
