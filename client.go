package latestgo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const versionPrefix = "go"

// Version returns the latest stable go version.
func Version(ctx context.Context) (string, error) {
	req, err := http.NewRequest("GET", "https://golang.org/dl/", nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	s := doc.Find("div.toggleVisible").First()
	ver, exists := s.Attr("id")
	if !exists {
		return "", errors.New("id attribute not found in div.toggleVisible")
	}

	if !strings.HasPrefix(ver, versionPrefix) {
		return "", fmt.Errorf("id attribute div.toggleVisible does not start with %s, ver=%s", versionPrefix, ver)
	}

	return ver[len(versionPrefix):], nil
}
