import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useHistory } from 'react-router-dom'
import {
  CRow,
  CCol,
  CCard,
  CCardHeader,
  CCardBody,
  CDataTable,
  CBadge,
  CPagination,
  CButton,
} from '@coreui/react'
import CardSpinner from '../../common/CardSpinner'
import { useDispatch, useSelector } from 'react-redux'
import { loadBankAccountPage } from '../../../data/actions/assets/bankAccounts'
import CIcon from '@coreui/icons-react'

const getBadge = (item) => {
  switch (item.status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'secondary'
    default:
      return 'primary'
  }
}

const BankAccounts = () => {
  const { t } = useTranslation(['assets', 'formats'])
  const { t: f } = useTranslation('formats')
  const dispatch = useDispatch()
  const [currentPage, setCurrentPage] = useState(1)
  const history = useHistory()
  const fields = [
    { key: 'accountName', label: t('bankAccounts.accountName') },
    { key: 'bankName', label: t('bankAccounts.name') },
    { key: 'accountHolderName', label: t('bankAccounts.accountHolder') },
    { key: 'accountNumber', label: t('bankAccounts.accountNumber') },
    {
      key: 'lastBalance',
      label: t('bankAccounts.lastBalance'),
      _classes: ['text-right'],
    },
    { key: 'status', label: t('bankAccounts.status') },
  ]

  const rawData = useSelector((state) => state.entities.bankAccounts)
  const paginations = useSelector(
    (state) => state.pagination.bankAccountsByKeyword
  )

  const parseData = (rawData, paginations, keyword) => {
    const currentPageInfo = paginations[keyword]
    if (currentPageInfo !== undefined) {
      return {
        items: currentPageInfo.ids.map((id) => {
          return rawData[id]
        }),
        pagination: currentPageInfo,
      }
    } else {
      return {
        items: [],
        pagination: {
          isFetching: true,
        },
      }
    }
  }

  const data = parseData(rawData, paginations, '')

  useEffect(() => {
    if (currentPage > 0) {
      dispatch(loadBankAccountPage('', currentPage))
    }
  }, [dispatch, currentPage])

  const dataTable = (
    <>
      <CDataTable
        items={data.items}
        fields={fields}
        hover
        striped
        bordered
        size="sm"
        itemsPerPage={parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)}
        clickableRows
        onRowClick={(item) => {
          history.push(`/assets/bankAccounts/${item.id}`)
        }}
        scopedSlots={{
          lastBalance: (item) => (
            <td className="text-right">
              {f('number.decimal.2fractions', {
                value: item.lastBalance,
              })}
              <br />
              <small>
                {f('date.long', {
                  value: item.lastBalanceDate,
                })}
              </small>
            </td>
          ),
          status: (item) => (
            <td>
              <CBadge color={getBadge(item)}>
                {t('common.states.' + item.status)}
              </CBadge>
            </td>
          ),
        }}
      />
      <CPagination
        size="sm"
        activePage={currentPage}
        pages={data.pagination ? data.pagination.pageCount : 0}
        onActivePageChange={setCurrentPage}
      />
    </>
  )

  return (
    <CRow>
      <CCol>
        <CCard>
          <CCardHeader>
            {t('bankAccounts.listOfBankAccounts')}
            <div className="card-header-actions">
              <CButton
                size="sm"
                color="primary"
                onClick={() => history.push(`/assets/bankAccounts/new`)}
              >
                <CIcon name="cil-plus" /> {t('common.actions.addNew')}
              </CButton>
            </div>
          </CCardHeader>
          <CCardBody>
            {data.pagination.isFetching ? <CardSpinner /> : dataTable}
          </CCardBody>
        </CCard>
      </CCol>
    </CRow>
  )
}

export default BankAccounts
