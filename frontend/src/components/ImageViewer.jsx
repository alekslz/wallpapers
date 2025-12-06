import { useEffect } from 'react'
import './ImageViewer.css'

const ImageViewer = ({ image, onClose, isSelected, onToggleSelection }) => {
  useEffect(() => {
    const handleKeyPress = (e) => {
      if (e.key === 'Escape') {
        onClose()
      }
    }

    window.addEventListener('keyup', handleKeyPress)
    return () => window.removeEventListener('keyup', handleKeyPress)
  }, [onClose])

  const handleLike = async () => {
    try {
      const formData = new FormData()
      formData.append('like', `/resource/images/${image.name}`)
      await fetch('/api/like', {
        method: 'POST',
        body: formData
      })
    } catch (err) {
      console.error('Error liking image:', err)
    }
  }

  const handleDislike = async () => {
    try {
      const formData = new FormData()
      formData.append('dislike', `/resource/images/${image.name}`)
      await fetch('/api/dislike', {
        method: 'POST',
        body: formData
      })
    } catch (err) {
      console.error('Error disliking image:', err)
    }
  }

  return (
    <div className="image-viewer-overlay" onClick={onClose}>
      <div className="image-viewer" onClick={(e) => e.stopPropagation()}>
        <button className="close-btn" onClick={onClose}>×</button>

        <img
          src={`/resource/images/${image.name}`}
          alt={image.name}
          className="full-image"
        />

        <div className="viewer-controls">
          <label className="checkbox-label">
            <input
              type="checkbox"
              checked={isSelected}
              onChange={onToggleSelection}
            />
            <span>P.F.D</span>
          </label>
          <button onClick={handleLike} className="like-btn">Like</button>
          <button onClick={handleDislike} className="dislike-btn">Dislike</button>
        </div>
      </div>
    </div>
  )
}

export default ImageViewer
