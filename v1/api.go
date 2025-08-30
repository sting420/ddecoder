package decoder

type EncodingSeed int32

const (
	Interstitial EncodingSeed = -248865861
	Captcha      EncodingSeed = -1380527010
	Tags         EncodingSeed = 374833564
)

func Decode(data string, seed EncodingSeed) (string, error) {
	payload, cid, hash, err := ParseData(data, seed)
	if err != nil {
		return "", err
	}

	decoded, err := DecodePayload(payload, cid, hash, int32(seed))
	if err != nil {
		return "", err
	}

	return decoded, nil
}
