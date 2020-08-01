import React from 'react'

const CardSpinner = () => {
  return (
    <>
      <div className="pt-3 text-center">
        <div className="spinner-border text-primary" role="status">
          <span className="sr-only">Loading...</span>
        </div>
      </div>
    </>
  )
}

export default CardSpinner
