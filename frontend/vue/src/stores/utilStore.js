import { ref } from "vue";
import { defineStore } from "pinia";
import { useUtilService } from "@/services/utilService";

export const useUtilStore = defineStore('util', () => {
    const utilService = useUtilService()

    // reactive state
    const serverIsHealthy = ref(false)

    // actions
    async function getServerHealth() {
        serverIsHealthy.value = await utilService.isServerHealthy()
    }

    return {
        serverIsHealthy,
        getServerHealth,
    }
})