package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"pokabook/issue-bot/pkg/asana"
)

type createTaskBody struct {
	WorkSpace    string                 `json:"workspace"`
	Projects     []string               `json:"projects"`
	Name         string                 `json:"name"`
	Assignee     string                 `json:"assignee"`
	Notes        string                 `json:"notes"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type taskCreateResult struct {
	Data struct {
		PermalinkURL string `json:"permalink_url"`
		CustomFields []struct {
			Name      string `json:"name"`
			TextValue string `json:"text_value"`
		} `json:"custom_fields"`
	} `json:"data"`
}

func CreateTask(title, tagID, assignee, description, major, priority, label string) (string, string) {
	url := "https://app.asana.com/api/1.0/tasks"

	customFields := map[string]interface{}{
		"1206775921028667": asana.PriorityIdMap[priority],
		"1206775921028676": asana.LabelIdMap[label],
	}

	body := map[string]interface{}{
		"data": createTaskBody{
			WorkSpace:    asana.AsanaWorkspaceID,
			Projects:     []string{asana.AsanaProjectWithTag[tagID]},
			Name:         title,
			Assignee:     asana.NicknameEmailMap[assignee],
			Notes:        description,
			CustomFields: customFields,
		},
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return "", ""
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))

	req.Header.Set("Authorization", "Bearer "+asana.AsanaToken)
	req.Header.Set("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)

	createResult := taskCreateResult{}
	err = json.Unmarshal(response, &createResult)
	if err != nil {
		return "", ""
	}

	var issueID string
	for _, field := range createResult.Data.CustomFields {
		if field.Name == major {
			issueID = field.TextValue
			break
		}
	}

	return issueID, createResult.Data.PermalinkURL
}
