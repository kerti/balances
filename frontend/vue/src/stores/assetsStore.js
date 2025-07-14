import { useDateUtils } from "@/composables/useDateUtils"
import { useAssetsService } from "@/services/assetsService"
import { defineStore } from "pinia"
import { ref } from "vue"

const dateUtils = useDateUtils()
const assetsService = useAssetsService()

export const useAssetsStore = defineStore('assets', () => {
    // use assets service
    // no toast needed

    ////// reactive state
    const assets = ref([])
    const assetValueHistory = ref([])

    // actions
    async function hydrate() {
        const data = await assetsService.getAssetsData()
        assets.value = data.assets

        const cashHistory = {
            label: "Cash",
            fill: true,
            data: data.cashHistory,
        }

        const vehicleHistory = {
            label: "Vehicles",
            fill: true,
            data: data.vehicleHistory,
        }

        const propertyHistory = {
            label: "Properties",
            fill: true,
            data: data.propertyHistory,
        }

        assetValueHistory.value = [vehicleHistory, propertyHistory, cashHistory]
    }

    async function dehydrate() {
        assets.value = []
        assetProportions.value = {}
        assetValueHistory.value = []
    }

    return {
        // reactive state
        assets,
        assetValueHistory,
        // actions
        hydrate,
        dehydrate,
    }
})