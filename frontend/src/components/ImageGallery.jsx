import { useState, useEffect, useRef, useCallback } from 'react'
import ImageCard from './ImageCard'
import './ImageGallery.css'

const ImageGallery = ({ selectedImages, onToggleSelection, onImageClick }) => {
  const [images, setImages] = useState([])
  const [groupNo, setGroupNo] = useState(0)
  const [totalGroups, setTotalGroups] = useState(0)
  const [loading, setLoading] = useState(false)
  const observerTarget = useRef(null)

  // Fetch total count on mount
  useEffect(() => {
    fetch('/api/count')
      .then(res => res.json())
      .then(total => {
        setTotalGroups(total)
        loadImages(0)
      })
      .catch(err => console.error('Error fetching count:', err))
  }, [])

  const loadImages = useCallback(async (group) => {
    if (loading) return
    setLoading(true)

    try {
      const formData = new FormData()
      formData.append('group_no', group)

      const response = await fetch('/api/images', {
        method: 'POST',
        body: formData
      })

      const html = await response.text()

      // Parse HTML to extract image names
      const parser = new DOMParser()
      const doc = parser.parseFromString(html, 'text/html')
      const imgElements = doc.querySelectorAll('img.img')

      const newImages = Array.from(imgElements).map((img, index) => ({
        id: group * 70 + index,
        name: img.src.split('/').pop(),
        src: img.src
      }))

      setImages(prev => [...prev, ...newImages])
      setGroupNo(group + 1)
    } catch (err) {
      console.error('Error loading images:', err)
    } finally {
      setLoading(false)
    }
  }, [loading])

  // Infinite scroll observer
  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting && groupNo <= totalGroups && !loading) {
          loadImages(groupNo)
        }
      },
      { threshold: 1.0 }
    )

    if (observerTarget.current) {
      observer.observe(observerTarget.current)
    }

    return () => {
      if (observerTarget.current) {
        observer.unobserve(observerTarget.current)
      }
    }
  }, [groupNo, totalGroups, loading, loadImages])

  return (
    <div className="image-gallery">
      <div className="images-grid">
        {images.map(image => (
          <ImageCard
            key={image.id}
            image={image}
            isSelected={!!selectedImages[image.name]}
            onToggleSelection={() => onToggleSelection(image)}
            onClick={() => onImageClick(image)}
          />
        ))}
      </div>
      <div ref={observerTarget} className="observer-target" />
      {loading && <div className="loading">Loading...</div>}
    </div>
  )
}

export default ImageGallery
