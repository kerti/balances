export function useToast() {
    const showToast = (message, type = 'success', duration = 3000) => {
        let container = document.getElementById('toast-container')

        // if toast container doesn't exist, create and add it to the body
        if (!container) {
            container = document.createElement('div')
            container.id = 'toast-container'
            container.className = 'toast toast-top toast-end fixed z-50 space-y-2 p-4'
            document.body.appendChild(container)
        }

        // create the toast element
        const toast = document.createElement('div')
        toast.className = `alert shadow-lg ${type === 'success'
            ? 'alert-success'
            : type === 'error'
                ? 'alert-error'
                : type === 'warning'
                    ? 'alert-warning'
                    : 'alert-info'
            }`

        toast.innerHTML = `
            <div>
                <span>${message}</span>
            </div>
        `

        // append to container
        container.appendChild(toast)
        if (container.children.length > 5) {
            container.removeChild(container.firstChild)
        }

        // auto-remove after duration
        setTimeout(() => {
            toast.remove()
        }, duration)
    }

    return {
        showToast
    }
}