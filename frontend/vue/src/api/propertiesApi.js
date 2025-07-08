import axiosInstance from '@/api/index'
import { useEnvUtils } from '@/composables/useEnvUtils'

//// properties CRUD

// create
export async function createPropertyWithAPI(property) {
    try {
        const { data } = await axiosInstance.post('properties', {
            name: property.name,
            address: property.address,
            totalArea: property.totalArea,
            buildingArea: property.buildingArea,
            areaUnit: property.areaUnit,
            type: property.type,
            titleHolder: property.titleHolder,
            taxIdentifier: property.taxIdentifier,
            purchaseDate: property.purchaseDate,
            initialValue: property.initialValue,
            initialValueDate: property.initialValueDate,
            currentValue: property.currentValue,
            currentValueDate: property.currentValueDate,
            annualAppreciationPercent: property.annualAppreciationPercent,
            status: property.status,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function searchPropertiesFromAPI(filter, pageSize) {
    try {
        const { data } = await axiosInstance.post('properties/search', {
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
export async function getPropertyFromAPI(propertyId, startDate, endDate, pageSize, page) {
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
        const { data } = await axiosInstance.get('properties/' + propertyId + '?' + params.toString())
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// update
export async function updatePropertyWithAPI(property) {
    try {
        const { data } = await axiosInstance.patch('properties/' + property.id, property)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// delete
export async function deletePropertyWithAPI(id) {
    try {
        const { data } = await axiosInstance.delete('properties/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

//// property values CRUD

// create
export async function createPropertyValueWithAPI(propertyValue) {
    try {
        const { data } = await axiosInstance.post('properties/values', {
            propertyId: propertyValue.propertyId,
            date: propertyValue.date,
            value: propertyValue.value,
        })
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// read
export async function searchPropertyValuesFromAPI(propertyIds, startDate, endDate, pageSize, page) {
    try {
        const { data } = await axiosInstance.post('properties/values/search', {
            propertyIds: propertyIds,
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
export async function getPropertyValueFromAPI(id) {
    try {
        const { data } = await axiosInstance.get('properties/values/' + id)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// update
export async function updatePropertyValueWithAPI(propertyValue) {
    try {
        const { data } = await axiosInstance.patch('properties/values/' + propertyValue.id, propertyValue)
        return data
    } catch (error) {
        return {
            error: error
        }
    }
}

// delete
export async function deletePropertyValueWithAPI(id) {
    try {
        const { data } = await axiosInstance.delete('properties/values/' + id)
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