package db

import "errors"

func GetOffset(page, limit int) (int, error) {
	if page < 1 {
		return 0, errors.New("err page")
	}
	if limit < 1 || limit > 100 {
		return 0, errors.New("err limit")
	}
	offset := (page - 1) * limit

	return offset, nil
}
