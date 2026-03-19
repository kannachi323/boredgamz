import { test, expect } from '@playwright/test'

test.describe('Auth flows', () => {
  test('user can sign up and gets redirected to login', async ({ page }) => {
    await page.route('http://localhost:3000/signup', async (route) => {
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({ ok: true }),
      })
    })

    await page.goto('/signup')

    await page.getByRole('textbox', { name: 'Email' }).fill(`new-user-${Date.now()}@example.com`)
    await page.getByRole('textbox', { name: 'Password' }).fill('pw123456789')
    await page.getByRole('button', { name: 'Sign Up' }).click()

    await expect(page).toHaveURL(/\/login$/)
  })

  test('user can log in and gets redirected to home games landing', async ({ page }) => {
    await page.route('**/api/login', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ ok: true }),
      })
    })

    await page.goto('/login')

    await page.getByRole('textbox', { name: 'Email' }).fill(`existing-user-${Date.now()}@example.com`)
    await page.getByRole('textbox', { name: 'Password' }).fill('pw123456789')
    await page.getByRole('button', { name: 'Log In' }).click()

    await expect(page).toHaveURL(/\/$/)
    await expect(page.getByRole('heading', { name: 'Our Games' })).toBeVisible()
  })

  test('failed login stays on login page', async ({ page }) => {
    await page.route('**/api/login', async (route) => {
      await route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'unauthorized' }),
      })
    })

    await page.goto('/login')

    await page.getByRole('textbox', { name: 'Email' }).fill('wrong@example.com')
    await page.getByRole('textbox', { name: 'Password' }).fill('wrong-password')
    await page.getByRole('button', { name: 'Log In' }).click()

    await expect(page).toHaveURL(/\/login$/)
  })
})
