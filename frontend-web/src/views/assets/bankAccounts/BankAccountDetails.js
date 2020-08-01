import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
  CRow,
  CCol,
  CCard,
  CCardHeader,
  CCardBody,
  CForm,
  CFormGroup,
  CLabel,
  CInput,
  CCardFooter,
  CButton,
  CDataTable,
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import { CChartLine } from '@coreui/react-chartjs'
import { useDispatch, useSelector } from 'react-redux'
import { useParams } from 'react-router-dom'
import {
  loadBankAccount,
  loadBankAccountBalancePage,
  updateBankAccount,
} from '../../../data/actions/assets/bankAccounts'
import CardSpinner from '../../common/CardSpinner'
import uniq from 'lodash/uniq'
import { loadUserPage } from '../../../data/actions/system/users'

const getDatasetFromBalanceHistory = (balanceHistoryData) => {
  return balanceHistoryData.map((hist) => ({
    x: new Date(hist.date),
    y: hist.balance,
  }))
}

const getLabelsFromBalanceHistory = (balanceHistoryData) => {
  return balanceHistoryData.map((hist) => new Date(hist.date))
}

const getChartOptions = (t) => {
  return {
    responsive: true,
    maintainAspectRatio: false,
    legend: {
      display: false,
    },
    scales: {
      xAxes: [{ type: 'time' }],
      yAxes: [
        {
          ticks: {
            callback: function (label, index, labels) {
              return label / 1000000 + 'M'
            },
          },
        },
      ],
    },
    tooltips: { enabled: true },
  }
}

const getBalanceHistoryFields = (t) => [
  { key: 'date', label: t('common.date') },
  { key: 'user', label: t('bankAccounts.balanceSetBy') },
  {
    key: 'balance',
    label: t('bankAccounts.balance'),
    _classes: ['text-right'],
  },
]

const getUserIDs = (balances) => {
  const userIDs = balances.map((balanceData) => {
    return balanceData.updatedBy ? balanceData.updatedBy : balanceData.createdBy
  })
  return uniq(userIDs)
}

