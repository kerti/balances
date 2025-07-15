import { searchBankAccountBalancesFromAPI, searchBankAccountsFromAPI } from "@/api/bankAccountsApi"
import { searchPropertiesFromAPI, searchPropertyValuesFromAPI } from "@/api/propertiesApi"
import { searchVehiclesFromAPI, searchVehicleValuesFromAPI } from "@/api/vehiclesApi"
import { useDateUtils } from "@/composables/useDateUtils"

const dateUtils = useDateUtils()

export function useAssetsService() {

    // fetch all data required to display the assets dashboard
    async function getAssetsData() {
        const filter = ''
        const pageSize = 31 * 120
        const startDate = dateUtils.getEpochXYearsAgo(10)
        const endDate = (new Date()).getTime()

        // assets

        const searchAssets = [
            searchBankAccountsFromAPI(filter, pageSize),
            searchVehiclesFromAPI(filter, pageSize),
            searchPropertiesFromAPI(filter, pageSize),
        ]

        const searchAssetsResult = await Promise.all(searchAssets)

        const bankAccounts = searchAssetsResult[0]
        const vehicles = searchAssetsResult[1]
        const properties = searchAssetsResult[2]

        const convertedBankAccounts = bankAccounts.data.items.map((account) => {
            return {
                name: account.accountName,
                class: 'cash',
                value: account.lastBalance,
                lastUpdated: account.lastBalanceDate,
            }
        })

        const convertedVehicles = vehicles.data.items.map((vehicle) => {
            return {
                name: vehicle.name,
                class: 'vehicle',
                value: vehicle.currentValue,
                lastUpdated: vehicle.currentValueDate,
            }
        })

        const convertedProperties = properties.data.items.map((property) => {
            return {
                name: property.name,
                class: 'property',
                value: property.currentValue,
                lastUpdated: property.currentValueDate,
            }
        })

        const convertedAssets = convertedBankAccounts.concat(convertedVehicles, convertedProperties)

        // asset values

        const bankAccountIds = bankAccounts.data.items.map((account) => account.id)
        const vehicleIds = vehicles.data.items.map((vehicle) => vehicle.id)
        const propertyIds = properties.data.items.map((property) => property.id)

        const searchValues = [
            searchBankAccountBalancesFromAPI(bankAccountIds, startDate, endDate, pageSize, 1),
            searchVehicleValuesFromAPI(vehicleIds, startDate, endDate, pageSize, 1),
            searchPropertyValuesFromAPI(propertyIds, startDate, endDate, pageSize, 1),
        ]

        const searchValuesResult = await Promise.all(searchValues)

        const rawBankAccountBalances = searchValuesResult[0]
        const rawVehicleValues = searchValuesResult[1]
        const rawPropertyValues = searchValuesResult[2]

        const bankAccountBalanceDates = rawBankAccountBalances.data.items.map((balance) => balance.date)
        const vehicleValueDates = rawVehicleValues.data.items.map((value) => value.date)
        const propertyValueDates = rawPropertyValues.data.items.map((value) => value.date)

        const bankAccountBalanceUniqueDates = [...new Set(bankAccountBalanceDates)]
        const vehicleValueUniqueDates = [...new Set(vehicleValueDates)]
        const propertyValueUniqueDates = [...new Set(propertyValueDates)]

        const bankAccountBalanceHistory = bankAccountBalanceUniqueDates.map((balanceDate) => {
            const balancesOnDate = rawBankAccountBalances.data.items.filter((bal) => bal.date === balanceDate)
            const sumOfBalancesOnDate = balancesOnDate.reduce((sum, val) => {
                return sum + val.balance
            }, 0)
            return {
                x: balanceDate,
                y: sumOfBalancesOnDate,
            }
        })

        const vehicleValueHistory = vehicleValueUniqueDates.map((valueDate) => {
            const valuesOnDate = rawVehicleValues.data.items.filter((val) => val.date === valueDate)
            const sumOfValuesOnDate = valuesOnDate.reduce((sum, val) => {
                return sum + val.value
            }, 0)
            return {
                x: valueDate,
                y: sumOfValuesOnDate,
            }
        })

        const propertyValueHistory = propertyValueUniqueDates.map((valueDate) => {
            const valuesOnDate = rawPropertyValues.data.items.filter((val) => val.date == valueDate)
            const sumOfValuesOnDate = valuesOnDate.reduce((sum, val) => {
                return sum + val.value
            }, 0)
            return {
                x: valueDate,
                y: sumOfValuesOnDate,
            }
        })

        return {
            assets: convertedAssets,
            cashHistory: bankAccountBalanceHistory.sort((a, b) => a.x - b.x),
            vehicleHistory: vehicleValueHistory.sort((a, b) => a.x - b.x),
            propertyHistory: propertyValueHistory.sort((a, b) => a.x - b.x),
        }
    }

    return {
        getAssetsData,
    }
}