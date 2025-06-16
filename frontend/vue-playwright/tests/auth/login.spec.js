// @ts-check
import { test, expect } from '@playwright/test';

test.describe("Login Page", () => {
  test('has correct title containing application name', async ({ page }) => {
    await page.goto('/login');

    await expect(page).toHaveTitle(/Balances/);
  });

  test('able to login with correct credentials', async ({ page }) => {
    await page.goto('/login');

    await page.getByTestId('frmLogin_txtUsername').click()
    await page.getByTestId('frmLogin_txtUsername').fill('admin')
    await page.getByTestId('frmLogin_txtPassword').click()
    await page.getByTestId('frmLogin_txtPassword').fill('admin')
    await page.getByTestId('frmLogin_btnLogin').click()

    // TODO: check for auth cookies?
    await expect(page).toHaveURL('/')

  })

  test('able to show correct error message when supplied wrong credentials', async ({ page }) => {
    await page.goto('/login');

    await page.getByTestId('frmLogin_txtUsername').click()
    await page.getByTestId('frmLogin_txtUsername').fill('admin123')
    await page.getByTestId('frmLogin_txtPassword').click()
    await page.getByTestId('frmLogin_txtPassword').fill('admin')
    await page.getByTestId('frmLogin_btnLogin').click()

    await expect(page.getByTestId('frmLogin_divErrorMessage')).toBeVisible()
    await expect(page.getByTestId('frmLogin_divErrorMessage')).toContainText(/Login failed/)

  })

  test.afterEach(async ({ page }) => {
    await page.close()
  })
})
