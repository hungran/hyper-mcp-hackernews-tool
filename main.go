package main

import (
	"encoding/json"
	"fmt"
	"strings"

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
	return fetchHackerNews()
}

func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "hackernews",
				Description: "Get top 5 stories from Hacker News",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"random_string": map[string]interface{}{
							"description": "Dummy parameter for no-parameter tools",
							"type":        "string",
						},
					},
					"required": []string{"random_string"},
				},
			},
		},
	}, nil
}

func fetchHackerNews() (CallToolResult, error) {
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

	// Get top 5 stories
	var stories []HNStory
	for i := 0; i < 5 && i < len(storyIDs); i++ {
		story, err := fetchStory(storyIDs[i])
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to fetch story %d: %v", storyIDs[i], err)
		}
		stories = append(stories, story)
	}

	// Format output
	var output strings.Builder
	output.WriteString("# Top 5 Hacker News Stories\n\n")
	for i, story := range stories {
		output.WriteString(fmt.Sprintf("## %d. %s\n", i+1, story.Title))
		output.WriteString(fmt.Sprintf("Score: %d | Author: %s\n", story.Score, story.By))
		if story.URL != "" {
			output.WriteString(fmt.Sprintf("URL: %s\n", story.URL))
		}
		output.WriteString("\n")
	}

	text := output.String()
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
