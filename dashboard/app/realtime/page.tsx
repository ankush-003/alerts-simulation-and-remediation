"use client"

import { useEffect, useState } from "react"

export default function Realtime() {
    useEffect(() => {
        console.log("Realtime page loaded")
        const eventSource = new EventSource("/api/stream")
        eventSource.addEventListener("message", (event) => {
            const data = JSON.parse(event.data)
            console.table(data)
        })

        return () => {
            eventSource.close()
        }
    }, []);

    return (
        <div className="flex flex-col items-center justify-center h-[40rem]">
            <h1>Realtime</h1>
        </div>
    )
}