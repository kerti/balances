import Cookies from 'js-cookie'
import { useEnvUtils } from './useEnvUtils'

const ev = useEnvUtils()

export function useAuthCookie() {
    const cookieToken = ev.getCookieToken()
    const cookieUserData = ev.getCookieUserData()

    const setAuthTokenToCookie = (token) => {
        Cookies.set(ev.getCookieToken(), token, {
            expires: 7,
            secure: true,
            sameSite: 'Strict',
            path: '/',
        })
    }

    const getAuthTokenFromCookie = () => Cookies.get(cookieToken)

    const removeAuthTokenFromCookie = () => Cookies.remove(cookieToken)

    const setUserDataToCookie = (userData) => {
        const encodedUserData = btoa(JSON.stringify(userData))
        Cookies.set(cookieUserData, encodedUserData, {
            expires: 7,
            secure: true,
            sameSite: 'Strict',
            path: '/',
        })
    }

    const getUserDataFromCookie = () => {
        const encodedUserData = Cookies.get(cookieUserData)
        return JSON.parse(atob(encodedUserData))
    }

    const removeUserDataFromCookie = () => Cookies.remove(cookieUserData)

    return {
        setAuthTokenToCookie,
        getAuthTokenFromCookie,
        removeAuthTokenFromCookie,
        setUserDataToCookie,
        getUserDataFromCookie,
        removeUserDataFromCookie,
    }
}