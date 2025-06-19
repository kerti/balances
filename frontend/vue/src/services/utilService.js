import { getServerHealthFromAPI } from "@/api/utilApi";

export function useUtilService() {
    async function getServerHealth() {
        const result = await getServerHealthFromAPI()
        if (result.error) {
            return result.error.message
        } else if (result.message === 'OK') {
            return 'Server is healthy.'
        } else if (result.message) {
            return 'Server is not healthy - ' + result.message
        } else {
            return 'Unable to fetch server health - unknown error'
        }
    }

    return {
        getServerHealth,
    }
}