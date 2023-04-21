package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jomei/notionapi"
)

type Config struct {
	APIKey notionapi.Token `json:"apiKey"`
}

func main() {
	// Load the configuration from the config.json file.
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the Notion client.
	client := notionapi.NewClient(config.APIKey)

	// Retrieve the list of databases from the Notion API.
	databases, err := listDatabases(client)
	if err != nil {
		log.Fatalf("Error retrieving databases: %v", err)
	}

	// Print the names of the databases and prompt the user to select one.
	fmt.Println("Databases:")
	for i, database := range databases {
		title := extractPlainTextTitle(database.Title)
		fmt.Printf("%d. %s (ID: %s)\n", i+1, title, database.ID)
	}
	fmt.Print("Select a database to save to a JSON file (enter the number): ")

	// Read the user's selection.
	var selection int
	_, err = fmt.Scan(&selection)
	if err != nil {
		log.Fatalf("Error reading selection: %v", err)
	}

	// Check if the selection is valid.
	if selection < 1 || selection > len(databases) {
		log.Fatalf("Invalid selection: %d. Please enter a number between 1 and %d.", selection, len(databases))
	}

	// Retrieve the selected database.
	selectedDatabase := databases[selection-1]
	plainTextTitle := extractPlainTextTitle(selectedDatabase.Title)

	// Ensure the "exported" directory exists.
	err = os.MkdirAll("exported", 0755)
	if err != nil {
		log.Fatalf("Error creating 'exported' directory: %v", err)
	}

	// Generate the file name.
	baseFileName := strings.ReplaceAll(strings.ToLower(plainTextTitle), " ", "_")
	fileName := "exported/" + baseFileName + ".json"

	// Check if the file already exists and generate a new version if necessary.
	version := 1
	for {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			break
		}
		version++
		fileName = fmt.Sprintf("exported/%s_v%d.json", baseFileName, version)
	}

	// Create the JSON file.
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Write the database information to the JSON file.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(selectedDatabase)
	if err != nil {
		log.Fatalf("Error writing JSON data: %v", err)
	}

	fmt.Printf("Saved database '%s' to file: %s\n", plainTextTitle, fileName)
}

func listDatabases(client *notionapi.Client) ([]*notionapi.Database, error) {
	var databases []*notionapi.Database
	var cursor notionapi.Cursor

	for {
		response, err := client.Search.Do(context.Background(), &notionapi.SearchRequest{
			Filter: notionapi.SearchFilter{
				Property: "object",
				Value:    "database",
			},
			PageSize:    100,
			StartCursor: cursor,
		})
		if err != nil {
			return nil, err
		}

		for _, result := range response.Results {
			if database, ok := result.(*notionapi.Database); ok {
				databases = append(databases, database)
			}
		}

		if !response.HasMore {
			break
		}
		cursor = response.NextCursor
	}

	return databases, nil
}

func loadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// extractPlainTextTitle extracts the plain text title from a RichText array.
func extractPlainTextTitle(title []notionapi.RichText) string {
	var plainTextTitle string
	for _, richText := range title {
		plainTextTitle += richText.PlainText
	}
	return plainTextTitle
}
