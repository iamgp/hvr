# Hamilton Venus Registry

Hamilton Venus Registry is a package management system for Hamilton Venus library files, similar to npm or PyPI. It allows users to upload, download, and search for libraries.

## Project Structure

The project is structured as follows:

- `cmd/server/`: Contains the main server application.
- `internal/`: Contains the core logic of the application.
  - `api/handlers/`: HTTP handlers for the server.
  - `models/`: Data models used in the application.
  - `services/`: Business logic for managing libraries.
  - `storage/`: Database interactions.
- `pkg/client/`: Contains the CLI client for interacting with the server.

## Getting Started

### Prerequisites

- Go 1.16 or later
- SQLite3

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/your-username/hamilton-venus-registry.git
   cd hamilton-venus-registry
   ```

2. Install dependencies:

   ```
   go mod download
   ```

3. Build the server:

   ```
   go build -o hvr-server ./cmd/server
   ```

4. Build the client:
   ```
   go build -o hvr ./pkg/client
   ```

### Running the Server

Run the server with:

```
./hvr-server
```

The server will start on `localhost:8080`.

### Using the CLI Client

The CLI client provides commands to interact with the server:

1. Upload a library:

   ```
   ./hvr upload path/to/library.zip --name my-library --version 1.0.0
   ```

2. Download a library:

   ```
   ./hvr download my-library 1.0.0
   ```

3. Search for libraries:
   ```
   ./hvr search query
   ```

## How It Works

1. **Server**: The server uses an SQLite database to store library information and file data. It provides HTTP endpoints for uploading, downloading, and searching libraries.

2. **Client**: The CLI client sends HTTP requests to the server to perform operations.

3. **Upload**: When a library is uploaded, it's stored in the database with its name, version, and file data.

4. **Download**: When a library is downloaded, the server retrieves the file data from the database and sends it to the client.

5. **Search**: The search functionality allows users to find libraries by name.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
