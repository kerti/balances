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

    return {
        numericToMoney,
    }
}