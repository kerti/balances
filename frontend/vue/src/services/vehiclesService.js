import { getVehicleFromAPI, searchVehiclesFromAPI, searchVehicleValuesFromAPI } from '@/api/vehiclesApi'

export function useVehiclesService() {

    //// vehicles CRUD

    // TODO: create

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

    return {
        // vehicles
        searchVehicles,
        getVehicle,
        // vehicle values
    }
}