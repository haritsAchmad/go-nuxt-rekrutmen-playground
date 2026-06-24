package request

import (
	"errors"
	"net/http"
	"strconv"
)

func GetIDFromQuery(r *http.Request) (int, error) {
	idText := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idText)

	if err != nil || id <= 0 {
		return 0, errors.New("ID tidak valid")
	}

	return id, nil
}