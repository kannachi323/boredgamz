import { useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'

type CarouselItem = {
  key: string
  image: string
  alt: string
}

type CarouselProps = {
  items: CarouselItem[]
  autoSlideMs?: number
  children?: ReactNode
  showIndicators?: boolean
}

export function Carousel({
  items,
  autoSlideMs = 3200,
  children,
  showIndicators = true,
}: CarouselProps) {
  const [index, setIndex] = useState(0)
  const [isTransitioning, setIsTransitioning] = useState(true)

  const slides = useMemo(() => {
    if (items.length === 0) return []
    if (items.length === 1) return items
    return [...items, ...items]
  }, [items])

  useEffect(() => {
    setIsTransitioning(true)

    setIndex(0)
  }, [items.length])

  useEffect(() => {
    if (items.length <= 1) return

    const timer = window.setInterval(() => {
      setIndex((prev) => prev + 1)
    }, autoSlideMs)

    return () => window.clearInterval(timer)
  }, [autoSlideMs, items.length])

  if (slides.length === 0) return null

  const activeRealIndex = items.length <= 1 ? 0 : index % items.length

  function handleTransitionEnd() {
    if (items.length <= 1) return

    if (index >= items.length) {
      setIsTransitioning(false)
      setIndex((prev) => prev - items.length)
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          setIsTransitioning(true)
        })
      })
    }
  }

  return (
    <div className="relative h-full w-full overflow-hidden">
      <div
        className="flex h-full"
        style={{
          transform: `translateX(-${index * 100}%)`,
          transition: isTransitioning ? 'transform 700ms ease' : 'none',
        }}
        onTransitionEnd={handleTransitionEnd}
      >
        {slides.map((slide, idx) => (
          <div key={`${slide.key}-${idx}`} className="min-w-full h-full">
            <img
              src={slide.image}
              alt={slide.alt}
              className="h-full w-full object-contain object-center"
            />
          </div>
        ))}
      </div>

      {children}

      {showIndicators && items.length > 1 && (
        <div className="absolute bottom-5 left-1/2 -translate-x-1/2 flex gap-2 z-30">
          {items.map((item, dotIdx) => (
            <button
              key={item.key}
              aria-label={`Show ${item.alt}`}
              onClick={() => {
                setIsTransitioning(true)
                setIndex((prev) => {
                  if (items.length <= 1) return dotIdx
                  const current = prev % items.length
                  const forwardDelta = (dotIdx - current + items.length) % items.length
                  return prev + forwardDelta
                })
              }}
              className={`h-2.5 w-8 rounded-full transition ${
                dotIdx === activeRealIndex ? 'bg-white' : 'bg-zinc-400/70 hover:bg-zinc-300'
              }`}
            />
          ))}
        </div>
      )}
    </div>
  )
}
