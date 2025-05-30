export function useEnvUtils() {

    const getAPIBaseURL = () => {
        const val = import.meta.env.VITE_API_BASE_URL
        if (val === undefined || val === '') {
            console.warn('Unset environment variable VITE_API_BASE_URL, reverting to defaults')
            return 'http://localhost:8080/'
        }
        return val
    }

    const getCookieToken = () => {
        const val = import.meta.env.VITE_COOKIE_TOKEN
        if (val === undefined || val === '') {
            console.warn('Unset environment variable VITE_COOKIE_TOKEN, reverting to defaults')
            return 'token'
        }
        return val
    }

    const getCookieUserData = () => {
        const val = import.meta.env.VITE_COOKIE_USERDATA
        if (val === undefined || val === '') {
            console.warn('Unset environment variable VITE_COOKIE_USERDATA, reverting to defaults')
            return 'userData'
        }
        return val
    }

    const getDefaultLocale = () => {
        const val = import.meta.env.VITE_DEFAULT_LOCALE
        if (val === undefined || val === '') {
            console.warn('Unset environment variable VITE_DEFAULT_LOCALE, reverting to defaults')
            return 'en-US'
        }
        return val
    }

    const getDefaultCurrency = () => {
        const val = import.meta.env.VITE_DEFAULT_CURRENCY
        if (val === undefined || val === '') {
            console.warn('Unset environment variable VITE_DEFAULT_CURRENCY, reverting to defaults')
            return 'USD'
        }
        return val
    }

    const getDefaultPageSize = () => {
        const val = import.meta.env.VITE_DEFAULT_PAGE_SIZE
        if (val === undefined || val === '') {
            console.warn('Unset environment variable VITE_DEFAULT_PAGE_SIZE, reverting to defaults')
        }
        const parsedVal = parseInt(val)
        if (isNaN(parsedVal)) {
            console.warn('Incorrect type for environment variable VITE_DEFAULT_PAGE_SIZE, reverting to defaults')
        }
        return parsedVal
    }

    return {
        getAPIBaseURL,
        getCookieToken,
        getCookieUserData,
        getDefaultLocale,
        getDefaultCurrency,
        getDefaultPageSize,
    }

}