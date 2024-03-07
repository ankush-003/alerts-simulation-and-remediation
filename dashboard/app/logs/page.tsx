import React from 'react'
import Alert from '@/components/Alert'

import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"

const fakeData = [
  {
    id: "1",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "2",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "3",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "4",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "5",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "6",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "7",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
  {
    id: "8",
    nodeID: "1",
    description: "Low Disk Space",
    severity: "critical",
    source: "disk",
    createdAt: new Date().toDateString()
  },
]

export default function logs() {
  return (
    <div className='p-2'>
      <div className='mt-4 p-8'>
        <Carousel>
          <CarouselContent>
            {fakeData.map((item) => (
              <CarouselItem key={item.id} className="md:basis-1/3 lg:basis-1/4">
                <Alert key={item.id} {...item} />
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselPrevious />
          <CarouselNext />
        </Carousel>
      </div>

    </div>
  )
}