import { ref } from "vue";
import { defineStore } from "pinia";
import { useUtilService } from "@/services/utilService";

export const useUtilStore = defineStore('util', () => {
    const utilService = useUtilService()

    // reactive state
    const serverHealth = ref('Unknown')

    // actions
    async function getServerHealth() {
        serverHealth.value = await utilService.getServerHealth()
    }

    return {
        serverHealth,
        getServerHealth,
    }
})