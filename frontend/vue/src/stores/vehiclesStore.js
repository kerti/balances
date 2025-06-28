import { useToast } from "@/composables/useToast"
import { useVehiclesService } from "@/services/vehiclesService"
import { defineStore } from "pinia"
import { ref } from "vue"

export const useVehiclesStore = defineStore('vehicles', () => {
    const svc = useVehiclesService()
    const toast = useToast()

    ////// templates

    ////// reactive state

    //// list view
    // main page
    const lvFilter = ref('')
    const lvValuesStartDate = ref(0)
    const lvValuesEndDate = ref(0)
    const lvPageSize = ref(10)
    const lvVehicles = ref([])
    const lvChartData = ref([])
    // add vehicle dialog box
    // delete vehicle dialog box

    //// detail view
    // main page
    // value editor dialog box

    ////// actions

    //// list view
    // hydration
    async function lvHydrate(initLVFilter, initLVValuesStartDate, initLVValuesEndDate, initLVPageSize) {
        lvFilter.value = initLVFilter
        lvValuesStartDate.value = initLVValuesStartDate
        lvValuesEndDate.value = initLVValuesEndDate
        lvPageSize.value = initLVPageSize
    }

    async function lvDehydrate() {
        lvFilter.value = ''
        lvValuesStartDate.value = 0
        lvValuesEndDate.value = 0
        lvPageSize.value = 10
        lvVehicles.value = []
        lvChartData.value = []
        // lvAddBankAccount.value = {}
    }

    // CRUD

    // TODO: add vehicle

    async function filterVehicles() {
        lvVehicles.value = await svc.searchVehicles(
            lvFilter.value,
            lvValuesStartDate.value,
            lvValuesEndDate.value,
            lvPageSize.value)
        extractLVChartData()
    }

    // TODO: everything else in list view

    // chart utils

    function extractLVChartData() {
        lvChartData.value = lvVehicles.value.map(veh => {
            return {
                label: veh.name,
                data: veh.values.map(value => {
                    return {
                        x: value.date,
                        y: value.value,
                    }
                })
            }
        })
    }

    //// detail view

    return {
        ////// reactive state

        //// list view
        // main page
        lvFilter,
        lvValuesStartDate,
        lvValuesEndDate,
        lvPageSize,
        lvVehicles,
        lvChartData,
        // add vehicle dialog box
        // delete vehicle dialog box

        //// detail view
        // main page
        // vehicle editor dialog box

        ////// actions

        //// list view
        // hydration
        lvHydrate,
        lvDehydrate,
        // CRUD
        // createVehicle,
        filterVehicles,
        // getVehicleToDeleteById,
        // deleteVehicle,
        // cache and prep

        //// detail view
        // ...
    }
})