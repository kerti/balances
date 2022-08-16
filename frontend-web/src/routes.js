import { lazy } from 'react';

const Toaster = lazy(() =>
  import('./views/notifications/toaster/Toaster')
)
const Tables = lazy(() => import('./views/base/tables/Tables'))

const Breadcrumbs = lazy(() =>
  import('./views/base/breadcrumbs/Breadcrumbs')
)
const Cards = lazy(() => import('./views/base/cards/Cards'))
const Carousels = lazy(() => import('./views/base/carousels/Carousels'))
const Collapses = lazy(() => import('./views/base/collapses/Collapses'))
const BasicForms = lazy(() => import('./views/base/forms/BasicForms'))

const Jumbotrons = lazy(() =>
  import('./views/base/jumbotrons/Jumbotrons')
)
const ListGroups = lazy(() =>
  import('./views/base/list-groups/ListGroups')
)
const Navbars = lazy(() => import('./views/base/navbars/Navbars'))
const Navs = lazy(() => import('./views/base/navs/Navs'))
const Paginations = lazy(() =>
  import('./views/base/paginations/Paginations')
)
const Popovers = lazy(() => import('./views/base/popovers/Popovers'))
const ProgressBar = lazy(() =>
  import('./views/base/progress-bar/ProgressBar')
)
const Switches = lazy(() => import('./views/base/switches/Switches'))

const Tabs = lazy(() => import('./views/base/tabs/Tabs'))
const Tooltips = lazy(() => import('./views/base/tooltips/Tooltips'))
const BrandButtons = lazy(() =>
  import('./views/buttons/brand-buttons/BrandButtons')
)
const ButtonDropdowns = lazy(() =>
  import('./views/buttons/button-dropdowns/ButtonDropdowns')
)
const ButtonGroups = lazy(() =>
  import('./views/buttons/button-groups/ButtonGroups')
)
const Buttons = lazy(() => import('./views/buttons/buttons/Buttons'))
const Charts = lazy(() => import('./views/charts/Charts'))
const Dashboard = lazy(() => import('./views/dashboard/Dashboard'))
const CoreUIIcons = lazy(() =>
  import('./views/icons/coreui-icons/CoreUIIcons')
)
const Flags = lazy(() => import('./views/icons/flags/Flags'))
const Brands = lazy(() => import('./views/icons/brands/Brands'))
const Alerts = lazy(() => import('./views/notifications/alerts/Alerts'))
const Badges = lazy(() => import('./views/notifications/badges/Badges'))
const Modals = lazy(() => import('./views/notifications/modals/Modals'))
const Colors = lazy(() => import('./views/theme/colors/Colors'))
const Typography = lazy(() =>
  import('./views/theme/typography/Typography')
)
const Widgets = lazy(() => import('./views/widgets/Widgets'))
const Users = lazy(() => import('./views/users/Users'))
const User = lazy(() => import('./views/users/User'))

// Assets Section
const BankAccounts = lazy(() =>
  import('./views/assets/bankAccounts/BankAccounts')
)
const Properties = lazy(() =>
  import('./views/assets/properties/Properties')
)
const Vehicles = lazy(() => import('./views/assets/vehicles/Vehicles'))

// Liabilities Section
const Institutional = lazy(() =>
  import('./views/liabilities/institutional/Institutional')
)
const Personal = lazy(() =>
  import('./views/liabilities/personal/Personal')
)

// Investments Section
const Deposits = lazy(() =>
  import('./views/investments/deposits/Deposits')
)
const Obligations = lazy(() =>
  import('./views/investments/obligations/Obligations')
)
const Gold = lazy(() => import('./views/investments/gold/Gold'))
const MutualFunds = lazy(() =>
  import('./views/investments/mutualFunds/MutualFunds')
)
const Stocks = lazy(() => import('./views/investments/stocks/Stocks'))
const P2pLending = lazy(() =>
  import('./views/investments/p2pLending/P2pLending')
)

