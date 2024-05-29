"use client"

import { useEffect, useState } from "react"
import { Alert as AlertData, columns } from "./columns"
import { DataTable } from "@/app/realtime/components/DataTable";
import { toast } from "sonner"
// import { TracingBeam } from "@/components/ui/tracing-beam"

async function handleAcknowledge(id: string) {
    const response = await fetch(`/api/ack/${id}`)
    if (!response.ok) {
        // throw new Error("Failed to acknowledge alert")
        console.error("Failed to acknowledge alert")
        // toast.error("Failed to acknowledge alert")
        return false;
    }
    // return response.json()
    // toast.success("Alert acknowledged")
    console.log("Alert acknowledged")
    return true;

}  

export default function Realtime() {
    const [data, setData] = useState<AlertData[]>([]);

    useEffect(() => {
        console.log("Realtime page loaded")
        const eventSource = new EventSource("/api/stream/alerts")
        eventSource.addEventListener("message", (event) => {
            const data = JSON.parse(event.data)
            console.table(data)
            if(data.hasOwnProperty("message")) {
                // console.log("Connection established")
                toast.success("Connection established")
                return
            }

            data.status = "open"
            setData((prevData) => [data, ...prevData])
            toast.error(`New alert: ${data.id}`)
        })

        return () => {
            eventSource.close()
        }
    }, [])

    return (
        <div className="container mx-auto py-10">
            <DataTable columns={columns} data={data} setData={setData} ackfn={handleAcknowledge} />
        </div>
    )
}