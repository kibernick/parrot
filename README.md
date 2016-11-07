# parrot

A simple API that accepts SMS messages submitted via a POST request. Uses the MessageBird API to sent proper SMS messages. It _parrots_ the Concatenated SMS functionality by splitting long messages (>160 characters) into parts and setting the correct message headers.

## Installation

1. `go get...` & `go install`
2. Create an account on MessageBird. Make a note your API key.
3. Create your `config.json` based on the `config_example.json` file, and add your API key there.
4. Run the server by executing the Go binary.
5. The default API endpoint will be available at: `http://127.0.0.1:8000`

## Usage

Example JSON payload:
```json
{
    "recipient": "123", // a valid phone number
    "originator": "123", // another valid phone number
    "message": "Hello, stranger!"
}
```
