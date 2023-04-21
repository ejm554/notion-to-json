# Notion Database Exporter

This script allows you to export a Notion database to a JSON file. You can list all databases in your Notion workspace and select the one you want to export. The exported JSON file will be saved in the `exported` directory.

## Prerequisites

- Go 1.11 or higher
- Notion API key

## Getting Started

1. Clone the repository:
```
git clone https://github.com/yourusername/yourrepository.git
cd yourrepository
```

2. Create a `config.json` file in the project directory with the following content, replacing `your_notion_api_key` with your actual Notion API key:
```json
{
  "apiKey": "your_notion_api_key"
}
```

3. Build the executable:
```
go build main.go
```

4. Run the script:

On Unix-based systems (Linux and macOS):
```
./main
```

On Windows:
```
main.exe
```

5. Follow the prompts to select a database and export it to a JSON file.

# How it Works
1. The script loads the Notion API key from the config.json file.
2. It initializes the Notion client using the API key.
3. It retrieves the list of databases from the Notion API.
4. It prints the names of the databases and prompts the user to select one.
5. It exports the selected database to a JSON file in the exported directory. If the file already exists, it creates a new version of the file, such as playable_spaceships_v2.json.

## Troubleshooting
If you encounter any issues, make sure you have the correct Notion API key in your config.json file and that your API key has the necessary permissions for the databases you want to access.

Note that Wiki pages aren't supported as of today due to dependency not recognizing the "verification" property.
