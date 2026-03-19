import { useCallback, useRef, useState } from 'react'
import type { ReactNode, RefObject } from 'react'
import { useClickOutside } from '../hooks/useClickOutside'

type DropdownItem = {
  key: string
  content: ReactNode
  onClick?: () => void
}

type DropdownProps = {
  trigger: ReactNode
  items: DropdownItem[]
  className?: string
  menuClassName?: string
  align?: 'left' | 'right' | 'center'
  openOnHover?: boolean
}

function alignClass(align: DropdownProps['align']) {
  switch (align) {
    case 'right':
      return 'right-0'
    case 'center':
      return 'left-1/2 -translate-x-1/2'
    default:
      return 'left-0'
  }
}

export function Dropdown({
  trigger,
  items,
  className,
  menuClassName,
  align = 'left',
  openOnHover = false,
}: DropdownProps) {
  const [open, setOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement>(null)
  const closeTimeoutRef = useRef<NodeJS.Timeout | null>(null)

  // Only use click-outside detection for click-based dropdowns
  useClickOutside(
    rootRef as RefObject<HTMLElement>,
    () => {
      if (!openOnHover) {
        setOpen(false)
      }
    }
  )

  const closeMenu = useCallback(() => setOpen(false), [])

  const handleMouseLeave = useCallback(() => {
    if (closeTimeoutRef.current) clearTimeout(closeTimeoutRef.current)
    // Add a small delay to allow moving mouse to menu items
    closeTimeoutRef.current = setTimeout(() => setOpen(false), 150)
  }, [])

  const handleMouseEnter = useCallback(() => {
    if (closeTimeoutRef.current) clearTimeout(closeTimeoutRef.current)
    setOpen(true)
  }, [])

  return (
    <div
      ref={rootRef}
      className={`relative inline-flex ${className || ''}`}
      onMouseEnter={openOnHover ? handleMouseEnter : undefined}
      onMouseLeave={openOnHover ? handleMouseLeave : undefined}
    >
      <button
        type="button"
        className="cursor-pointer"
        onClick={() => setOpen((prev) => !prev)}
      >
        {trigger}
      </button>

      {open && items.length > 0 && (
        <div
          className={`absolute top-full mt-2 min-w-36 rounded-xl border border-zinc-800 bg-zinc-900/95 p-1.5 shadow-xl backdrop-blur z-50 ${alignClass(align)} ${menuClassName || ''}`}
        >
          {items.map((item) => (
            <button
              key={item.key}
              type="button"
              className="w-full text-left px-3 py-2 text-zinc-100 rounded-lg hover:bg-zinc-800 transition"
              onClick={() => {
                item.onClick?.()
                closeMenu()
              }}
            >
              {item.content}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
