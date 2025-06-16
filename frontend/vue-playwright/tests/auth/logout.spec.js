// @ts-check
import { test, expect } from '@playwright/test';

test.describe("Logout", () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/login');

        await page.getByTestId('frmLogin_txtUsername').click()
        await page.getByTestId('frmLogin_txtUsername').fill('admin')
        await page.getByTestId('frmLogin_txtPassword').click()
        await page.getByTestId('frmLogin_txtPassword').fill('admin')
        await page.getByTestId('frmLogin_btnLogin').click()
    })

    test('able to logout by clicking the logout link in the sidebar', async ({ page }) => {
        await page.getByTestId('divSidebar_btnLogout').click()

        await expect(page).toHaveURL('/login')
    })

    test.afterEach(async ({ page }) => {
        await page.close()
    })
})