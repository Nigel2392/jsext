package encoding

import (
	"time"

	"github.com/Nigel2392/jsext/v2"
)

func MarshalCookie(name string, value any, ttl time.Duration) error {
	var encoded, err = EncodeJSON[string](value)
	if err != nil {
		return err
	}
	encoded, err = EncodeBase64[string](encoded)
	if err != nil {
		return err
	}
	return jsext.SetCookie(name, encoded, ttl)
}

func UnmarshalCookie(name string, dst any) error {
	var encoded = jsext.GetCookie(name)
	if encoded == "" {
		return nil
	}
	var decoded, err = DecodeBase64[string](encoded)
	if err != nil {
		return err
	}
	return DecodeJSON(decoded, dst)
}
