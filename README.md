# AI Chat Application

This project is an AI-powered chat application that utilizes both Go and Python to handle chat sessions and AI responses, respectively. The application stores chat history in a SQLite database and interfaces with OpenAI's GPT models for generating responses.

## Features

- Interactive chat interface with session management.
- Backend implemented in Go for handling chat sessions and database interactions.
- Python script for querying OpenAI's GPT models.
- SQLite database for persisting chat sessions and history.

## Prerequisites

To run this application, you will need:

- Go (version 1.x)
- Python (version 3.x)
- Access to OpenAI API

## Installation and Setup

1. **Clone the Repository:**


2. **Set Up the Database:**

Ensure SQLite is installed on your system. The Go application will automatically create a database file named `conversation.db`.

3. **Set Up OpenAI API Key:**

You will need to export your OpenAI API key as an environment variable:

```bash
export OPENAI_API_KEY='your_api_key_here'
```


4. **Install Python Dependencies:**
'''bash
pip install -r requirements.txt
'''


## Running the Application

1. **Start the Go Server:**

From the root directory of the project, run:

'''bash
go run main.go`
'''

This will start the server on `localhost:8080`.

2. **Accessing the Chat Interface:**

Open your browser and navigate to `http://localhost:8080` to access the chat interface.

## Usage

- Create new chat sessions, send messages, and receive AI-generated responses.
- Manage chat sessions from the sidebar in the web interface.
- View and interact with chat history for each session.

## Acknowledgements

- This project uses OpenAI's GPT models for generating AI responses.
- Built with Go and Python.







