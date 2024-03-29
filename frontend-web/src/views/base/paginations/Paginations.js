import { useState } from 'react';
import { CCard, CCardBody, CCardHeader, CPagination } from '@coreui/react'

const Paginations = () => {
  const [currentPage, setCurrentPage] = useState(2)

  return (
    <>
      <CCard>
        <CCardHeader>
          Pagination
          <div className="card-header-actions">
            <a
              href="https://coreui.github.io/components/pagination/"
              rel="noreferrer noopener"
              target="_blank"
              className="card-header-action"
            >
              <small className="text-muted">docs</small>
            </a>
          </div>
        </CCardHeader>
        <CCardBody>
          <h6>Default</h6>
          <CPagination
            activePage={currentPage}
            pages={10}
            onActivePageChange={setCurrentPage}
          />
          <br></br>

          <h6>Small</h6>
          <CPagination
            size="sm"
            activePage={currentPage}
            pages={10}
            onActivePageChange={setCurrentPage}
          />
          <br></br>

          <div className="d-md-down-none">
            <h6>Large</h6>
            <CPagination
              size="lg"
              activePage={currentPage}
              pages={10}
              onActivePageChange={setCurrentPage}
            />
            <br></br>
          </div>

          <div>currentPage: {currentPage}</div>
        </CCardBody>
      </CCard>
      <CCard>
        <CCardHeader>
          <strong> Pagination </strong>
          <small>alignment</small>
        </CCardHeader>
        <CCardBody>
          <h6>Left alignment (default)</h6>
          <CPagination
            activePage={currentPage}
            pages={10}
            onActivePageChange={setCurrentPage}
          />
          <br></br>

          <h6>Center alignment</h6>
          <CPagination
            align="center"
            addListClass="some-class"
            activePage={currentPage}
            pages={10}
            onActivePageChange={setCurrentPage}
          />
          <br></br>

          <h6>Right (end) alignment</h6>
          <CPagination
            align="end"
            activePage={currentPage}
            pages={10}
            onActivePageChange={setCurrentPage}
          />
          <br></br>

          <div>currentPage: {currentPage}</div>
        </CCardBody>
      </CCard>
    </>
  )
}

export default Paginations
