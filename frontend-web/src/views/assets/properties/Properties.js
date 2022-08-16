import { useTranslation } from 'react-i18next'

const Properties = () => {
  const { t } = useTranslation('assets')
  return (
    <>
      <div className="card">
        <div className="card-header">{t('properties.properties')}</div>
        <div className="card-body">
          <p>{t('properties.description')}</p>
        </div>
      </div>
    </>
  )
}

export default Properties
