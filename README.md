# fliQt HR System

## Development Guide

### Prerequisites

- Go (1.23)
- Make
- Docker

## Setup Instructions

1. Build and Run the Application:
    - Ensure Docker is installed and running.
    - Build and start the containers:
      ```bash
      make up
      ```
    - To stop and remove the containers:
      ```bash
      make down
      ```

2. Access the App:
    - Backend API: `http://localhost:8080`


## API Documentation

The API specification is defined using OpenAPI 3.0 and can be found in:
```
api/api.yaml
```

## Getting Started

1. Install dependencies:
```bash
make dependencies
```

2. Making API Changes

After making changes to `api/api.yaml`, you need to regenerate the API structures:
```bash
make generate
```
This will update the request and response body structures based on your OpenAPI definitions.

## Important Notes

- Always run `make generate` after modifying the API specification
- Commit both the API specification and generated code changes
