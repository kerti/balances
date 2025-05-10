import { useEffect, useMemo, useState } from 'react'
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
  CModal,
  CModalHeader,
  CModalTitle,
  CModalBody,
  CAlert,
  CSelect,
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import { CChartLine } from '@coreui/react-chartjs'
import { useDispatch, useSelector } from 'react-redux'
import { useParams } from 'react-router-dom'
import {
  loadBankAccount,
  createBankAccount,
  updateBankAccount,
  createBankAccountBalance,
  showBankAccountBalanceModal,
  hideBankAccountBalanceModal,
  updateBankAccountBalance,
  deleteBankAccountBalance,
} from '../../../data/actions/assets/bankAccounts'
import CardSpinner from '../../common/CardSpinner'
import uniq from 'lodash/uniq'
import { loadUserPage } from '../../../data/actions/system/users'
import { resetErrorMessage } from '../../../data/actions/system/errors'

const getDatasetFromBalanceHistory = (balanceHistoryData) => {
  return balanceHistoryData.map((hist) => ({
    x: new Date(hist.date),
    y: hist.balance,
  }))
}

const getLabelsFromBalanceHistory = (balanceHistoryData) => {
  return balanceHistoryData.map((hist) => new Date(hist.date))
}

// TODO: dynamic scaling?
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
  {
    key: 'balanceActions',
    label: '',
  },
]

const getUserIDs = (balances) => {
  const userIDs = balances.map((balanceData) => {
    return balanceData.updatedBy ? balanceData.updatedBy : balanceData.createdBy
  })
  return uniq(userIDs)
}

