package bingxgo

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
)

func parseKlineData(data KlineDataRaw, interval string) (KlineData, error) {
	if len(data) < 8 {
		return KlineData{}, errors.New("invalid kline received")
	}

	return KlineData{
		StartTime: int64(data[0]),
		EndTime:   int64(data[6]),
		Interval:  interval,
		Open:      data[1],
		High:      data[2],
		Low:       data[3],
		Close:     data[4],
		Volume:    data[7],
	}, nil
}

func DecodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decodedMsg, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decodedMsg, nil
}
