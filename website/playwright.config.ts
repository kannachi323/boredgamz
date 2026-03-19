import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : undefined,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://127.0.0.1:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: process.env.CI
    ? undefined
    : [
        {
          command:
            'DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5433/boredgamz?sslmode=disable" JWT_SECRET_KEY="local-e2e-secret" go run main.go',
          cwd: '../server',
          url: 'http://127.0.0.1:3000/api/hello',
          reuseExistingServer: true,
          timeout: 120_000,
        },
        {
          command:
            'VITE_SERVER_ROOT="http://127.0.0.1:3000/api" VITE_WEBSOCKET_ROOT="ws://127.0.0.1:3000/api" npm run dev -- --host 127.0.0.1 --port 5173',
          cwd: '.',
          url: 'http://127.0.0.1:5173',
          reuseExistingServer: true,
          timeout: 120_000,
        },
      ],
})
