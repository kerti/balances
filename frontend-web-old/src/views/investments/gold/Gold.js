import React from 'react'
import { useTranslation } from 'react-i18next'

const Gold = () => {
  const { t } = useTranslation('investments')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('gold.gold')}</div>
        <div className="card-body">
          <p>{t('gold.description')}</p>
        </div>
      </div>
    </>
  )
}

export default Gold
