import { memo } from 'react';
import { useSelector, useDispatch } from 'react-redux'
import {
  CCreateElement,
  CSidebar,
  CSidebarBrand,
  CSidebarNav,
  CSidebarNavDivider,
  CSidebarNavTitle,
  CSidebarMinimizer,
  CSidebarNavDropdown,
  CSidebarNavItem,
} from '@coreui/react'

import CIcon from '@coreui/icons-react'

import { useTranslation } from 'react-i18next'

// sidebar nav config
import getNavigation from './_nav'

// action
import { setSidebarShow } from '../data/actions/ui'

const TheSidebar = () => {
  const dispatch = useDispatch()
  const show = useSelector((state) => state.ui.sidebarShow)
  const { t } = useTranslation('navigation')

  const translatedNavigation = getNavigation(t)

  return (
    <CSidebar show={show} onShowChange={(val) => dispatch(setSidebarShow(val))}>
      <CSidebarBrand className="d-md-down-none" to="/">
        <CIcon
          className="c-sidebar-brand-full"
          name="logo-negative"
          height={35}
        />
        <CIcon
          className="c-sidebar-brand-minimized"
          name="sygnet"
          height={35}
        />
      </CSidebarBrand>
      <CSidebarNav>
        <CCreateElement
          items={translatedNavigation}
          components={{
            CSidebarNavDivider,
            CSidebarNavDropdown,
            CSidebarNavItem,
            CSidebarNavTitle,
          }}
        />
      </CSidebarNav>
      <CSidebarMinimizer className="c-d-md-down-none" />
    </CSidebar>
  )
}

export default memo(TheSidebar)
