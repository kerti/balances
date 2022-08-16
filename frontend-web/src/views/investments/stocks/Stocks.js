import { useTranslation } from 'react-i18next'

const Stocks = () => {
  const { t } = useTranslation('investments')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('stocks.stocks')}</div>
        <div className="card-body">
          <p>{t('stocks.description')}</p>
        </div>
      </div>
    </>
  )
}

export default Stocks
