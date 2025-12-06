# Wallpaper Gallery - Go Edition

This application has been migrated from PHP/MySQL to Go with JSON storage.

## Changes from PHP Version

### Backend
- **Language**: PHP → Go
- **Database**: MySQL → JSON file (`data/images.json`)
- **Server**: Apache/Nginx → Built-in Go HTTP server
- **Port**: Runs on `http://localhost:8080` by default

### Features Retained
- Image gallery with infinite scroll pagination
- Like/Dislike functionality for images
- Select multiple images and download as ZIP
- Full-size image viewer with keyboard navigation
- Thumbnail previews

## Directory Structure

```
wallpapers/
├── main.go                    # Main HTTP server and API handlers
├── models.go                  # Data models and JSON storage operations
├── init_data.go              # Utility to initialize JSON data from images
├── data/
│   └── images.json           # JSON database file
├── resource/
│   ├── images/               # Full-size images
│   ├── imagesmall/           # Thumbnail images
│   ├── css/                  # Stylesheets
│   ├── js/                   # JavaScript files (updated for Go API)
│   └── icons/                # UI icons
├── download/                 # Temporary ZIP files
└── index.html               # Main HTML page
```

## Installation & Setup

### Prerequisites
- Go 1.16 or higher
- Image files in `resource/images/` and `resource/imagesmall/` directories

### Steps

1. **Initialize data from existing images**:
   ```bash
   go run init_data.go
   ```
   This scans the `resource/images/` directory and creates `data/images.json`

   **Note:** `init_data.go` is a separate utility program. Do NOT run it together with main.go.

2. **Run the web server**:
   ```bash
   go run main.go models.go
   ```

   **Important:** You must include both `main.go` and `models.go` when running the server.

3. **Access the application**:
   Open your browser to `http://localhost:8080`

### Building Binaries

To create standalone executables:
```bash
# Build the data initializer
go build -o init_data init_data.go

# Build the web server
go build -o wallpapers main.go models.go

# Run them:
./init_data      # Initialize/update the JSON database
./wallpapers     # Start the web server
```

## API Endpoints

All API endpoints are RESTful and return appropriate HTTP status codes:

- `GET /` - Serves the main HTML page
- `GET /api/count` - Returns total number of pagination groups
- `POST /api/images` - Returns paginated HTML for images
  - Form param: `group_no` (integer)
- `POST /api/like` - Increments like counter for an image
  - Form param: `like` (image path)
- `POST /api/dislike` - Increments dislike counter for an image
  - Form param: `dislike` (image path)
- `POST /api/download` - Creates and downloads ZIP of selected images
  - Form param: `obj` (JSON string of image objects)

## JSON Data Format

The `data/images.json` file stores image metadata:

```json
[
  {
    "id": 1,
    "name": "image1.jpg",
    "liked": 5,
    "disliked": 2
  },
  {
    "id": 2,
    "name": "image2.jpg",
    "liked": 10,
    "disliked": 1
  }
]
```

## Migration Notes

### What Was Removed
- PHP files (`*.php`) are no longer needed
- MySQL database and configuration
- PHP sessions (state is now managed client-side)
- `Web.config` (IIS configuration)

### What Changed
- Frontend JavaScript updated to call Go API endpoints instead of PHP scripts
- Image pagination now uses JSON data instead of SQL queries
- Downloads are created on-demand and automatically cleaned up

### Data Migration
If you have existing data in MySQL:
1. Export image filenames from the `files` table
2. Create a JSON file manually or modify `init_data.go` to preserve like/dislike counts
3. Run the data initialization

## Performance Improvements

- **Concurrent-safe**: JSON operations use mutex locks for thread safety
- **Efficient pagination**: Filters images in memory (suitable for moderate datasets)
- **Static file caching**: Go's FileServer handles caching headers automatically
- **Background cleanup**: ZIP files are deleted asynchronously after download

## Future Enhancements

Potential improvements:
- Database (SQLite, PostgreSQL) for larger datasets
- Image upload functionality
- User authentication
- Search and filtering
- Image scraping from 4chan (currently removed)

## Troubleshooting

### Port already in use
Change the port in `main.go`:
```go
port := "8080"  // Change to another port
```

### Images not loading
- Ensure `resource/images/` and `resource/imagesmall/` directories exist
- Check file permissions (images must be readable)
- Run `go run init_data.go` to rebuild the JSON database

### JSON file locked
The application uses mutexes to prevent concurrent write conflicts. If issues persist, restart the server.

## License

Same as the original project.
