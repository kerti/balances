export function useDateUtils() {
    const epochToLocalDate = (epoch) => {
        const d = new Date(epoch)
        return d.getDate() + ' ' + d.toLocaleDateString('en-US', { month: 'long' })
            + ' ' + d.getFullYear()
    }

    return {
        epochToLocalDate,
    }
}