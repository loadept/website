const formatDate = (utcDate: string, format: string) => {
  switch (format) {
    case "long":
      return new Date(utcDate).toLocaleString('es-PE', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    case "numeric":
      return new Date(utcDate).toLocaleString('es-PE', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    default:
      return new Date(utcDate).toLocaleDateString()
}
}

export default formatDate
