import Cookies from 'js-cookie'

export function useAuthCookie() {
    const setAuthTokenToCookie = (token) => {
        Cookies.set(import.meta.env.VITE_COOKIE_TOKEN, token, {
            expires: 7,
            secure: true,
            sameSite: 'Strict',
            path: '/',
        })
    }

    const getAuthTokenFromCookie = () => Cookies.get(import.meta.env.VITE_COOKIE_TOKEN)

    const removeAuthTokenFromCookie = () => Cookies.remove(import.meta.env.VITE_COOKIE_TOKEN)

    const setUserDataToCookie = (userData) => {
        const encodedUserData = btoa(JSON.stringify(userData))
        Cookies.set(import.meta.env.VITE_COOKIE_USERDATA, encodedUserData, {
            expires: 7,
            secure: true,
            sameSite: 'Strict',
            path: '/',
        })
    }

    const getUserDataFromCookie = () => {
        const encodedUserData = Cookies.get(import.meta.env.VITE_COOKIE_USERDATA)
        return JSON.parse(atob(encodedUserData))
    }

    const removeUserDataFromCookie = () => Cookies.remove(import.meta.env.VITE_COOKIE_USERDATA)

    return {
        setAuthTokenToCookie,
        getAuthTokenFromCookie,
        removeAuthTokenFromCookie,
        setUserDataToCookie,
        getUserDataFromCookie,
        removeUserDataFromCookie,
    }
}