package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

var usernames = []string{
	"mathewka",
	"Krish06m",
}

var db *sql.DB

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

func main() {
	var err error
	db, err = sql.Open("sqlite3", "lc.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
	startScheduler()

	http.HandleFunc("/history", withCORS(historyHandler))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS daily_stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		timestamp TEXT,
		easy INTEGER,
		medium INTEGER,
		hard INTEGER,
		total INTEGER
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func startScheduler() {
	c := cron.New()

	// Runs once every day at midnight
	c.AddFunc("* */2 * * *", func() {
		log.Println("Daily update running...")
		updateAllUsers()
	})

	c.Start()

	// Run once at startup
	updateAllUsers()
}

func updateAllUsers() {
	for _, username := range usernames {
		err := fetchAndStore(username)
		if err != nil {
			log.Println("Error:", err)
		}
		time.Sleep(2 * time.Second)
	}
}

func fetchAndStore(username string) error {
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

	req, _ := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result LeetCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	var easy, medium, hard int

	for _, stat := range result.Data.MatchedUser.SubmitStats.AcSubmissionNum {
		switch stat.Difficulty {
		case "Easy":
			easy = stat.Count
		case "Medium":
			medium = stat.Count
		case "Hard":
			hard = stat.Count
		}
	}

	total := easy + medium + hard

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(`
	INSERT INTO daily_stats (username, timestamp, easy, medium, hard, total)
	VALUES (?, ?, ?, ?, ?, ?)
	`,
		username, timestamp, easy, medium, hard, total,
	)

	return err
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
	SELECT username, timestamp, total
	FROM daily_stats
	ORDER BY timestamp ASC
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	result := make(map[string][]map[string]interface{})

	for rows.Next() {
		var username, timestamp string
		var total int
		rows.Scan(&username, &timestamp, &total)

		result[username] = append(result[username], map[string]interface{}{
			"timestamp":  timestamp,
			"total": total,
		})
	}

	json.NewEncoder(w).Encode(result)
}
