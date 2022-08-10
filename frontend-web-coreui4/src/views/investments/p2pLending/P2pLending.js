import React from 'react'
import { useTranslation } from 'react-i18next'

const P2pLending = () => {
  const { t } = useTranslation('investments')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('p2pLending.p2pLending')}</div>
        <div className="card-body">
          <p>{t('p2pLending.description')}</p>
        </div>
      </div>
    </>
  )
}

export default P2pLending
