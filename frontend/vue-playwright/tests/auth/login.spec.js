// @ts-check
import { test, expect } from '@playwright/test';
import { LoginPage } from '../../pages/LoginPage'

test.describe("Login Page", () => {
  test('has correct title containing application name', async ({ page }) => {
    await page.goto('/login');

    await expect(page).toHaveTitle(/Balances/);
  });

  test('able to login with correct credentials', async ({ page }) => {
    const loginPage = new LoginPage(page)
    await loginPage.login('admin', 'admin')

    // TODO: check for auth cookies?
    await expect(page).toHaveURL('/')

  })

  test('able to show correct error message when supplied wrong credentials', async ({ page }) => {
    const loginPage = new LoginPage(page)
    await loginPage.login('admin', 'admin123')

    await expect(page.getByTestId('frmLogin_divErrorMessage')).toBeVisible()
    await expect(page.getByTestId('frmLogin_divErrorMessage')).toContainText(/Login failed/)

  })

  test.afterEach(async ({ page }) => {
    await page.close()
  })
})
