package wiki

import (
	"net/url"
)

// Конвертируем запрос для использования в качестве части URL
func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
