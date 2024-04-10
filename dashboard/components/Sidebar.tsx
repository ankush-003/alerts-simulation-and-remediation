"use client"

import { useState } from "react"
import {
    Home,
    FlaskConical,
    Github,
    Box,
    ChevronRight,
    ChevronLeft,
    PlugZap
  } from "lucide-react"
import { Nav } from '@/components/Nav'
import { Button } from '@/components/ui/button'


export default function Sidebar() {
  const [isCollapsed, setIsCollapsed] = useState(true)
  return (
    <div className="relative min-w-[80px] border-r px-3 pb-10 pt-24 bg-background text-foreground">
        <div className="absolute right-[-20px] top-8">
        <Button variant={"secondary"} className="rounded-full p-2 hover:border-2 max-md:hidden hover:border-red-500">
            {isCollapsed ? (
                <ChevronRight
                className="h-6 w-6"
                onClick={() => setIsCollapsed(false)}
                />
            ) : (
                <ChevronLeft
                className="h-6 w-6"
                onClick={() => setIsCollapsed(true)}
                />
            )}
        </Button>
        </div>
        <Nav
            isCollapsed={isCollapsed}
            links={[
              {
                title: "dashboard",
                icon: Home,
                variant: "default",
                href: "/home",
              },
              {
                title: "alert config",
                icon: FlaskConical,
                variant: "default",
                href: "/alert-config",
              },
              {
                title: "realtime alerts",
                icon: PlugZap,
                variant: "default",
                href: "/realtime",
              },
              {
                title: "logs",
                icon: Box,
                variant: "default",
                href: "/logs",
              },
              {
                title: "about",
                icon: Github,
                variant: "default",
                href: "https://github.com/ankush-003"
              }
            ]}
          />
    </div>
  )
}