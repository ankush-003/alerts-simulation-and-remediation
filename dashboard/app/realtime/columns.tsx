import { ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { cn } from "@/lib/utils";
import { parse, format, addDays, subMonths } from 'date-fns';

const dateSchema = 'yyyy-MM-dd HH:mm:ss'

const alertSchema = z.object({
    id: z.string(),
    node_id: z.string(),
    description: z.string(),
    severity: z.string(),
    source: z.string(),
    created_at: z.string(),
    status: z.enum(["open", "ack"]),
    runtime_metrics: z.object({
        num_goroutine: z.number(),
        allocated_mem_bytes: z.number(),
        total_allocated_mem_bytes: z.number(),
        sys_mem_bytes: z.number(),
    }),
})

export type Alert = z.infer<typeof alertSchema>;

export const columns: ColumnDef<Alert>[] = [
    {
        header: "Id",
        accessorKey: "id",
    },
    {
        header: "Node Id",
        accessorKey: "node_id",
    },
    {
        header: "Description",
        accessorKey: "description",
    },
    {
        header: "Severity",
        accessorKey: "severity",
    },
    {
        header: "Source",
        accessorKey: "source",
    },
    {
        header: "Created At",
        accessorKey: "created_at",
        cell: ({ row }) => {
            const dateString:string = row.getValue("created_at")
            const date = parse(dateString, dateSchema, new Date())

            return (
                <span className="text-gray-500">
                    {format(date, "yyyy-MM-dd HH:mm:ss")}
                </span>
            )
        }
    },
    {
        header: "Status",
        accessorKey: "status",
    },
    {
        header: "Num Goroutine",
        accessorKey: "runtime_metrics.num_goroutine",
    },
    {
        header: "Allocated Mem",
        accessorKey: "runtime_metrics.allocated_mem_bytes",
    },
    {
        header: "Total Allocated Mem",
        accessorKey: "runtime_metrics.total_allocated_mem_bytes",
    },
    {
        header: "Sys Mem",
        accessorKey: "runtime_metrics.sys_mem_bytes",
    },
]