import {
    // vehicles
    createVehicleWithAPI,
    searchVehiclesFromAPI,
    getVehicleFromAPI,
    updateVehicleWithAPI,
    deleteVehicleWithAPI,
    // vehicle values
    searchVehicleValuesFromAPI,
} from '@/api/vehiclesApi'

export function useVehiclesService() {

    //// vehicles CRUD

    // create
    async function createVehicle(vehicle) {
        console.log(vehicle)
        const payload = {
            name: vehicle.name,
            make: vehicle.make,
            model: vehicle.model,
            year: parseInt(vehicle.year),
            type: vehicle.type,
            titleHolder: vehicle.titleHolder,
            licensePlateNumber: vehicle.licensePlateNumber,
            purchaseDate: vehicle.purchaseDate,
            initialValue: parseFloat(vehicle.initialValue),
            initialValueDate: vehicle.initialValueDate,
            currentValue: parseFloat(vehicle.currentValue),
            currentValueDate: vehicle.currentValueDate,
            annualDepreciationPercent: parseFloat(vehicle.annualDepreciationPercent),
            status: vehicle.status,
        }

        const result = await createVehicleWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // read
    async function searchVehicles(filter, valuesStartDate, valuesEndDate, pageSize) {
        // get the vehicles data
        const vehicles = await searchVehiclesFromAPI(filter, pageSize)

        if (!vehicles.error) {
            // then get the values data
            const vehicleIds = vehicles.data.items.map(vehicle => vehicle.id)
            const values = await searchVehicleValuesFromAPI(vehicleIds, valuesStartDate, valuesEndDate, 9999, 1)
            if (!values.error) {
                return vehicles.data.items.map(vehicle => {
                    vehicle.values = values.data.items.filter(function (val) {
                        return val.vehicleId == vehicle.id
                    }).sort((a, b) => a.date - b.date)
                    return vehicle
                })
            } else {
                return {
                    error: values.error
                }
            }
        } else {
            return {
                error: vehicles.error
            }
        }
    }

    // read
    async function getVehicle(id, valueStartDate, valueEndDate, pageSize) {
        const vehicle = await getVehicleFromAPI(id, valueStartDate, valueEndDate, pageSize)

        if (!vehicle.error) {
            vehicle.data.values.sort((a, b) => a.date - b.date)
            return vehicle.data
        } else {
            return {
                error: vehicle.error
            }
        }
    }

    // update
    async function updateVehicle(vehicle) {
        const payload = {
            id: vehicle.id,
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
        }

        const result = await updateVehicleWithAPI(payload)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    // delete
    async function deleteVehicle(id) {
        const result = await deleteVehicleWithAPI(id)

        if (!result.error) {
            return result.data
        } else {
            return {
                error: result.error
            }
        }
    }

    return {
        // vehicles
        createVehicle,
        searchVehicles,
        getVehicle,
        updateVehicle,
        deleteVehicle,
        // vehicle values
    }
}