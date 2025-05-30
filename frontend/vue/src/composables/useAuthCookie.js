import Cookies from 'js-cookie'
import { useEnvUtils } from './useEnvUtils'

const ev = useEnvUtils()

export function useAuthCookie() {
    const setAuthTokenToCookie = (token) => {
        Cookies.set(ev.getCookieToken(), token, {
            expires: 7,
            secure: true,
            sameSite: 'Strict',
            path: '/',
        })
    }

    const getAuthTokenFromCookie = () => Cookies.get(ev.getCookieToken())

    const removeAuthTokenFromCookie = () => Cookies.remove(ev.getCookieToken())

    const setUserDataToCookie = (userData) => {
        const encodedUserData = btoa(JSON.stringify(userData))
        Cookies.set(ev.getCookieUserData(), encodedUserData, {
            expires: 7,
            secure: true,
            sameSite: 'Strict',
            path: '/',
        })
    }

    const getUserDataFromCookie = () => {
        const encodedUserData = Cookies.get(ev.getCookieUserData())
        return JSON.parse(atob(encodedUserData))
    }

    const removeUserDataFromCookie = () => Cookies.remove(ev.getCookieUserData())

    return {
        setAuthTokenToCookie,
        getAuthTokenFromCookie,
        removeAuthTokenFromCookie,
        setUserDataToCookie,
        getUserDataFromCookie,
        removeUserDataFromCookie,
    }
}