import { useDateUtils } from "@/composables/useDateUtils"
import { defineStore } from "pinia"
import { ref } from "vue"

const dateUtils = useDateUtils()

export const useAssetsStore = defineStore('assets', () => {
    // use assets service
    // no toast needed

    ////// reactive state
    const assets = ref([])
    const assetValueHistory = ref([])

    // actions
    async function hydrate() {
        // for now, set data templates to display the data in the view
        const asset1 = {
            name: "Shopping Account",
            class: "cash",
            value: 1200,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        const asset2 = {
            name: "Savings Account",
            class: "cash",
            value: 13000,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        const asset3 = {
            name: "Retirement Account",
            class: "cash",
            value: 100000,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        const asset4 = {
            name: "Mike's Car",
            class: "vehicle",
            value: 36000,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        const asset5 = {
            name: "John's Motorbike",
            class: "vehicle",
            value: 12500,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        const asset6 = {
            name: "Margaret's House",
            class: "property",
            value: 85900,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        const asset7 = {
            name: "Mike's Apartment",
            class: "property",
            value: 99450,
            lastUpdated: dateUtils.getEpochOneYearAgo(),
        }
        assets.value = [asset1, asset2, asset3, asset4, asset5, asset6, asset7]

        const cashHistory = {
            label: "Cash",
            fill: true,
            data: [
                { x: 1722445199000, y: 105843 },
                { x: 1725123599000, y: 110900 },
                { x: 1727715599000, y: 114200 },
            ],
        }

        const vehicleHistory = {
            label: "Vehicles",
            fill: true,
            data: [
                { x: 1722445199000, y: 50700 },
                { x: 1725123599000, y: 49500 },
                { x: 1727715599000, y: 48500 },
            ],
        }

        const propertyHistory = {
            label: "Properties",
            fill: true,
            data: [
                { x: 1722445199000, y: 183746 },
                { x: 1725123599000, y: 184854 },
                { x: 1727715599000, y: 185350 },
            ],
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