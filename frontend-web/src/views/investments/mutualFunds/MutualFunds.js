import React from 'react'
import { useTranslation } from 'react-i18next'

const MutualFunds = () => {
  const { t } = useTranslation('investments')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('mutualFunds.mutualFunds')}</div>
        <div className="card-body">
          <p>{t('mutualFunds.description')}</p>
        </div>
      </div>
    </>
  )
}

export default MutualFunds