const getRoutes = (t) => {
  return [
    { path: '/', exact: true, name: t('home') },
    { path: '/dashboard', name: 'Dashboard', component: Dashboard },
    // Assets Section
    {
      path: '/assets',
      name: t('assets.assets'),
      component: BankAccounts,
      exact: true,
    },
    {
      path: '/assets/bankAccounts',
      name: t('assets.bankAccounts'),
      component: BankAccounts,
    },
    {
      path: '/assets/properties',
      name: t('assets.properties'),
      component: Properties,
    },
    {
      path: '/assets/vehicles',
      name: t('assets.vehicles'),
      component: Vehicles,
    },
    // Liabilities Section
    {
      path: '/liabilities',
      name: t('liabilities.liabilities'),
      component: Institutional,
      exact: true,
    },
    {
      path: '/liabilities/institutional',
      name: t('liabilities.institutional'),
      component: Institutional,
    },
    {
      path: '/liabilities/personal',
      name: t('liabilities.personal'),
      component: Personal,
    },
    // Investments Section
    {
      path: '/investments',
      name: t('investments.investments'),
      component: Deposits,
      exact: true,
    },
    {
      path: '/investments/deposits',
      name: t('investments.deposits'),
      component: Deposits,
    },
    {
      path: '/investments/obligations',
      name: t('investments.obligations'),
      component: Obligations,
    },
    {
      path: '/investments/gold',
      name: t('investments.gold'),
      component: Gold,
    },
    {
      path: '/investments/mutualFunds',
      name: t('investments.mutualFunds'),
      component: MutualFunds,
    },
    {
      path: '/investments/stocks',
      name: t('investments.stocks'),
      component: Stocks,
    },
    {
      path: '/investments/p2pLending',
      name: t('investments.p2pLending'),
      component: P2pLending,
    },
    // To-be-deleted Section
    { path: '/theme', name: 'Theme', component: Colors, exact: true },
    { path: '/theme/colors', name: 'Colors', component: Colors },
    { path: '/theme/typography', name: 'Typography', component: Typography },
    { path: '/base', name: 'Base', component: Cards, exact: true },
    { path: '/base/breadcrumbs', name: 'Breadcrumbs', component: Breadcrumbs },
    { path: '/base/cards', name: 'Cards', component: Cards },
    { path: '/base/carousels', name: 'Carousel', component: Carousels },
    { path: '/base/collapses', name: 'Collapse', component: Collapses },
    { path: '/base/forms', name: 'Forms', component: BasicForms },
    { path: '/base/jumbotrons', name: 'Jumbotrons', component: Jumbotrons },
    { path: '/base/list-groups', name: 'List Groups', component: ListGroups },
    { path: '/base/navbars', name: 'Navbars', component: Navbars },
    { path: '/base/navs', name: 'Navs', component: Navs },
    { path: '/base/paginations', name: 'Paginations', component: Paginations },
    { path: '/base/popovers', name: 'Popovers', component: Popovers },
    {
      path: '/base/progress-bar',
      name: 'Progress Bar',
      component: ProgressBar,
    },
    { path: '/base/switches', name: 'Switches', component: Switches },
    { path: '/base/tables', name: 'Tables', component: Tables },
    { path: '/base/tabs', name: 'Tabs', component: Tabs },
    { path: '/base/tooltips', name: 'Tooltips', component: Tooltips },
    { path: '/buttons', name: 'Buttons', component: Buttons, exact: true },
    { path: '/buttons/buttons', name: 'Buttons', component: Buttons },
    {
      path: '/buttons/button-dropdowns',
      name: 'Dropdowns',
      component: ButtonDropdowns,
    },
    {
      path: '/buttons/button-groups',
      name: 'Button Groups',
      component: ButtonGroups,
    },
    {
      path: '/buttons/brand-buttons',
      name: 'Brand Buttons',
      component: BrandButtons,
    },
    { path: '/charts', name: 'Charts', component: Charts },
    { path: '/icons', exact: true, name: 'Icons', component: CoreUIIcons },
    {
      path: '/icons/coreui-icons',
      name: 'CoreUI Icons',
      component: CoreUIIcons,
    },
    { path: '/icons/flags', name: 'Flags', component: Flags },
    { path: '/icons/brands', name: 'Brands', component: Brands },
    {
      path: '/notifications',
      name: 'Notifications',
      component: Alerts,
      exact: true,
    },
    { path: '/notifications/alerts', name: 'Alerts', component: Alerts },
    { path: '/notifications/badges', name: 'Badges', component: Badges },
    { path: '/notifications/modals', name: 'Modals', component: Modals },
    { path: '/notifications/toaster', name: 'Toaster', component: Toaster },
    { path: '/widgets', name: 'Widgets', component: Widgets },
    { path: '/users', exact: true, name: 'Users', component: Users },
    { path: '/users/:id', exact: true, name: 'User Details', component: User },
  ]
}

export default getRoutes
