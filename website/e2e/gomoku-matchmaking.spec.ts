import { test, expect, Page, BrowserContext } from '@playwright/test'

async function joinGomokuLobbyAsGuest(page: Page, color: 'white' | 'black') {
  await page.goto('/games/gomoku')
  
  // Wait for auth check to complete and guest session to be created
  await expect(page.getByRole('heading', { name: 'Gomoku' })).toBeVisible({ timeout: 10_000 })

  if (color === 'white') {
    await page.getByAltText('white stone').click()
  } else {
    await page.getByAltText('black stone').click()
  }

  await page.getByRole('button', { name: 'Play Now' }).click()
}

async function createMatch(browser: import('@playwright/test').Browser) {
  const contextA = await browser.newContext()
  const contextB = await browser.newContext()

  const pageA = await contextA.newPage()
  const pageB = await contextB.newPage()

  await Promise.all([
    joinGomokuLobbyAsGuest(pageA, 'black'),
    joinGomokuLobbyAsGuest(pageB, 'white'),
  ])

  await expect.poll(() => new URL(pageA.url()).pathname, {
    message: 'player A did not reach a gomoku room',
    timeout: 60_000,
  }).toMatch(/^\/games\/gomoku\/gomoku-/)

  await expect.poll(() => new URL(pageB.url()).pathname, {
    message: 'player B did not reach a gomoku room',
    timeout: 60_000,
  }).toMatch(/^\/games\/gomoku\/gomoku-/)

  await expect.poll(() => {
    const pathA = new URL(pageA.url()).pathname
    const pathB = new URL(pageB.url()).pathname
    return pathA === pathB
  }, {
    message: 'players did not land in the same room route',
    timeout: 30_000,
  }).toBe(true)

  await expect(pageA.getByAltText('gomoku board')).toBeVisible()
  await expect(pageB.getByAltText('gomoku board')).toBeVisible()

  return { contextA, contextB, pageA, pageB }
}

test.describe('Gomoku multiplayer matchmaking', () => {
  test('two guest users can join queue and enter the same room', async ({ browser }) => {
    const { contextA, contextB, pageA, pageB } = await createMatch(browser)

    try {
      await expect(pageA.getByAltText('gomoku board')).toBeVisible()
      await expect(pageB.getByAltText('gomoku board')).toBeVisible()
    } finally {
      await closeContext(contextA)
      await closeContext(contextB)
    }
  })

  test('a disconnected guest can retry and reconnect to an active match', async ({ browser }) => {
    const { contextA, contextB, pageA } = await createMatch(browser)

    try {
      const roomPathBeforeDrop = new URL(pageA.url()).pathname

      await contextA.setOffline(true)

      await expect(pageA.getByRole('button', { name: 'Retry' })).toBeVisible({ timeout: 25_000 })

      await contextA.setOffline(false)

      await pageA.getByRole('button', { name: 'Retry' }).click()

      await expect.poll(() => new URL(pageA.url()).pathname, {
        message: 'reconnected player did not stay in gomoku room route',
        timeout: 25_000,
      }).toBe(roomPathBeforeDrop)

      await expect(pageA.getByAltText('gomoku board')).toBeVisible()
      await expect(pageA.getByRole('button', { name: 'Retry' })).not.toBeVisible({ timeout: 25_000 })
    } finally {
      await closeContext(contextA)
      await closeContext(contextB)
    }
  })
})

async function closeContext(context: BrowserContext) {
  try {
    await context.close()
  } catch {
    // no-op
  }
}
