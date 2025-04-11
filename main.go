package main

import (
	"encoding/json"
	"fmt"

	pdk "github.com/extism/go-pdk"
	"github.com/tidwall/gjson"
)

type HNStory struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
	By    string `json:"by"`
}

func Call(input CallToolRequest) (CallToolResult, error) {
	numStories := 10 // default value
	if args, ok := input.Params.Arguments.(map[string]interface{}); ok {
		if n, exists := args["num_stories"]; exists {
			if val, ok := n.(float64); ok {
				numStories = int(val)
				if numStories > 100 {
					numStories = 100
				}
				if numStories < 1 {
					numStories = 10
				}
			}
		}
	}
	return fetchHackerNews(numStories)
}

func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "hackernews",
				Description: "Get top stories from Hacker News (max 100)",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"num_stories": map[string]interface{}{
							"description": "Number of top stories to fetch (max 100, defaults to 10)",
							"type":        "integer",
							"minimum":     1,
							"maximum":     100,
							"default":     10,
						},
					},
					"required": []string{"num_stories"},
				},
			},
		},
	}, nil
}

func fetchHackerNews(numStories int) (CallToolResult, error) {
	// Fetch top stories
	req := pdk.NewHTTPRequest(pdk.MethodGet, "https://hacker-news.firebaseio.com/v0/topstories.json")
	resp := req.Send()

	if resp.Status() != 200 {
		return CallToolResult{}, fmt.Errorf("failed to fetch top stories: HTTP %d", resp.Status())
	}

	var storyIDs []int
	if err := json.Unmarshal(resp.Body(), &storyIDs); err != nil {
		return CallToolResult{}, fmt.Errorf("failed to parse story IDs: %v", err)
	}

	// Get requested number of stories
	var stories []HNStory
	for i := 0; i < numStories && i < len(storyIDs); i++ {
		story, err := fetchStory(storyIDs[i])
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to fetch story %d: %v", storyIDs[i], err)
		}
		stories = append(stories, story)
	}

	// Convert stories to JSON string
	jsonBytes, err := json.Marshal(stories)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to marshal stories: %v", err)
	}
	text := string(jsonBytes)

	return CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &text,
			},
		},
	}, nil
}

func fetchStory(id int) (HNStory, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	resp := req.Send()

	if resp.Status() != 200 {
		return HNStory{}, fmt.Errorf("failed to fetch story: HTTP %d", resp.Status())
	}

	// Parse story data
	body := resp.Body()
	story := HNStory{
		ID:    int(gjson.GetBytes(body, "id").Int()),
		Title: gjson.GetBytes(body, "title").String(),
		URL:   gjson.GetBytes(body, "url").String(),
		Score: int(gjson.GetBytes(body, "score").Int()),
		By:    gjson.GetBytes(body, "by").String(),
	}

	return story, nil
}

//export call
func _Call() int32 {
	var err error
	_ = err

	input := pdk.InputString()
	var request CallToolRequest
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		pdk.SetError(err)
		return -1
	}

	output, err := Call(request)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.OutputString(string(jsonBytes))
	return 0
}

//export describe
func _Describe() int32 {
	var err error
	_ = err

	output, err := Describe()
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	pdk.OutputString(string(jsonBytes))
	return 0
}

type CallToolRequest struct {
	Method *string `json:"method,omitempty"`
	Params Params  `json:"params"`
}

type CallToolResult struct {
	Content []Content `json:"content"`
	IsError *bool     `json:"isError,omitempty"`
}

type Content struct {
	Type string  `json:"type"`
	Text *string `json:"text,omitempty"`
}

type Params struct {
	Arguments interface{} `json:"arguments,omitempty"`
	Name      string      `json:"name"`
}

type ListToolsResult struct {
	Tools []ToolDescription `json:"tools"`
}

type ToolDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

const ContentTypeText = "text"

func main() {}
