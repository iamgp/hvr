# Hamilton Venus Registry

Hamilton Venus Registry is a package management system for storing and retrieving libraries.

## Features

- Upload libraries with metadata (name, version, description, author, repository URL)
- Download libraries (specific version or latest)
- Search for libraries
- File integrity verification using SHA-256 hash
- Metadata-based uploads using JSON files
- Preservation of file modification times

## Installation

1. Clone the repository:

   ```
   git clone github.com/iamgp/hvr.git
   ```

2. Build the server and client:
   ```
   make all
   ```

## Usage

### Starting the Server

Run the server using:

```
make run
```

The server will start on `localhost:8080`.

### Using the CLI Client

The CLI client provides commands to interact with the server:

1. Upload a library:

   ```
   ./hvr upload <file> --name <library-name> --version <version>
   ```

2. Upload a library using a metadata file:

   ```
   ./hvr uploadmeta <metadata-file>
   ```

   Example metadata file (library_meta.json):

   ```json
   {
     "name": "my-library",
     "version": "1.0.0",
     "description": "A useful library",
     "author": "John Doe",
     "repo_url": "https://github.com/johndoe/my-library",
     "files": ["src/*.go", "README.md", "LICENSE"]
   }
   ```

3. Download a library:

   ```
   ./hvr download <library-name> [version]
   ```

   If version is omitted, it will download the latest version.

4. Search for libraries:
   ```
   ./hvr search <query>
   ```

## How It Works

1. **Server**: The server uses an SQLite database to store library information and a local file system to store library files. It provides HTTP endpoints for uploading, downloading, and searching libraries.

2. **Client**: The CLI client sends HTTP requests to the server to perform operations.

3. **Upload**: When a library is uploaded, it's stored in the database with its name, version, description, author, repository URL, file path, and a SHA-256 hash of the file contents.

4. **Download**: When a library is downloaded, the server retrieves the file from storage, and the client verifies the file integrity using the stored hash.

5. **Search**: The search functionality allows users to find libraries by name.

6. **Metadata Upload**: Users can provide a JSON metadata file that specifies multiple files to be included in the library, along with other metadata.

7. **File Integrity**: SHA-256 hashes are used to ensure the integrity of downloaded files.

8. **Modification Time**: The original modification time of uploaded files is preserved and restored upon download.

## Development

- Reset the database:
  ```
  make reset-db
  ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
