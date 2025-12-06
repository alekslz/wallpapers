import './DownloadButton.css'

const DownloadButton = ({ selectedImages, onClear }) => {
  const handleDownload = async () => {
    try {
      const formData = new FormData()
      formData.append('obj', JSON.stringify(selectedImages))

      const response = await fetch('/api/download', {
        method: 'POST',
        body: formData
      })

      if (!response.ok) throw new Error('Download failed')

      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `wallpapers-${Date.now()}.zip`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)

      onClear()
    } catch (err) {
      console.error('Error downloading images:', err)
      alert('Failed to download images')
    }
  }

  const count = Object.keys(selectedImages).length

  return (
    <div className="download-container">
      <div className="selected-count">{count} image{count !== 1 ? 's' : ''} selected</div>
      <button onClick={handleDownload} className="download-btn">
        Download ZIP
      </button>
      <button onClick={onClear} className="clear-btn">Clear Selection</button>
    </div>
  )
}

export default DownloadButton
