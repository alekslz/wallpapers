# Wallpaper Gallery - Modern Stack

A modern wallpaper gallery application with Go backend and React frontend.

## Tech Stack

### Backend
- **Language**: Go
- **Database**: JSON file (`data/images.json`)
- **Server**: Built-in Go HTTP server
- **Port**: `http://localhost:8080`

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **Styling**: CSS Modules

### Features
- Image gallery with infinite scroll pagination
- Like/Dislike functionality for images
- Select multiple images and download as ZIP
- Full-size image viewer with keyboard navigation
- Thumbnail previews
- **4chan wallpaper scraper** (re-implemented in Go)

## Directory Structure

```
wallpapers/
├── main.go                    # Main HTTP server and API handlers
├── models.go                  # Data models and JSON storage operations
├── init_data.go              # Utility to initialize JSON data from images
├── scraper.go                # Utility to download wallpapers from 4chan /wg/
├── create_thumbnails.go      # Utility to generate thumbnail images
├── frontend/                 # React frontend (Vite)
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── App.jsx          # Main app component
│   │   └── main.jsx         # Entry point
│   ├── package.json
│   └── vite.config.js
├── dist/                     # Built frontend (generated)
├── data/
│   └── images.json           # JSON database file
├── resource/
│   ├── images/               # Full-size images
│   └── imagesmall/           # Thumbnail images
└── download/                 # Temporary ZIP files
```

## Utilities

The application includes several utility programs:

1. **scraper.go** - Downloads wallpapers from 4chan.org/wg/
   - Scrapes the first 10 pages of the board
   - Downloads up to 15 threads worth of images
   - Skips images that already exist
   - Rate-limited to be respectful to 4chan's servers

2. **create_thumbnails.go** - Generates thumbnail images
   - Creates 300x300 max thumbnails maintaining aspect ratio
   - High-quality CatmullRom scaling
   - Skips thumbnails that already exist

3. **init_data.go** - Initializes the JSON database
   - Scans `resource/images/` directory
   - Creates metadata for each image in `data/images.json`
   - Sets initial like/dislike counts to 0

## Installation & Setup

### Prerequisites
- Go 1.16 or higher
- Node.js 16 or higher
- npm or yarn

### Quick Start

#### Option A: Download Wallpapers from 4chan

1. **Scrape and download wallpapers**:
   ```bash
   go run scraper.go
   ```

2. **Create thumbnails**:
   ```bash
   go run create_thumbnails.go
   ```

3. **Initialize the database**:
   ```bash
   go run init_data.go
   ```

4. **Build the frontend**:
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

5. **Start the server**:
   ```bash
   go run main.go models.go
   ```

6. **Access the app**: http://localhost:8080

#### Option B: Use Existing Images

1. **Add images** to `resource/images/` and `resource/imagesmall/`

2. **Initialize the database**:
   ```bash
   go run init_data.go
   ```

3. **Build the frontend**:
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

4. **Start the server**:
   ```bash
   go run main.go models.go
   ```

5. **Access the app**: http://localhost:8080

### Development Mode

For development with hot module replacement:

1. **Start the Go backend**:
   ```bash
   go run main.go models.go
   ```

2. **In a new terminal, start the Vite dev server**:
   ```bash
   cd frontend
   npm run dev
   ```

3. **Access the development server**: http://localhost:5173

The Vite dev server will proxy API requests to the Go backend at :8080.

### Building Binaries

To create standalone executables:
```bash
# Build all utilities
go build -o scraper scraper.go
go build -o create_thumbnails create_thumbnails.go
go build -o init_data init_data.go
go build -o wallpapers main.go models.go

# Run them in order:
./scraper            # Download wallpapers from 4chan
./create_thumbnails  # Create thumbnail images
./init_data          # Initialize/update the JSON database
./wallpapers         # Start the web server
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
