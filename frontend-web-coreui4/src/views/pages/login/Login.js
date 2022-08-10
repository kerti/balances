import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  CButton,
  CCard,
  CCardBody,
  CCardGroup,
  CCol,
  CContainer,
  CForm,
  CFormInput,
  CInputGroup,
  CInputGroupText,
  CRow,
  CDropdown,
  CDropdownToggle,
  CDropdownMenu,
  CDropdownItem,
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import { cilLockLocked, cilUser } from '@coreui/icons'

import { useTranslation } from 'react-i18next'
import { useDispatch, useSelector } from 'react-redux'
import flagIconMap from '../../../translations/flags.js'

// actions
import { requestLogin } from '../../../data/actions/auth'
import { setLang } from '../../../data/actions/ui'

const Login = () => {
  const { t } = useTranslation('app')
  const dispatch = useDispatch()
  const authLoading = useSelector((state) => state.auth.loading)
  const lang = useSelector((state) => state.ui.lang)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()

  const handleSubmit = (e) => {
    e.preventDefault()
    dispatch(requestLogin(username, password, navigate))
  }

  const selectLang = (lang) => {
    dispatch(setLang(lang))
  }

  return (
    <div className="bg-light min-vh-100 d-flex flex-row align-items-center">
      <CContainer>
        <CRow className="justify-content-center">
          <CCol md={8}>
            <CCardGroup>
              <CCard className="p-4">
                <CCardBody>
                  <CForm onSubmit={handleSubmit}>
                    <CRow>
                      <CCol xs={6}>
                        <h1>{t('login.login')}</h1>
                        <p className="text-muted">{t('login.signInToYourAccount')}</p>
                      </CCol>
                      <CCol xs={6} className="text-end">
                        <CDropdown className="m-1">
                          <CDropdownToggle>
                            <CIcon icon={flagIconMap[lang]} size="lg" />
                          </CDropdownToggle>
                          <CDropdownMenu>
                            <CDropdownItem onClick={() => selectLang('en')}>English</CDropdownItem>
                            <CDropdownItem onClick={() => selectLang('id')}>
                              Bahasa Indonesia
                            </CDropdownItem>
                          </CDropdownMenu>
                        </CDropdown>
                      </CCol>
                    </CRow>
                    <CInputGroup className="mb-3">
                      <CInputGroupText>
                        <CIcon icon={cilUser} />
                      </CInputGroupText>
                      <CFormInput
                        type="text"
                        placeholder={t('login.username')}
                        autoComplete="username"
                        value={username}
                        onChange={(u) => setUsername(u.target.value)}
                      />
                    </CInputGroup>
                    <CInputGroup className="mb-4">
                      <CInputGroupText>
                        <CIcon icon={cilLockLocked} />
                      </CInputGroupText>
                      <CFormInput
                        type="password"
                        placeholder={t('login.password')}
                        autoComplete="current-password"
                        value={password}
                        onChange={(p) => setPassword(p.target.value)}
                      />
                    </CInputGroup>
                    <CRow>
                      <CCol xs={6}>
                        <CButton
                          type="submit"
                          color="primary"
                          className="px-4"
                          disabled={authLoading}
                        >
                          {t('login.login')}
                        </CButton>
                      </CCol>
                      <CCol xs={6} className="text-end">
                        <CButton color="link" className="px-0" hidden>
                          Forgot password?
                        </CButton>
                      </CCol>
                    </CRow>
                  </CForm>
                </CCardBody>
              </CCard>
              <CCard className="text-white bg-primary py-5" style={{ width: '44%' }} hidden>
                <CCardBody className="text-center">
                  <div>
                    <h2>Sign up</h2>
                    <p>
                      Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
                      tempor incididunt ut labore et dolore magna aliqua.
                    </p>
                    <Link to="/register">
                      <CButton color="primary" className="mt-3" active tabIndex={-1}>
                        Register Now!
                      </CButton>
                    </Link>
                  </div>
                </CCardBody>
              </CCard>
            </CCardGroup>
          </CCol>
        </CRow>
      </CContainer>
    </div>
  )
}

export default Login
