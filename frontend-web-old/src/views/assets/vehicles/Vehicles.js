import React from 'react'
import { useTranslation } from 'react-i18next'

const Vehicles = () => {
  const { t } = useTranslation('assets')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('vehicles.vehicles')}</div>
        <div className="card-body">
          <p>{t('vehicles.description')}</p>
        </div>
      </div>
    </>
  )
}

export default Vehicles
