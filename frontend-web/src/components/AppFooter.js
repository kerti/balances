import React from 'react'
import { CFooter } from '@coreui/react'
import CIcon from '@coreui/icons-react'
import { cibGithub } from '@coreui/icons'

const AppFooter = () => {
  return (
    <CFooter>
      <div className="small">
        <span className="ms-1">&copy; 2022 Raditya Kertiyasa</span>
        <span className="ms-1">
          <a
            href="https://github.com/kerti/balances/blob/master/LICENSE"
            target="_blank"
            rel="noopener noreferrer"
          >
            License
          </a>
        </span>
      </div>
      <div className="ms-auto small">
        <span className="me-2">
          <a href="https://balances-app.io" target="_blank" rel="noopener noreferrer">
            Balances
          </a>
        </span>
        <span className="me-2">version 0.0.1</span>
        <a href="https://github.com/kerti/balances" target="_blank" rel="noopener noreferrer">
          {<CIcon icon={cibGithub} />}
        </a>
      </div>
    </CFooter>
  )
}

export default React.memo(AppFooter)
