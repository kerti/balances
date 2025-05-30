export function useDateUtils() {
    // TODO: use some kind of formatting string such as dd MMM yyyy
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