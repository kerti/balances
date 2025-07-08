import {
    // properties
    createPropertyWithAPI,
    searchPropertiesFromAPI,
    getPropertyFromAPI,
    updatePropertyWithAPI,
    deletePropertyWithAPI,
    // property values
    createPropertyValueWithAPI,
    searchPropertyValuesFromAPI,
    getPropertyValueFromAPI,
    updatePropertyValueWithAPI,
    deletePropertyValueWithAPI,
} from '@/api/propertiesApi'

export function usePropertiesService() {

    //// properties CRUD

    // create
    async function createProperty(property) {
        console.log(property)
        const payload = {
            name: property.name,
            address: property.address,
            totalArea: parseFloat(property.totalArea),
            buildingArea: parseFloat(property.buildingArea),
            areaUnit: property.areaUnit,
            type: property.type,
            titleHolder: property.titleHolder,
            taxIdentifier: property.taxIdentifier,
            purchaseDate: property.purchaseDate,
            initialValue: parseFloat(property.initialValue),
            initialValueDate: property.initialValueDate,
            currentValue: parseFloat(property.currentValue),
            currentValueDate: property.currentValueDate,
            annualAppreciationPercent: parseFloat(property.annualAppreciationPercent),
            status: property.status,
        }

        const result = await createPropertyWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // read
    async function searchProperties(filter, valuesStartDate, valuesEndDate, pageSize) {
        // get the properties data
        const properties = await searchPropertiesFromAPI(filter, pageSize)

        if (!properties.error) {
            // then get the values data
            const propertyIds = properties.data.items.map(property => property.id)
            const values = await searchPropertyValuesFromAPI(propertyIds, valuesStartDate, valuesEndDate, 9999, 1)
            if (!values.error) {
                return properties.data.items.map(property => {
                    property.values = values.data.items.filter(function (val) {
                        return val.propertyId == property.id
                    }).sort((a, b) => a.date - b.date)
                    return property
                })
            } else {
                return {
                    error: values.error
                }
            }
        } else {
            return {
                error: properties.error
            }
        }
    }

    // read
    async function getProperty(id, valueStartDate, valueEndDate, pageSize) {
        const property = await getPropertyFromAPI(id, valueStartDate, valueEndDate, pageSize)

        if (!property.error) {
            property.data.values.sort((a, b) => a.date - b.date)
            return property.data
        } else {
            return {
                error: property.error
            }
        }
    }

    // update
    async function updateProperty(property) {
        const payload = {
            id: property.id,
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
        }

        const result = await updatePropertyWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // delete
    async function deleteProperty(id) {
        const result = await deletePropertyWithAPI(id)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    //// property values CRUD

    // create
    async function createPropertyValue(propertyValue) {
        const payload = {
            propertyId: propertyValue.propertyId,
            date: propertyValue.date,
            value: parseInt(propertyValue.value)
        }

        const result = await createPropertyValueWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // read
    async function getPropertyValue(id) {
        const propertyValue = await getPropertyValueFromAPI(id)

        if (!propertyValue.error) {
            return propertyValue.data
        } else {
            return {
                error: propertyValue.error
            }
        }
    }

    // update
    async function updatePropertyValue(propertyValue) {
        const payload = {
            id: propertyValue.id,
            propertyId: propertyValue.propertyId,
            date: propertyValue.date,
            value: parseInt(propertyValue.value),
        }

        const result = await updatePropertyValueWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // delete
    async function deletePropertyValue(id) {
        const result = await deletePropertyValueWithAPI(id)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    return {
        // properties
        createProperty,
        searchProperties,
        getProperty,
        updateProperty,
        deleteProperty,
        // property values
        createPropertyValue,
        getPropertyValue,
        updatePropertyValue,
        deletePropertyValue,
    }
}