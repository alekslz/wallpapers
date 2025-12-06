import { useState } from 'react'
import ImageGallery from './components/ImageGallery'
import ImageViewer from './components/ImageViewer'
import DownloadButton from './components/DownloadButton'
import './App.css'

function App() {
  const [selectedImages, setSelectedImages] = useState({})
  const [viewerImage, setViewerImage] = useState(null)

  const toggleImageSelection = (image) => {
    setSelectedImages(prev => {
      const newSelection = { ...prev }
      if (newSelection[image.name]) {
        delete newSelection[image.name]
      } else {
        newSelection[image.name] = {
          src: `/resource/images/${image.name}`
        }
      }
      return newSelection
    })
  }

  const openViewer = (image) => {
    setViewerImage(image)
  }

  const closeViewer = () => {
    setViewerImage(null)
  }

  return (
    <div className="app">
      <div className="container">
        <ImageGallery
          selectedImages={selectedImages}
          onToggleSelection={toggleImageSelection}
          onImageClick={openViewer}
        />

        {Object.keys(selectedImages).length > 0 && (
          <DownloadButton
            selectedImages={selectedImages}
            onClear={() => setSelectedImages({})}
          />
        )}

        {viewerImage && (
          <ImageViewer
            image={viewerImage}
            onClose={closeViewer}
            isSelected={!!selectedImages[viewerImage.name]}
            onToggleSelection={() => toggleImageSelection(viewerImage)}
          />
        )}
      </div>
    </div>
  )
}

export default App
