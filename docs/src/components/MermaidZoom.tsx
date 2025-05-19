'use client'

import { cn } from '@/lib/utils'
import React from 'react'

import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch'
import { Button } from './ui/button'
import {
  ArrowPathIcon,
  MagnifyingGlassMinusIcon,
  MagnifyingGlassPlusIcon,
} from '@heroicons/react/24/outline'

interface MermaidZoomProps extends React.ComponentPropsWithoutRef<'svg'> {}

const MermaidZoom: React.FC<MermaidZoomProps> = (props) => (
  <div className="relative flex items-center justify-center">
    <TransformWrapper centerOnInit initialScale={0.9} minScale={0.75}>
      {({ zoomIn, zoomOut, resetTransform }) => (
        <React.Fragment>
          <div className="absolute right-2 top-2 z-10 space-x-1">
            <Button
              aria-label="Zoom In"
              variant="outline"
              size="icon"
              onClick={() => zoomIn()}
            >
              <MagnifyingGlassPlusIcon className="size-5" />
            </Button>
            <Button
              aria-label="Zoom Out"
              variant="outline"
              size="icon"
              onClick={() => zoomOut()}
            >
              <MagnifyingGlassMinusIcon className="size-5" />
            </Button>
            <Button
              aria-label="Reset Zoom"
              variant="outline"
              size="icon"
              onClick={() => resetTransform()}
            >
              <ArrowPathIcon className="size-5" />
            </Button>
          </div>
          <TransformComponent
            wrapperClass="bg-white rounded-lg !w-full cursor-move"
            contentClass="!w-full !h-full"
          >
            <svg
              {...props}
              className={cn('mx-auto px-2 py-4', props.className)}
            />
          </TransformComponent>
        </React.Fragment>
      )}
    </TransformWrapper>
  </div>
)

export default MermaidZoom
