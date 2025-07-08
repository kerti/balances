import { useToast } from "@/composables/useToast"
import { usePropertiesService } from "@/services/propertiesService"
import { defineStore } from "pinia"
import { ref } from "vue"

export const usePropertiesStore = defineStore('properties', () => {
    const svc = usePropertiesService()
    const toast = useToast()

    ////// templates
    const blankProperty = {
        id: '',
        name: '',
        address: '',
        totalArea: 0,
        buildingArea: 0,
        areaUnit: 'sqm',
        type: '',
        titleHolder: '',
        taxIdentifier: '',
        purchaseDate: (new Date()).getTime(),
        initialValue: 0,
        initialValueDate: (new Date()).getTime(),
        currentValue: 0,
        currentValueDate: (new Date()).getTime(),
        annualAppreciationPercent: 0,
        status: 'in_use',
    }
    const blankPropertyValue = {
        id: '',
        propertyId: '',
        date: (new Date()).getTime(),
        value: 0,
    }

    ////// reactive state

    //// list view
    // main page
    const lvFilter = ref('')
    const lvValuesStartDate = ref(0)
    const lvValuesEndDate = ref(0)
    const lvPageSize = ref(10)
    const lvProperties = ref([])
    const lvChartData = ref([])
    // add property dialog box
    const lvAddProperty = ref({})
    // delete property dialog box
    const lvDeleteProperty = ref({})

    //// detail view
    // main page
    const dvPropertyId = ref('')
    const dvValuesStartDate = ref(0)
    const dvValuesEndDate = ref(0)
    const dvPageSize = ref(10)
    const dvProperty = ref({})
    const dvPropertyCache = ref({})
    const dvChartData = ref([])
    // value editor dialog box
    const dvValueEditorMode = ref('Add')
    const dvEditPropertyValue = ref({})
    const dvEditPropertyValueCache = ref({})

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
        lvProperties.value = []
        lvChartData.value = []
        lvAddProperty.value = {}
    }

    // CRUD

    async function createProperty() {
        const res = await svc.createProperty({
            id: lvAddProperty.value.id,
            name: lvAddProperty.value.name,
            address: lvAddProperty.value.address,
            totalArea: lvAddProperty.value.totalArea,
            buildingArea: lvAddProperty.value.buildingArea,
            areaUnit: lvAddProperty.value.areaUnit,
            type: lvAddProperty.value.type,
            titleHolder: lvAddProperty.value.titleHolder,
            taxIdentifier: lvAddProperty.value.taxIdentifier,
            purchaseDate: lvAddProperty.value.purchaseDate,
            initialValue: lvAddProperty.value.initialValue,
            initialValueDate: lvAddProperty.value.initialValueDate,
            currentValue: lvAddProperty.value.currentValue,
            currentValueDate: lvAddProperty.value.currentValueDate,
            annualAppreciationPercent: lvAddProperty.value.annualAppreciationPercent,
            status: lvAddProperty.value.status,
        })
        if (!res.error) {
            filterProperties()
            toast.showToast('Property created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create property: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function filterProperties() {
        lvProperties.value = await svc.searchProperties(
            lvFilter.value,
            lvValuesStartDate.value,
            lvValuesEndDate.value,
            lvPageSize.value)
        extractLVChartData()
    }

    async function getPropertyToDeleteById(id) {
        const res = await svc.getProperty(id, null, null, 0)
        if (!res.error) {
            lvDeleteProperty.value = res
            return res
        } else {
            toast.showToast('Failed to retrieve property for deletion: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function deleteProperty() {
        const res = await svc.deleteProperty(lvDeleteProperty.value.id)
        if (!res.error) {
            filterProperties()
            toast.showToast('Property deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete property: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    // cache prep and reset

    function resetLVAddPropertyDialog() {
        lvAddProperty.value = JSON.parse(JSON.stringify(blankProperty))
    }

    function resetLVDeletePropertyDialog() {
        lvDeleteProperty.value = JSON.parse(JSON.stringify(blankProperty))
    }

    // chart utils

    function extractLVChartData() {
        lvChartData.value = lvProperties.value.map(veh => {
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

    async function dvHydrate(initDVPropertyId, initDVValuesStartDate, initDVValuesEndDate, initDVPageSize) {
        dvPropertyId.value = initDVPropertyId
        dvValuesStartDate.value = initDVValuesStartDate
        dvValuesEndDate.value = initDVValuesEndDate
        dvPageSize.value = initDVPageSize
    }

    function dvDehydrate() {
        dvPropertyId.value = ''
        dvValuesStartDate.value = 0
        dvValuesEndDate.value = 0
        dvPageSize.value = 10
        dvProperty.value = {}
        dvPropertyCache.value = {}
        dvChartData.value = []
        dvValueEditorMode.value = 'Add'
        dvEditPropertyValue.value = {}
        dvEditPropertyValueCache.value = {}
    }

    // CRUD

    async function createPropertyValue() {
        const res = await svc.createPropertyValue({
            propertyId: dvEditPropertyValue.value.propertyId,
            date: dvEditPropertyValue.value.date,
            value: dvEditPropertyValue.value.value,
        })
        if (!res.error) {
            getPropertyForDV()
            toast.showToast('Property value created!', 'success')
            return res
        } else {
            toast.showToast('Failed to create property value: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function getPropertyForDV() {
        const fetchedProperty = await svc.getProperty(
            dvPropertyId.value,
            dvValuesStartDate.value,
            dvValuesEndDate.value,
            dvPageSize.value)

        dvProperty.value = JSON.parse(JSON.stringify(fetchedProperty))
        dvPropertyCache.value = JSON.parse(JSON.stringify(fetchedProperty))

        extractDVChartData()
    }

    async function getPropertyValueById(id) {
        const fetchedValue = await svc.getPropertyValue(id)
        dvEditPropertyValue.value = JSON.parse(JSON.stringify(fetchedValue))
        dvEditPropertyValueCache.value = JSON.parse(JSON.stringify(fetchedValue))
    }

    async function updateProperty() {
        const res = await svc.updateProperty(dvProperty.value)
        if (!res.error) {
            // preserve values records not fetched during property update
            res.values = JSON.parse(JSON.stringify(dvPropertyCache.value.values))
            // then sync the store to the latest data from update
            dvProperty.value = JSON.parse(JSON.stringify(res))
            dvPropertyCache.value = JSON.parse(JSON.stringify(res))
            toast.showToast('Property updated!', 'success')
        } else {
            toast.showToast('Failed to save property: ' + res.error.message, 'error')
        }
    }

    async function updatePropertyValue() {
        const res = await svc.updatePropertyValue(dvEditPropertyValue.value)
        if (!res.error) {
            getPropertyForDV()
            getPropertyValueById(res.id)
            toast.showToast('Property value updated!', 'success')
            return res
        } else {
            toast.showToast('Failed to save property value: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    async function deletePropertyValue() {
        const res = await svc.deletePropertyValue(dvEditPropertyValue.value.id)
        if (!res.error) {
            getPropertyForDV()
            toast.showToast('Property value deleted!', 'success')
            return res
        } else {
            toast.showToast('Failed to delete property value: ' + res.error.message, 'error')
            return {
                error: res.error
            }
        }
    }

    // cache prep and reset

    function revertDVPropertyToCache() {
        if (dvPropertyCache.value) {
            dvProperty.value = JSON.parse(JSON.stringify(dvPropertyCache.value))
        }
    }

    function revertDVPropertyValueToCache() {
        if (dvEditPropertyValueCache.value) {
            dvEditPropertyValue.value = JSON.parse(JSON.stringify(dvEditPropertyValueCache.value))
        }
    }

    function prepDVBlankPropertyValue() {
        const template = JSON.parse(JSON.stringify(blankPropertyValue))
        template.propertyId = dvPropertyId.value
        dvEditPropertyValue.value = JSON.parse(JSON.stringify(template))
        dvEditPropertyValueCache.value = JSON.parse(JSON.stringify(template))
    }

    // chart utils

    function extractDVChartData() {
        dvChartData.value = [{
            label: dvProperty.value.name,
            data: dvProperty.value.values.map(value => {
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
        lvProperties,
        lvChartData,
        // add property dialog box
        lvAddProperty,
        // delete property dialog box
        lvDeleteProperty,

        //// detail view
        // main page
        dvPropertyId,
        dvValuesStartDate,
        dvValuesEndDate,
        dvPageSize,
        dvProperty,
        dvPropertyCache,
        dvChartData,
        // property editor dialog box
        dvValueEditorMode,
        dvEditPropertyValue,
        dvEditPropertyValueCache,

        ////// actions

        //// list view
        // hydration
        lvHydrate,
        lvDehydrate,
        // CRUD
        createProperty,
        filterProperties,
        getPropertyToDeleteById,
        deleteProperty,
        // cache and prep
        resetLVAddPropertyDialog,
        resetLVDeletePropertyDialog,

        //// detail view
        dvHydrate,
        dvDehydrate,
        // CRUD
        createPropertyValue,
        getPropertyForDV,
        getPropertyValueById,
        updateProperty,
        updatePropertyValue,
        deletePropertyValue,
        // cache and prep
        revertDVPropertyToCache,
        revertDVPropertyValueToCache,
        prepDVBlankPropertyValue,
    }
})