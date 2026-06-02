# OpenAI API Request Example in Go

This project demonstrates how to send a request to an OpenAI language model using the REST API and display the response in the console.

## Prerequisites

* Go 1.22 or newer
* OpenAI API key

## Installation

Clone the repository and install dependencies:

```bash
go mod tidy
```

Create a `.env` file in the project root:

```text
OPENAI_API_KEY=your_api_key_here
```

## Running the Application

```bash
go run main.go
```

## What the Application Does

1. Loads the API key from the `.env` file.
2. Sends an HTTP POST request to the OpenAI API.
3. Requests a response from the language model.
4. Prints the API response to the console.

## Example Prompt

```text
Explain what REST API with constraints or without.
```

## Expected Output

```text
Status: 200 OK
...
```
