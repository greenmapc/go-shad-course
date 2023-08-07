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

	keyNotFound := errors.Is(err, kvapi.ErrKeyNotFound)

	if err != nil && errors.As(err, &AuthError) {
		return err
	}

	if keyNotFound || err == nil {
		var newValue string

		if keyNotFound {
			updateWithNilResult, updateWithNilErr := updateFn(nil)
			if updateWithNilErr != nil {
				return updateWithNilErr
			}
			newValue = updateWithNilResult
		} else {
			updateResult, updateErr := updateFn(&response.Value)
			if updateErr != nil {
				return updateErr
			}
			newValue = updateResult
		}

		var oldVersionIfKeyNotFound uuid.UUID
		var oldVersion uuid.UUID
		newVersion := uuid.Must(uuid.NewV4())

		if !keyNotFound {
			oldVersion = response.Version
		}
		for {
			_, setError := c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: oldVersion,
				NewVersion: newVersion,
			})

			if setError != nil && errors.As(setError, &AuthError) {
				return setError
			}
			if errors.Is(setError, kvapi.ErrKeyNotFound) {

				nilValue, err := updateFn(nil)
				if err != nil {
					return err
				}
				for {
					_, setRepeatError := c.Set(&kvapi.SetRequest{
						Key:        key,
						Value:      nilValue,
						OldVersion: oldVersionIfKeyNotFound,
						NewVersion: uuid.Must(uuid.NewV4()),
					})

					if setRepeatError != nil && errors.As(setRepeatError, &AuthError) {
						return setRepeatError
					}
					if setRepeatError != nil {
						continue
					}
					return nil
				}
			}
			if errors.As(setError, &OutdatedVersion) {
				if OutdatedVersion.ExpectedVersion == newVersion {
					return nil
				}
				return UpdateValue(c, key, updateFn)
			}
			if setError != nil {
				continue
			}
			return nil
		}
	} else {
		return UpdateValue(c, key, updateFn)
	}
}
