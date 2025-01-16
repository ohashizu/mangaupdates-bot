package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/cenkalti/backoff/v4"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/valyala/fasthttp"
)

var (
	username = os.Getenv("MU_USERNAME")
	password = os.Getenv("MU_PASSWORD")
)

func MuConvertOldId(oldID int64) string {
	const base = 36
	const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

	if oldID == 0 {
		return "0"
	}

	var result strings.Builder
	for oldID > 0 {
		remainder := oldID % base
		result.WriteString(string(digits[remainder]))
		oldID /= base
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func MuConvertNewId(newID string) (int64, error) {
	const base = 36
	return strconv.ParseInt(newID, base, 64)
}

func MuCleanupDescription(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	const maxLength = 1024
	var stop bool

	var f func(*html.Node)
	f = func(n *html.Node) {
		if stop {
			return
		}
		if n.Type == html.TextNode {
			remaining := maxLength - buf.Len()
			if remaining <= 0 {
				stop = true
				return
			}
			text := n.Data
			if len(text) > remaining {
				buf.WriteString(text[:remaining])
				stop = true
				return
			}
			buf.WriteString(text)
		} else if n.Type == html.ElementNode && strings.ToLower(n.Data) == "br" {
			if buf.Len() < maxLength {
				buf.WriteString("\n")
				if buf.Len() >= maxLength {
					stop = true
					return
				}
			} else {
				stop = true
				return
			}
		}
		for c := n.FirstChild; c != nil && !stop; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return buf.String(), nil
}

func MuLogin() (*MuLoginResponse, error) {

	login := MuLoginRequest{
		Username: username,
		Password: password,
	}
	loginJSON, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("PUT")
	req.SetRequestURI("https://api.mangaupdates.com/v1/account/login")
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(loginJSON)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	respBody := &MuLoginResponse{}
	if err = json.Unmarshal(resp.Body(), respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func MuLogout(b *mubot.Bot) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.SetRequestURI("https://api.mangaupdates.com/v1/account/logout")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.MuToken)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return err
	}

	respBody := &MuLogoutResponse{}
	if err = json.Unmarshal(resp.Body(), respBody); err != nil {
		return err
	}

	return nil
}

func MuGetSeriesInfo(b *mubot.Bot, seriesId int64) (*MuSeriesInfoResponse, error) {
	var respBody *MuSeriesInfoResponse

	operation := func() error {
		resp, statusCode, err := muGetRequest(
			"https://api.mangaupdates.com/v1/series/"+strconv.FormatInt(seriesId, 10),
			b.MuToken,
		)

		if err != nil {
			return fmt.Errorf(
				"failed to fetch series info: %s, %s, %d",
				err.Error(),
				string(resp),
				seriesId,
			)
		}

		if statusCode == 200 {
			respBody = &MuSeriesInfoResponse{}
			if err := json.Unmarshal(resp, respBody); err != nil {
				return fmt.Errorf("failed to unmarshal series info: %s", err.Error())
			}
			return nil
		} else if statusCode >= 500 && statusCode < 600 {
			return fmt.Errorf("series info server error: %d", statusCode)
		} else {
			return backoff.Permanent(fmt.Errorf("series info client error: %d", statusCode))
		}
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 2 * time.Minute
	backoffConfig.MaxElapsedTime = 5 * time.Minute
	backoffConfig.Multiplier = 1.5
	backoffConfig.RandomizationFactor = 0.5

	retryPolicy := backoff.WithMaxRetries(backoffConfig, 20)

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get series info after retries: %v, Series ID: %d",
			err,
			seriesId,
		)
	}
	if respBody == nil {
		return nil, fmt.Errorf("MuGetSeriesInfo: respBody is nil after successful unmarshal")
	}

	return respBody, nil
}

func MuPostSearchGroups(b *mubot.Bot, groupName string) (*MuSearchGroupsResponse, error) {
	var respBody *MuSearchGroupsResponse

	operation := func() error {
		resp, statusCode, err := muPostRequest(
			"https://api.mangaupdates.com/v1/groups/search",
			b.MuToken,
			MuSearchGroupsRequest{
				Search:  groupName,
				PerPage: 10,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to post search groups: %s, %s, %s",
				err.Error(),
				string(resp),
				groupName,
			)
		}

		if statusCode == 200 {
			respBody = &MuSearchGroupsResponse{}
			if err = json.Unmarshal(resp, respBody); err != nil {
				return fmt.Errorf("failed to unmarshal search groups: %s", err.Error())
			}
			return nil
		} else if statusCode >= 500 && statusCode < 600 {
			return fmt.Errorf("search groups server error: %d", statusCode)
		} else {
			return backoff.Permanent(fmt.Errorf("search groups client error: %d", statusCode))
		}
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 2 * time.Minute
	backoffConfig.MaxElapsedTime = 5 * time.Minute
	backoffConfig.Multiplier = 1.5
	backoffConfig.RandomizationFactor = 0.5

	retryPolicy := backoff.WithMaxRetries(backoffConfig, 20)

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get search groups after retries: %v, group name: %s",
			err,
			groupName,
		)
	}
	if respBody == nil {
		return nil, fmt.Errorf("MuPostSearchGroups: respBody is nil after successful unmarshal")
	}

	return respBody, nil
}