const Properties = () => {
  const { t } = useTranslation('assets')
  const { t: f } = useTranslation('formats')
  const dispatch = useDispatch()
  const { id } = useParams()

  const rawBankAccounts = useSelector((state) => state.entities.bankAccounts)
  const rawBankAccountBalances = useSelector(
    (state) => state.entities.bankAccountBalances
  )
  const rawUsers = useSelector((state) => state.entities.users)

  const [state, setState] = useState({
    accountName: '',
    bankName: '',
    accountHolderName: '',
    accountNumber: '',
  })

  const chartOptions = getChartOptions(t)

  const balanceHistoryFields = getBalanceHistoryFields(t)

  const account = rawBankAccounts[id]

  const accountReady = account !== undefined && account.balances.length > 0

  const balances = accountReady
    ? account.balances.map((balanceId) => rawBankAccountBalances[balanceId])
    : []

  const balancesReady =
    accountReady && balances
      ? account.balances.length === balances.length &&
        !balances.includes(undefined)
      : false

  const users = balancesReady
    ? balances.map((balance) =>
        balance.updatedBy
          ? rawUsers[balance.updatedBy]
          : rawUsers[balance.createdBy]
      )
    : []

  const usersReady =
    balancesReady && users
      ? balances.length === users.length && !users.includes(undefined)
      : false

  const stateReady =
    usersReady &&
    state.accountName !== '' &&
    state.bankName !== '' &&
    state.accountHolderName !== '' &&
    state.accountNumber !== ''

  const ready = accountReady && balancesReady && usersReady && stateReady

  useEffect(() => {
    if (id) {
      // fetch bank account if necessary
      if (!accountReady) {
        console.log('dispatching loadBankAccount')
        dispatch(loadBankAccount(id, true, 36))
      }

      // fetch balances if necessary
      if (accountReady && !balancesReady) {
        console.log('dispatching loadBankAccountBalancePage')
        dispatch(loadBankAccountBalancePage(id, 1, 36))
      }

      // fetch users if necessary
      if (balancesReady && !usersReady) {
        console.log('dispatching loadUserPage')
        dispatch(loadUserPage(getUserIDs(balances), '', 1, 36))
      }

      // set state if necessary
      if (usersReady && !stateReady) {
        setState({
          accountName: account.accountName,
          bankName: account.bankName,
          accountHolderName: account.accountHolderName,
          accountNumber: account.accountNumber,
        })
      }
    }
  }, [
    dispatch,
    id,
    account,
    balances,
    accountReady,
    balancesReady,
    usersReady,
    stateReady,
  ])

  const handleAccountSubmit = (e) => {
    e.preventDefault()
    console.log('handling submit')
    console.log(e.target.accountName.value)
    if (id) {
      dispatch(
        updateBankAccount(
          id,
          e.target.accountName.value,
          e.target.bankName.value,
          e.target.accountHolderName.value,
          e.target.accountNumber.value,
          'active'
        )
      )
    } else {
      console.log('should be sending new account here')
    }
  }

  const balanceDataTable = (
    <>
      <CDataTable
        items={ready ? balances : []}
        fields={balanceHistoryFields}
        hover
        striped
        bordered
        size="sm"
        itemsPerPage={10}
        pagination
        scopedSlots={{
          date: (item) => <td>{f('date.long', { value: item.date })}</td>,
          user: (item) => (
            <td>
              {usersReady
                ? rawUsers[item.createdBy]
                  ? rawUsers[item.createdBy].name
                  : ''
                : ''}
            </td>
          ),
          balance: (item) => (
            <td className="text-right">
              {f('number.decimal.2fractions', { value: item.balance })}
            </td>
          ),
        }}
      ></CDataTable>
    </>
  )

  const balanceChart = (
    <>
      <CChartLine
        type="line"
        style={{ height: '300px' }}
        datasets={[
          {
            label: t('bankAccounts.balanceHistory'),
            backgroundColor: 'rgb(0,156,195,0.6)',
            data: getDatasetFromBalanceHistory(ready ? balances : []),
          },
        ]}
        options={chartOptions}
        labels={getLabelsFromBalanceHistory(ready ? balances : [])}
      />
    </>
  )

  return (
    <>
      <CRow>
        <CCol xs="12" md="12">
          <CCard>
            <CCardHeader>{t('bankAccounts.balanceHistoryGraph')}</CCardHeader>
            <CCardBody>{ready ? balanceChart : <CardSpinner />}</CCardBody>
          </CCard>
        </CCol>
      </CRow>
      <CRow>
        <CCol xs="12" md="6">
          <CForm
            action=""
            method="post"
            className="form-horizontal"
            onSubmit={handleAccountSubmit}
          >
            <CCard>
              <CCardHeader>{t('bankAccounts.details')}</CCardHeader>
              <CCardBody>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="bank-name">
                      {t('bankAccounts.accountName')}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="accountName"
                      name="accountName"
                      placeholder={t('bankAccounts.accountNamePlaceholder')}
                      defaultValue={state.accountName}
                    />
                  </CCol>
                </CFormGroup>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="bank-name">
                      {t('bankAccounts.name')}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="bankName"
                      name="bankName"
                      placeholder={t('bankAccounts.namePlaceholder')}
                      defaultValue={state.bankName}
                    />
                  </CCol>
                </CFormGroup>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="account-holder">
                      {t('bankAccounts.accountHolder')}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="accountHolderName"
                      name="accountHolderName"
                      placeholder={t('bankAccounts.accountHolderPlaceholder')}
                      defaultValue={state.accountHolderName}
                    />
                  </CCol>
                </CFormGroup>
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="account-number">
                      {t('bankAccounts.accountNumber')}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CInput
                      type="text"
                      id="accountNumber"
                      name="accountNumber"
                      placeholder={t('bankAccounts.accountNumberPlaceholder')}
                      defaultValue={state.accountNumber}
                    />
                  </CCol>
                </CFormGroup>
              </CCardBody>
              <CCardFooter>
                <CButton type="submit" size="sm" color="primary">
                  <CIcon name="cil-scrubber" /> {t('common.actions.submit')}
                </CButton>{' '}
                <CButton type="reset" size="sm" color="danger">
                  <CIcon name="cil-ban" /> {t('common.actions.reset')}
                </CButton>
              </CCardFooter>
            </CCard>
          </CForm>
        </CCol>
        <CCol xs="12" md="6">
          <CCard>
            <CCardHeader>{t('bankAccounts.balanceHistory')}</CCardHeader>
            <CCardBody>{ready ? balanceDataTable : <CardSpinner />}</CCardBody>
          </CCard>
        </CCol>
      </CRow>
    </>
  )
}

export default Properties
