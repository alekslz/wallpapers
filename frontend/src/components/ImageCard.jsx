import { useState } from 'react'
import './ImageCard.css'

const ImageCard = ({ image, isSelected, onToggleSelection, onClick }) => {
  const [hovering, setHovering] = useState(false)

  const handleLike = async (e) => {
    e.stopPropagation()
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

  const handleDislike = async (e) => {
    e.stopPropagation()
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

  const handleCheckboxChange = (e) => {
    e.stopPropagation()
    onToggleSelection()
  }

  return (
    <div
      className="image-card"
      onMouseEnter={() => setHovering(true)}
      onMouseLeave={() => setHovering(false)}
    >
      <img
        src={`/resource/imagesmall/${image.name}`}
        alt={image.name}
        onClick={onClick}
        className="thumbnail"
      />

      {hovering && (
        <div className="controls">
          <label className="checkbox-label">
            <input
              type="checkbox"
              checked={isSelected}
              onChange={handleCheckboxChange}
            />
            <span>Select</span>
          </label>
          <button onClick={handleLike} className="like-btn">
            <span>❤️</span>
            <span>Love</span>
          </button>
          <button onClick={handleDislike} className="dislike-btn">
            <span>👎</span>
            <span>Pass</span>
          </button>
        </div>
      )}
    </div>
  )
}

export default ImageCard
