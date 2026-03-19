import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import Login from './index'

const mockedNavigate = vi.fn()
const mockedLogin = vi.fn()
const mockedSetAuthenticated = vi.fn()
const mockedSetUser = vi.fn()

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockedNavigate,
  }
})

vi.mock('../../stores/useAuthStore', () => ({
  useAuthStore: () => ({
    login: mockedLogin,
    setIsAuthenticated: mockedSetAuthenticated,
    setUser: mockedSetUser,
  }),
}))

describe('Login page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('submits credentials and navigates to games on successful login', async () => {
    mockedLogin.mockResolvedValue(true)

    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>,
    )

    const user = userEvent.setup()
    const passwordInput = document.querySelector('input[name="password"]') as HTMLInputElement | null

    expect(passwordInput).not.toBeNull()

    await user.type(screen.getByRole('textbox', { name: 'Email' }), 'tester@example.com')
    await user.type(passwordInput!, 'pw123456789')
    await user.click(screen.getByRole('button', { name: 'Log In' }))

    expect(mockedLogin).toHaveBeenCalledWith('tester@example.com', 'pw123456789')
    expect(mockedSetAuthenticated).toHaveBeenCalledWith(true)
    expect(mockedSetUser).toHaveBeenCalledWith({ id: 'tester@example.com', username: 'tester' })
    expect(mockedNavigate).toHaveBeenCalledWith('/')
  })

  it('does not navigate when login fails', async () => {
    mockedLogin.mockResolvedValue(false)

    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>,
    )

    const user = userEvent.setup()
    const passwordInput = document.querySelector('input[name="password"]') as HTMLInputElement | null

    expect(passwordInput).not.toBeNull()

    await user.type(screen.getByRole('textbox', { name: 'Email' }), 'tester@example.com')
    await user.type(passwordInput!, 'wrong')
    await user.click(screen.getByRole('button', { name: 'Log In' }))

    expect(mockedLogin).toHaveBeenCalledWith('tester@example.com', 'wrong')
    expect(mockedNavigate).not.toHaveBeenCalled()
    expect(mockedSetAuthenticated).not.toHaveBeenCalled()
    expect(mockedSetUser).not.toHaveBeenCalled()
  })
})
