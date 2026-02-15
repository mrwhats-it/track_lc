package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var usernames = []string{
	"mathewka", 
	"Krish06m",
	"Mann_Mehta",

}

type LeetCodeResponse struct {
	Data struct {
		MatchedUser struct {
			SubmitStats struct {
				AcSubmissionNum []struct {
					Difficulty string `json:"difficulty"`
					Count      int    `json:"count"`
				} `json:"acSubmissionNum"`
			} `json:"submitStats"`
		} `json:"matchedUser"`
	} `json:"data"`
}

func fetchUser(username string) (int, error) {
	query := `
	query getUserProfile($username: String!) {
		matchedUser(username: $username) {
			submitStats {
				acSubmissionNum {
					difficulty
					count
				}
			}
		}
	}`

	body := map[string]interface{}{
		"query": query,
		"variables": map[string]string{
			"username": username,
		},
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(
		"https://leetcode.com/graphql",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	responseBytes, _ := io.ReadAll(resp.Body)

	var result LeetCodeResponse
	err = json.Unmarshal(responseBytes, &result)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, stat := range result.Data.MatchedUser.SubmitStats.AcSubmissionNum {
		if stat.Difficulty == "Easy" ||
			stat.Difficulty == "Medium" ||
			stat.Difficulty == "Hard" {
			total += stat.Count
		}
	}

	return total, nil
}

const dataPath = "frontend/public/data.json"

func main() {
	file, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string][]map[string]interface{}
	json.Unmarshal(file, &data)

	if data == nil {
		data = make(map[string][]map[string]interface{})
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)

	for _, user := range usernames {
		total, err := fetchUser(user)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		entry := map[string]interface{}{
			"timestamp": timestamp,
			"total":     total,
		}

		data[user] = append(data[user], entry)
	}

	output, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(dataPath, output, 0644)
}
