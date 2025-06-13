
import { createRoot } from 'react-dom/client'
import React from 'react'
import './index.css'
import router from './router.tsx'
import { RouterProvider } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthProvider.tsx'

createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <AuthProvider>
            <RouterProvider router={router} />
        </AuthProvider>
      
    </React.StrictMode>
)
