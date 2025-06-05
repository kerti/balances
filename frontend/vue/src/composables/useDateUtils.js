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

    const epochToShortLocalDate = (epoch) => {
        const d = new Date(epoch)
        return d.getDate() + ' '
            + d.toLocaleDateString(
                import.meta.env.VITE_DEFAULT_LOCALE,
                { month: 'short' })
            + ' ' + d.getFullYear()
    }

    const getEpochOneYearAgo = () => {
        const now = new Date();
        const oneYearAgo = new Date(
            now.getFullYear() - 1,
            now.getMonth(),
            now.getDate(),
            23, 59, 59, 0);
        return oneYearAgo.getTime();
    }

    return {
        epochToLocalDate,
        epochToShortLocalDate,
        getEpochOneYearAgo,
    }
}