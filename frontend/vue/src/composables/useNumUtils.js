import { useEnvUtils } from './useEnvUtils'

const ev = useEnvUtils()

export function useNumUtils() {

    const numericToMoney = (num) => {
        return Intl.NumberFormat(
            ev.getDefaultLocale(),
            {
                style: 'currency',
                currency: ev.getDefaultCurrency()
            }
        ).format(num)
    }

    const queryParamToInt = (queryParam, defaultValue) => {
        if (defaultValue === undefined) {
            console.warn('numUtils.queryParamToInt: no default value supplied, reverting to 0')
            defaultValue = 0
        }

        const parsedQueryParam = parseInt(queryParam)

        if (isNaN(parsedQueryParam)) {
            return defaultValue
        } else {
            return parsedQueryParam
        }
    }

    return {
        numericToMoney,
        queryParamToInt
    }
}