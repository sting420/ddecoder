package decoder

import (
	"errors"
	"net/url"
)

func ParseData(data string, seed EncodingSeed) (string, string, string, error) {

	var payload string
	var cid string
	var hash string
	var query url.Values
	var err error

	var missingValError = errors.New("missing value in query")
	if seed == Captcha {
		urlVals, err := url.Parse(data)
		if err != nil {
			return "", "", "", err
		}
		query = urlVals.Query()
		if !query.Has("ddCaptchaEncodedPayload") || !query.Has("hash") {
			return "", "", "", missingValError
		}
		payload, err = url.QueryUnescape(query.Get("ddCaptchaEncodedPayload"))
		if err != nil {
			return "", "", "", err
		}
		hash = query.Get("hash")

	} else {
		query, err = url.ParseQuery(data)
		if err != nil {
			return "", "", "", err
		}
		if seed == Interstitial {
			if !query.Has("payload") || !query.Has("hash") {
				return "", "", "", missingValError
			}
			payload = query.Get("payload")
			hash = query.Get("hash")

		} else {
			if !query.Has("jspl") || !query.Has("ddk") {
				return "", "", "", missingValError
			}
			payload = query.Get("jspl")
			hash = query.Get("ddk")

		}
	}
	if !query.Has("cid") {
		return "", "", "", missingValError
	}
	cid = query.Get("cid")

	return payload, cid, hash, nil
}
