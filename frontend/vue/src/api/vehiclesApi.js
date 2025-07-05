import axiosInstance from '@/api/index'
import { useEnvUtils } from '@/composables/useEnvUtils'

//// vehicles CRUD

// create
export async function createVehicleWithAPI(vehicle) {
    try {
        const { data } = await axiosInstance.post('vehicles', {
            name: vehicle.name,
            make: vehicle.make,
            model: vehicle.model,
            year: vehicle.year,
            type: vehicle.type,
            titleHolder: vehicle.titleHolder,
            licensePlateNumber: vehicle.licensePlateNumber,
            purchaseDate: vehicle.purchaseDate,
            initialValue: vehicle.initialValue,
            initialValueDate: vehicle.initialValueDate,
            currentValue: vehicle.currentValue,
            currentValueDate: vehicle.currentValueDate,
            annualDepreciationPercent: vehicle.annualDepreciationPercent,
            status: vehicle.status,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

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

// update
export async function updateVehicleWithAPI(vehicle) {
    try {
        const { data } = await axiosInstance.patch('vehicles/' + vehicle.id, vehicle)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// delete
export async function deleteVehicleWithAPI(id) {
    try {
        const { data } = await axiosInstance.delete('vehicles/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

//// vehicle values CRUD

// create
export async function createVehicleValueWithAPI(vehicleValue) {
    try {
        const { data } = await axiosInstance.post('vehicles/values', {
            vehicleId: vehicleValue.id,
            date: vehicleValue.date,
            value: vehicleValue.value,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
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

// read
export async function getVehicleValueFromAPI(id) {
    try {
        const { data } = await axiosInstance.get('vehicles/values/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// update
export async function updateVehicleValueWithAPI(vehicleValue) {
    try {
        const { data } = await axiosInstance.patch('vehicles/values/' + vehicleValue.id, vehicleValue)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// delete
export async function deleteVehicleValueWithAPI(id) {
    try {
        const { data } = await axiosInstance.delete('vehicles/values/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

//// utilities

function getPageSize(pageSize) {
    const ev = useEnvUtils()
    return pageSize ? pageSize : ev.getDefaultPageSize()
}