const BankAccounts = () => {
  const { t } = useTranslation('assets')
  const { t: f } = useTranslation('formats')
  const dispatch = useDispatch()
  const { id } = useParams()

  const [modalState, setModalState] = useState({
    id: '',
    date: '',
    balance: '',
    isDelete: false,
  })

  const rawBankAccounts = useSelector((state) => state.entities.bankAccounts)
  const rawBankAccountBalances = useSelector(
    (state) => state.entities.bankAccountBalances
  )
  const rawUsers = useSelector((state) => state.entities.users)

  const [state, setState] = useState({
    accountName: '',
    accountBankName: '',
    accountHolderName: '',
    accountNumber: '',
    balances: '',
    accountStatus: '',
  })

  const errorMessage = useSelector((state) => state.errorMessage)
  const updateBankAccountState = useSelector(
    (state) => state.api.updateBankAccount
  )
  const createBankAccountBalanceState = useSelector(
    (state) => state.api.createBankAccountBalance
  )
  const balanceModalState = useSelector(
    (state) => state.ui.modals.bankAccountsBalance
  )

  const chartOptions = getChartOptions(t)
  const balanceHistoryFields = getBalanceHistoryFields(t)

  const account = rawBankAccounts[id]

  const accountReady = account !== undefined && account.balances.length > 0

  const balances = useMemo(
    () =>
      accountReady
        ? account.balances.map((balanceId) => rawBankAccountBalances[balanceId])
        : [],
    [accountReady, account, rawBankAccountBalances]
  )

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
    state.accountBankName !== '' &&
    state.accountHolderName !== '' &&
    state.accountNumber !== '' &&
    state.accountStatus !== ''

  const ready = accountReady && balancesReady && usersReady && stateReady

  useEffect(() => {
    if (id) {
      if (!ready) {
        dispatch([
          loadBankAccount(id, true, 36),
          loadUserPage(getUserIDs(balances), '', 1, 36),
        ])
      }

      // set state if necessary
      if (usersReady && !stateReady) {
        setState({
          accountName: account.accountName,
          accountBankName: account.bankName,
          accountHolderName: account.accountHolderName,
          accountNumber: account.accountNumber,
          accountStatus: account.status,
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
    ready,
  ])

  const handleAccountSubmit = (e) => {
    e.preventDefault()
    if (id) {
      dispatch([
        updateBankAccount(
          id,
          e.target.accountName.value,
          e.target.accountBankName.value,
          e.target.accountHolderName.value,
          e.target.accountNumber.value,
          e.target.accountStatus.value
        ),
        loadBankAccount(id, true, 36, [], true),
      ])
    } else {
      dispatch([
        createBankAccount(
          e.target.accountName.value,
          e.target.accountBankName.value,
          e.target.accountHolderName.value,
          e.target.accountNumber.value,
          e.target.accountStatus.value
        ),
        loadBankAccount(id, true, 36, [], true),
      ])
    }
  }

  const handleBalanceModalSubmit = (e) => {
    e.preventDefault()
    if (modalState.isDelete) {
      dispatch([
        deleteBankAccountBalance(modalState.id),
        loadBankAccount(id, true, 36, [], true),
        resetErrorMessage(),
        hideBankAccountBalanceModal(),
      ])
    } else {
      if (e.target.balanceModalId.value === '') {
        dispatch([
          createBankAccountBalance(
            id,
            Date.parse(modalState.date),
            parseInt(modalState.balance),
            { nextAction: hideBankAccountBalanceModal }
          ),
          loadBankAccount(id, true, 36, [], true),
          resetErrorMessage(),
          hideBankAccountBalanceModal(),
        ])
      } else {
        dispatch([
          updateBankAccountBalance(
            modalState.id,
            id,
            Date.parse(modalState.date),
            parseInt(modalState.balance),
            { nextAction: hideBankAccountBalanceModal }
          ),
          loadBankAccount(id, true, 36, [], true),
          resetErrorMessage(),
          hideBankAccountBalanceModal(),
        ])
      }
    }
  }

  const handleInitEdit = (item) => {
    // TODO: move this to i18next interpolator somehow
    const yourDate = new Date(item.date)
    const yodaDate = yourDate.toISOString().split('T')[0]
    setModalState({
      id: item.id,
      date: yodaDate,
      balance: item.balance,
      isDelete: false,
    })
    dispatch([resetErrorMessage(), showBankAccountBalanceModal()])
  }

  const handleInitDelete = (item) => {
    // TODO: move this to i18next interpolator somehow
    const yourDate = new Date(item.date)
    const yodaDate = yourDate.toISOString().split('T')[0]
    setModalState({
      id: item.id,
      date: yodaDate,
      balance: item.balance,
      isDelete: true,
    })
    dispatch([resetErrorMessage(), showBankAccountBalanceModal()])
  }

  const handleBalanceModalClose = (e) => {
    dispatch(hideBankAccountBalanceModal())
    dispatch(resetErrorMessage())
    setModalState({
      id: '',
      date: '',
      balance: '',
      isDelete: false,
    })
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
        clickableRows
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
              {f('number.decimal.2fractions', { value: item.balance })}{' '}
            </td>
          ),
          balanceActions: (item) => (
            <td className="text-right">
              {
                <CButton
                  className="btn-ghost-primary"
                  size="sm"
                  onClick={() => {
                    handleInitEdit(item)
                  }}
                >
                  <CIcon name="cil-pencil" />
                </CButton>
              }{' '}
              {
                <CButton
                  className="btn-ghost-danger"
                  size="sm"
                  onClick={() => {
                    handleInitDelete(item)
                  }}
                >
                  <CIcon name="cil-trash" />
                </CButton>
              }
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
                      id="accountBankName"
                      name="accountBankName"
                      placeholder={t('bankAccounts.namePlaceholder')}
                      defaultValue={state.accountBankName}
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
                <CFormGroup row>
                  <CCol md="3">
                    <CLabel htmlFor="account-status">
                      {t('bankAccounts.status')}
                    </CLabel>
                  </CCol>
                  <CCol xs="12" md="9">
                    <CSelect
                      custom
                      id="accountStatus"
                      name="accountStatus"
                      value={state.accountStatus}
                      onChange={(i) => {
                        setState({ ...state, accountStatus: i.value })
                      }}
                    >
                      <option value="">
                        {t('common.actions.selectAnOption')}
                      </option>
                      <option value="active">
                        {t('common.states.active')}
                      </option>
                      <option value="inactive">
                        {t('common.states.inactive')}
                      </option>
                    </CSelect>
                  </CCol>
                </CFormGroup>
              </CCardBody>
              <CCardFooter>
                <CButton
                  type="submit"
                  size="sm"
                  color="primary"
                  disabled={updateBankAccountState.isFetching}
                >
                  <CIcon name="cil-scrubber" /> {t('common.actions.submit')}
                </CButton>{' '}
                <CButton
                  type="reset"
                  size="sm"
                  color="danger"
                  disabled={updateBankAccountState.isFetching}
                >
                  <CIcon name="cil-ban" /> {t('common.actions.reset')}
                </CButton>
              </CCardFooter>
            </CCard>
          </CForm>
        </CCol>
        <CCol xs="12" md="6">
          <CCard>
            <CCardHeader>
              {t('bankAccounts.balanceHistory')}
              <div className="card-header-actions">
                <CButton
                  size="sm"
                  color="primary"
                  onClick={() => {
                    dispatch([
                      resetErrorMessage(),
                      showBankAccountBalanceModal(),
                    ])
                  }}
                >
                  <CIcon name="cil-plus" /> {t('common.actions.addNew')}
                </CButton>
              </div>
            </CCardHeader>
            <CCardBody>{ready ? balanceDataTable : <CardSpinner />}</CCardBody>
          </CCard>
        </CCol>
      </CRow>
      <CModal
        show={balanceModalState.show}
        onClose={handleBalanceModalClose}
        size="lg"
      >
        <CModalHeader closeButton>
          <CModalTitle>
            {modalState.isDelete
              ? t('bankAccounts.modalTitle.deleteBalance')
              : modalState.id === ''
              ? t('bankAccounts.modalTitle.addBalance')
              : t('bankAccounts.modalTitle.editBalance')}
          </CModalTitle>
        </CModalHeader>
        <CModalBody>
          <CAlert
            color="danger"
            closeButton
            onClick={() => {
              dispatch(resetErrorMessage())
            }}
            show={errorMessage !== null}
          >
            {errorMessage}
          </CAlert>
          <CForm
            action=""
            method="post"
            className="form-horizontal"
            onSubmit={handleBalanceModalSubmit}
          >
            <CFormGroup row>
              <CCol md="3">
                <CLabel htmlFor="date">{t('common.date')}</CLabel>
              </CCol>
              <CCol xs="12" md="9">
                <CInput
                  type="hidden"
                  id="balanceModalId"
                  name="balanceModalId"
                  value={modalState.id}
                  onChange={(i) => {
                    const ms = { ...modalState, id: i.target.value }
                    setModalState(ms)
                  }}
                />
                <CInput
                  type="date"
                  id="balanceModalDate"
                  name="balanceModalDate"
                  value={modalState.date}
                  readOnly={modalState.isDelete}
                  onChange={(i) => {
                    const ms = { ...modalState, date: i.target.value }
                    setModalState(ms)
                  }}
                />
              </CCol>
            </CFormGroup>
            <CFormGroup row>
              <CCol md="3">
                <CLabel htmlFor="date">{t('bankAccounts.balance')}</CLabel>
              </CCol>
              <CCol xs="12" md="9">
                <CInput
                  type="text"
                  id="balanceModalBalanceAtDate"
                  name="balanceModalBalanceAtDate"
                  placeholder={t('bankAccounts.balancePlaceholder')}
                  value={modalState.balance}
                  readOnly={modalState.isDelete}
                  onChange={(i) => {
                    const ms = { ...modalState, balance: i.target.value }
                    setModalState(ms)
                  }}
                />
              </CCol>
            </CFormGroup>
            <CFormGroup row>
              <CCol className="text-right">
                <CButton
                  color={modalState.isDelete ? 'danger' : 'primary'}
                  type="submit"
                  disabled={createBankAccountBalanceState.isFetching}
                >
                  {modalState.isDelete
                    ? t('common.actions.confirm')
                    : t('common.actions.submit')}
                </CButton>{' '}
                <CButton
                  color="secondary"
                  onClick={handleBalanceModalClose}
                  disabled={createBankAccountBalanceState.isFetching}
                >
                  {t('common.actions.cancel')}
                </CButton>
              </CCol>
            </CFormGroup>
          </CForm>
        </CModalBody>
      </CModal>
    </>
  )
}

export default BankAccounts
