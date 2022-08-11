import React from 'react'
import { useTranslation } from 'react-i18next'

const BankAccounts = () => {
  const { t } = useTranslation('assets')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('banks.bankAccounts')}</div>
        <div className="card-body">
          <p>{t('banks.description')}</p>
        </div>
      </div>
    </>
  )
}

export default BankAccounts
