# Text-to-Speech Web Application

This repository contains the source code for a web-based Text-to-Speech (TTS) application. The application allows users to input text and generate spoken audio using OpenAI's TTS service.

## Features

- Convert text to speech and play it directly in the browser.
- Save TTS requests and access them later.
- List all TTS records with playback functionality.

## Prerequisites

Before running this application, make sure you have the following installed:
- Go (version 1.21.0 or higher)
- PostgreSQL

Additionally, you will need to set up environment variables for database connection and OpenAI API key.

## Installation

To get started with this project, clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/tts-web.git
cd tts-web
```

### Install the required Go modules:
```bash
go mod tidy
```

## Configuration
Set the following environment variables for your database connection and OpenAI API key:

```bash
export DB_HOST="your_database_host"
export DB_PORT="your_database_port"
export DB_USER="your_database_user"
export DB_PASSWORD="your_database_password"
export TTS_DB_NAME="your_database_name"
export TTS_API_KEY="your_openai_api_key"
export TTS_PORT="8081" # Optional: default port is 8081
```

## Running the Application
```angular2html
go run .
```

The application will start and be accessible at http://localhost:8081 or another port if you specified one in the environment variables.

## Endpoints
```bash
GET /: Serves the homepage where you can submit new TTS requests.
POST /api/tts: Endpoint to create a new TTS request.
GET /api/tts/records: Endpoint to list all TTS records.
GET /api/tts/records/{recordId}: Endpoint to get details of a specific TTS record.
GET /api/tts/records/{recordId}/audio: Endpoint to get the audio content of a specific TTS record.
```

## Static Files
Static files such as CSS and client-side JavaScript are located in the static directory.

## Templates
HTML templates are located in the templates directory.

## License
This project is licensed under the [MIT LICENSE.md](LICENSE.md) - see the LICENSE file for details.
