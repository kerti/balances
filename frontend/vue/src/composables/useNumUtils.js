export function useNumUtils() {

    const numericToMoney = (num) => {
        return Intl.NumberFormat(
            import.meta.env.VITE_DEFAULT_LOCALE,
            {
                style: 'currency',
                currency: import.meta.env.VITE_DEFAULT_CURRENCY
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