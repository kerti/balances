import { expect } from '@playwright/test'

class LoginPage {

    /**
     * 
     * @param {import('@playwright/test').Page} page
     */
    constructor(page) {
        this.page = page
        this.frmLogin = page.getByTestId('frmLogin_form')
        this.txtUsername = page.getByTestId('frmLogin_txtUsername')
        this.txtPassword = page.getByTestId('frmLogin_txtPassword')
        this.btnLogin = page.getByTestId('frmLogin_btnLogin')
        this.divErrorMessage = page.getByTestId('frmLogin_divErrorMessage')
    }

    async login(username, password) {
        await this.page.goto('/login')
        await expect(this.frmLogin).toBeVisible()
        await this.txtUsername.click()
        await this.txtUsername.fill(username)
        await this.txtPassword.click()
        await this.txtPassword.fill(password)
        await this.btnLogin.click()
    }

    async assertTxtUsernameExists() {
        await expect(this.txtUsername).toBeVisible()
    }
}

export { LoginPage }