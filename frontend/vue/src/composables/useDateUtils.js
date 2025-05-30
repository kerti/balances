export function useDateUtils() {
    const epochToLocalDate = (epoch) => {
        const d = new Date(epoch)
        return d.getDate() + ' '
            + d.toLocaleDateString(
                import.meta.env.VITE_DEFAULT_LOCALE,
                { month: 'long' })
            + ' ' + d.getFullYear()
    }

    return {
        epochToLocalDate,
    }
}