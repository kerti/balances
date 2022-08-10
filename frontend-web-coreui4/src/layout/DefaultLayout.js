import React from 'react'
import { Cookies } from 'react-cookie'
import { AppContent, AppSidebar, AppFooter, AppHeader } from '../components/index'
import { Navigate } from 'react-router-dom'
import cookieNames from '../data/cookies'

const DefaultLayout = () => {
  const cookies = new Cookies()
  const token = cookies.get(cookieNames.auth.token)

  return token ? (
    <div>
      <AppSidebar />
      <div className="wrapper d-flex flex-column min-vh-100 bg-light">
        <AppHeader />
        <div className="body flex-grow-1 px-3">
          <AppContent />
        </div>
        <AppFooter />
      </div>
    </div>
  ) : (
    <Navigate to="/login" />
  )
}

export default DefaultLayout
