# StormSync Provider

![License](https://img.shields.io/github/license/stormsync/provider)
![Issues](https://img.shields.io/github/issues/stormsync/provider)
![Forks](https://img.shields.io/github/forks/stormsync/provider)
![Stars](https://img.shields.io/github/stars/stormsync/provider)

StormSync Provider is a Go-based library designed to facilitate data synchronization and management. This project includes components for API interactions, data consumption, and storage handling.
---

[Storm Sync API Documentation](https://app.swaggerhub.com/apis-docs/StormSync/stormsync/v1.0.0)
Based on the information from the repository, here is a detailed README for the StormSync Provider:


## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Installation

To get started with StormSync Provider, clone the repository and build the project using Go:

```bash
git clone https://github.com/stormsync/provider.git
cd provider
go build
```

## Usage

### Running the Server

To run the server, navigate to the `cmd/server` directory and execute:

```bash
cd cmd/server
go build -o app main.go
chmod +x app
./app
```

This will start the StormSync server, which can then be accessed via the configured API endpoints.

## Configuration

Configuration settings for the StormSync Provider are typically defined in environment variables. Example configuration includes setting up the database connection, API keys, and other needed settings.


### Environment Variables
These env vars are required to run this application.
```bash
API_KEY="xxxx" 
DB_PASS="xxxx" 
DB_ADDRESS="xxxx" 
DB_NAME="xxxxx" 
DB_USER="xxxxxx" 
WEB_SERVER_ADDRESS="0.0.0.0:8080", 
KAFKA_ADDRESS="xxxxx"
KAFKA_USER="xxxxxx"
KAFKA_PASSWORD="xxxxxx"
CONSUMER_TOPIC="transformed-weather-data"  
```


### Running Tests

Run tests using the following command:

```bash
go test ./...
```

---

