export function useNumUtils() {
    const numericToMoney = (num) => {
        return Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(num)
    }

    return {
        numericToMoney,
    }
}