import { ref } from "vue"
import { defineStore } from "pinia"

export const useUiStore = defineStore('ui', () => {
    const theme = ref('corporate')
    const lightTheme = ref('corporate')
    const darkTheme = ref('sunset')

    function switchTheme() {
        theme === lightTheme ? darkTheme : lightTheme
    }

    return { theme, lightTheme, darkTheme, switchTheme }
})