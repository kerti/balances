import axiosInstance from '@/api/index'
import { useEnvUtils } from '@/composables/useEnvUtils'

//// vehicles CRUD

// TODO: create

// read
export async function searchVehiclesFromAPI(filter, pageSize) {
    try {
        const { data } = await axiosInstance.post('vehicles/search', {
            keyword: filter ? filter : '',
            pageSize: getPageSize(pageSize)
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function getVehicleFromAPI(vehicleId, startDate, endDate, pageSize, page) {
    try {
        const params = new URLSearchParams()
        params.append('withValues', true)
        if (startDate) {
            params.append('valuesStartDate', startDate)
        }
        if (endDate) {
            params.append('valuesEndDate', endDate)
        }
        if (pageSize) {
            params.append('pageSize', pageSize)
        }
        if (page) {
            params.append('page', page)
        }
        const { data } = await axiosInstance.get('vehicles/' + vehicleId + '?' + params.toString())
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// TODO: update

// TODO: delete

//// vehicle values CRUD

// TODO: create

// TODO: read
export async function searchVehicleValuesFromAPI(vehicleIds, startDate, endDate, pageSize, page) {
    try {
        const { data } = await axiosInstance.post('vehicles/values/search', {
            vehicleIds: vehicleIds,
            startDate: startDate ? startDate : null,
            endDate: endDate ? endDate : null,
            pageSize: getPageSize(pageSize),
            page: page ? page : 1,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// TODO: read

// TODO: update

// TODO: delete

//// utilities

function getPageSize(pageSize) {
    const ev = useEnvUtils()
    return pageSize ? pageSize : ev.getDefaultPageSize()
}