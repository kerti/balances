import { useToast } from "@/composables/useToast"
import { useVehiclesService } from "@/services/vehiclesService"
import { defineStore } from "pinia"
import { ref } from "vue"

export const useVehiclesStore = defineStore('vehicles', () => {
    const svc = useVehiclesService()
    const toast = useToast()

    ////// templates
    const blankVehicle = {
        id: '',
        name: '',
        make: '',
        model: '',
        year: 0,
        type: '',
        titleHolder: '',
        licensePlateNumber: '',
        purchaseDate: (new Date()).getTime(),
        initialValue: 0,
        initialValueDate: (new Date()).getTime(),
        currentValue: 0,
        currentValueDate: (new Date()).getTime(),
        annualDepreciationPercent: 0,
        status: 'in_use',
    }

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
    const lvAddVehicle = ref({})
    // delete vehicle dialog box
    const lvDeleteVehicle = ref({})

    //// detail view
    // main page
    const dvVehicleId = ref('')
    const dvValuesStartDate = ref(0)
    const dvValuesEndDate = ref(0)
    const dvPageSize = ref(10)
    const dvVehicle = ref({})
    const dvVehicleCache = ref({})
    const dvChartData = ref([])
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
        lvAddBankAccount.value = {}
    }

    // CRUD

    async function createVehicle() {
        const res = await svc.createVehicle({
            id: lvAddVehicle.value.id,
            name: lvAddVehicle.value.name,
            make: lvAddVehicle.value.make,
            model: lvAddVehicle.value.model,
            year: lvAddVehicle.value.year,
            type: lvAddVehicle.value.type,
            titleHolder: lvAddVehicle.value.titleHolder,
            licensePlateNumber: lvAddVehicle.value.licensePlateNumber,
            purchaseDate: lvAddVehicle.value.purchaseDate,
            initialValue: lvAddVehicle.value.initialValue,
            initialValueDate: lvAddVehicle.value.initialValueDate,
            currentValue: lvAddVehicle.value.currentValue,
            currentValueDate: lvAddVehicle.value.currentValueDate,
            annualDepreciationPercent: lvAddVehicle.value.annualDepreciationPercent,
            status: lvAddVehicle.value.status,
        })
        if (!res.error) {
            filterVehicles()
            toast.showToast('Vehicle created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create vehicle: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function filterVehicles() {
        lvVehicles.value = await svc.searchVehicles(
            lvFilter.value,
            lvValuesStartDate.value,
            lvValuesEndDate.value,
            lvPageSize.value)
        extractLVChartData()
    }

    // TODO: everything else in list view

    // cache prep and reset

    function resetLVAddVehicleDialog() {
        lvAddVehicle.value = JSON.parse(JSON.stringify(blankVehicle))
    }

    function resetLVDeleteVehicleDialog() {
        lvDeleteVehicle.value = JSON.parse(JSON.stringify(blankVehicle))
    }

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

    // hydration

    async function dvHydrate(initDVVehicleId, initDVValuesStartDate, initDVValuesEndDate, initDVPageSize) {
        dvVehicleId.value = initDVVehicleId
        dvValuesStartDate.value = initDVValuesStartDate
        dvValuesEndDate.value = initDVValuesEndDate
        dvPageSize.value = initDVPageSize
    }

    function dvDehydrate() {
        dvVehicleId.value = ''
        dvValuesStartDate.value = 0
        dvValuesEndDate.value = 0
        dvPageSize.value = 10
        dvVehicle.value = {}
        dvVehicleCache.value = {}
        dvChartData.value = []
        // TODO: set dialog boxes
    }

    // CRUD

    // TODO: create value

    async function getVehicleForDV() {
        const fetchedVehicle = await svc.getVehicle(
            dvVehicleId.value,
            dvValuesStartDate.value,
            dvValuesEndDate.value,
            dvPageSize.value)

        dvVehicle.value = JSON.parse(JSON.stringify(fetchedVehicle))
        dvVehicleCache.value = JSON.parse(JSON.stringify(fetchedVehicle))

        extractDVChartData()
    }

    // TODO: get value by ID

    async function updateVehicle() {
        const res = await svc.updateVehicle(dvVehicle.value)
        if (!res.error) {
            // preserve values records not fetched during vehicle update
            res.values = JSON.parse(JSON.stringify(dvVehicleCache.value.values))
            // then sync the store to the latest data from update
            dvVehicle.value = JSON.parse(JSON.stringify(res))
            dvVehicleCache.value = JSON.parse(JSON.stringify(res))
            toast.showToast('Vehicle updated!', 'success')
        } else {
            toast.showToast('Failed to save vehicle: ' + res.error.message, 'error')
        }
    }

    // TODO: update vehicle value

    // TODO: delete vehicle value

    // cache prep and reset

    function revertDVVehicleToCache() {
        if (dvVehicleCache.value) {
            dvVehicle.value = JSON.parse(JSON.stringify(dvVehicleCache.value))
        }
    }

    // chart utils

    function extractDVChartData() {
        dvChartData.value = [{
            label: dvVehicle.value.name,
            data: dvVehicle.value.values.map(value => {
                return {
                    x: value.date,
                    y: value.value,
                }
            })
        }]
    }

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
        lvAddVehicle,
        // delete vehicle dialog box
        lvDeleteVehicle,

        //// detail view
        // main page
        dvVehicleId,
        dvValuesStartDate,
        dvValuesEndDate,
        dvPageSize,
        dvVehicle,
        dvVehicleCache,
        dvChartData,
        // vehicle editor dialog box


        ////// actions

        //// list view
        // hydration
        lvHydrate,
        lvDehydrate,
        // CRUD
        createVehicle,
        filterVehicles,
        // getVehicleToDeleteById,
        // deleteVehicle,
        // cache and prep
        resetLVAddVehicleDialog,
        resetLVDeleteVehicleDialog,

        //// detail view
        dvHydrate,
        dvDehydrate,
        // CRUD
        getVehicleForDV,
        updateVehicle,
        // cache and prep
        revertDVVehicleToCache,
    }
})