func MuPostSearchSeries(b *mubot.Bot, seriesName string) (*MuSearchSeriesResponse, error) {
	var respBody *MuSearchSeriesResponse

	operation := func() error {
		resp, statusCode, err := muPostRequest(
			"https://api.mangaupdates.com/v1/series/search",
			b.MuToken,
			MuSearchSeriesRequest{
				Search:  seriesName,
				PerPage: 10,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to post search series: %s, %s, %s",
				err.Error(),
				string(resp),
				seriesName,
			)
		}

		if statusCode == 200 {
			respBody = &MuSearchSeriesResponse{}
			if err := json.Unmarshal(resp, respBody); err != nil {
				return fmt.Errorf("failed to unmarshal search series: %s", err.Error())
			}
			return nil
		} else if statusCode >= 500 && statusCode < 600 {
			return fmt.Errorf("search series server error: %d", statusCode)
		} else {
			return backoff.Permanent(fmt.Errorf("search series client error: %d", statusCode))
		}
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 2 * time.Minute
	backoffConfig.MaxElapsedTime = 5 * time.Minute
	backoffConfig.Multiplier = 1.5
	backoffConfig.RandomizationFactor = 0.5

	retryPolicy := backoff.WithMaxRetries(backoffConfig, 20)

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get search series after retries: %v, series name: %s",
			err,
			seriesName,
		)
	}
	if respBody == nil {
		return nil, fmt.Errorf("MuPostSearchSeries: respBody is nil after successful unmarshal")
	}

	return respBody, nil
}

func MuGetSeriesGroups(b *mubot.Bot, seriesId int64) (*MuSeriesGroupsResponse, error) {
	var respBody *MuSeriesGroupsResponse

	operation := func() error {
		resp, statusCode, err := muGetRequest(
			"https://api.mangaupdates.com/v1/series/"+strconv.FormatInt(seriesId, 10)+"/groups",
			b.MuToken,
		)

		if err != nil {
			return fmt.Errorf(
				"failed to fetch series groups: %s, %s, %d",
				err.Error(),
				string(resp),
				seriesId,
			)
		}

		if statusCode == 200 {
			respBody = &MuSeriesGroupsResponse{}
			if err := json.Unmarshal(resp, respBody); err != nil {
				return fmt.Errorf("failed to unmarshal series groups: %s", err.Error())
			}
			return nil
		} else if statusCode >= 500 && statusCode < 600 {
			return fmt.Errorf("series groups server error: %d", statusCode)
		} else {
			return backoff.Permanent(fmt.Errorf("series groups client error: %d", statusCode))
		}
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 2 * time.Minute
	backoffConfig.MaxElapsedTime = 5 * time.Minute
	backoffConfig.Multiplier = 1.5
	backoffConfig.RandomizationFactor = 0.5

	retryPolicy := backoff.WithMaxRetries(backoffConfig, 20)

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get series groups after retries: %v, series id: %d",
			err,
			seriesId,
		)
	}
	if respBody == nil {
		return nil, fmt.Errorf("MuGetSeriesGroups: respBody is nil after successful unmarshal")
	}

	return respBody, nil
}

func MuGetGroupInfo(b *mubot.Bot, groupId int64) (*MuGroupInfoResponse, error) {
	var respBody *MuGroupInfoResponse

	operation := func() error {
		resp, statusCode, err := muGetRequest(
			"https://api.mangaupdates.com/v1/groups/"+strconv.FormatInt(groupId, 10),
			b.MuToken,
		)

		if err != nil {
			return fmt.Errorf(
				"failed to fetch group info: %s, %s, %d",
				err.Error(),
				string(resp),
				groupId,
			)
		}

		if statusCode == 200 {
			respBody = &MuGroupInfoResponse{}
			if err := json.Unmarshal(resp, respBody); err != nil {
				return fmt.Errorf("failed to unmarshal group info: %s", err.Error())
			}
			return nil
		} else if statusCode >= 500 && statusCode < 600 {
			return fmt.Errorf("series info server error: %d", statusCode)
		} else {
			return backoff.Permanent(fmt.Errorf("series info client error: %d", statusCode))
		}
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 2 * time.Minute
	backoffConfig.MaxElapsedTime = 5 * time.Minute
	backoffConfig.Multiplier = 1.5
	backoffConfig.RandomizationFactor = 0.5

	retryPolicy := backoff.WithMaxRetries(backoffConfig, 20)

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get group info after retries: %v, group id: %d",
			err,
			groupId,
		)
	}
	if respBody == nil {
		return nil, fmt.Errorf("MuGetGroupInfo: respBody is nil after successful unmarshal")
	}

	return respBody, nil
}
