package metube

import (
	"fmt"
	"log"
	"time"

	"github.com/imroc/req/v3"
)

type Client struct {
	URL    string `json:"url" xml:"url"`
	Debug  bool
	Client *req.Client
}

func New(url string, debug bool) *Client {
	r := req.C().
		SetBaseURL(url).
		SetTimeout(5 * time.Second)
	if debug {
		r.DevMode()
	}
	return &Client{
		URL:    url,
		Debug:  debug,
		Client: r,
	}
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

type AddRequest struct {
	URL       string `json:"url"`
	AutoStart bool   `json:"auto_start"`
	Quality   string `json:"quality"`
	Format    string `json:"format"`
	Folder    string `json:"folder"`
	Name      string `json:"custom_name_prefix"`
}

func (m *Client) Add(name, url string, autoStart bool) error {
	result := &Response{}
	resp, err := m.Client.R().
		SetHeader("Accept", "application/json"). // Chainable request settings.
		SetSuccessResult(result).                // Unmarshal response body into userInfo automatically if status code is between 200 and 299.
		SetBody(&AddRequest{url, autoStart, "best", "any", "", name}).
		Post("add")

	if err != nil { // Error handling.
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return fmt.Errorf("failed: %w", err)
	}

	if resp.IsErrorState() { // Status code >= 400.
		return fmt.Errorf("failed: %d: %s: %s", resp.StatusCode, resp.Status, resp.String())
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		if result.Status != "ok" {
			return fmt.Errorf("failed: %s", result.Message)
		}
	}

	return nil
}

type DeleteRequest struct {
	IDs   []string `json:"ids"`
	Where string   `json:"where"`
}

func (m *Client) Delete(ids []string, where string) error {
	result := &Response{}
	resp, err := m.Client.R().
		SetHeader("Accept", "application/json"). // Chainable request settings.
		SetSuccessResult(result).                // Unmarshal response body into userInfo automatically if status code is between 200 and 299.
		SetBody(&DeleteRequest{ids, where}).
		Post("delete")

	if err != nil { // Error handling.
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return fmt.Errorf("failed: %w", err)
	}

	if resp.IsErrorState() { // Status code >= 400.
		return fmt.Errorf("failed: %d: %s: %s", resp.StatusCode, resp.Status, resp.String())
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		if result.Status != "ok" {
			return fmt.Errorf("failed: %s", result.Message)
		}
	}
	return nil
}

type StartRequest struct {
	IDs []string `json:"ids"`
}

func (m *Client) Start(ids []string) error {
	result := &Response{}
	resp, err := m.Client.R().
		SetHeader("Accept", "application/json"). // Chainable request settings.
		SetSuccessResult(result).                // Unmarshal response body into userInfo automatically if status code is between 200 and 299.
		SetBody(&StartRequest{ids}).
		Post("delete")

	if err != nil { // Error handling.
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return fmt.Errorf("failed: %w", err)
	}

	if resp.IsErrorState() { // Status code >= 400.
		return fmt.Errorf("failed: %d: %s: %s", resp.StatusCode, resp.Status, resp.String())
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		if result.Status != "ok" {
			return fmt.Errorf("failed: %s", result.Message)
		}
	}
	return nil
}

type HistoryResponse struct {
	Done  []*Download `json:"done"`
	Queue []*Download `json:"queue"`
}
type Download struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	URL              string  `json:"url"`
	Quality          string  `json:"quality"`
	Format           string  `json:"format"`
	Folder           string  `json:"folder"`
	CustomNamePrefix string  `json:"custom_name_prefix"`
	Msg              string  `json:"msg"`
	Percent          float64 `json:"percent"`
	Speed            int64   `json:"speed"`
	Eta              int64   `json:"eta"`
	Status           string  `json:"status"`
	Size             int64   `json:"size"`
	Timestamp        int64   `json:"timestamp"`
	Error            string  `json:"error"`
	Filename         string  `json:"filename"`
}

func (m *Client) History() (*HistoryResponse, error) {
	result := &HistoryResponse{}
	resp, err := m.Client.R().
		SetHeader("Accept", "application/json").
		SetContentType("application/json").
		SetSuccessResult(result).
		Get("history")

	if err != nil { // Error handling.
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return nil, fmt.Errorf("failed: %w", err)
	}

	if resp.IsErrorState() { // Status code >= 400.
		return nil, fmt.Errorf("failed: %d: %s: %s", resp.StatusCode, resp.Status, resp.String())
	}

	return result, nil
}
