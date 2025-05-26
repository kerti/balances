import { getServerHealth } from "@/api/utilApi";

export function useUtilService() {
    async function isServerHealthy() {
        // TODO: transform result to boolean here?
        return await getServerHealth()
    }

    return {
        isServerHealthy,
    }
}