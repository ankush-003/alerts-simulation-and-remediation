import { ColumnDef } from "@tanstack/react-table"
import { z } from "zod"
import { cn } from "@/lib/utils";
import { parse, format, addDays, subMonths } from 'date-fns';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"

const dateSchema = 'yyyy-MM-dd HH:mm:ss'

const alertSchema = z.object({
    id: z.string(),
    origin: z.string(),
    category: z.string(),
    // severity: z.string(),
    source: z.string(),
    createdAt: z.string(),
    handled: z.boolean(),
    status: z.enum(["open", "ack"]),
    params: z.object({
        num_goroutine: z.number(),
        cpu_usage: z.number(),
        ram_usage: z.number(),
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
        accessorKey: "origin",
    },
    {
        header: "Category",
        accessorKey: "category",
    },
    // {
    //     header: "Severity",
    //     accessorKey: "severity",
    // },
    {
        header: "Source",
        accessorKey: "source",
    },
    {
        header: "Created At",
        accessorKey: "createdAt",
        cell: ({ row }) => {
            const dateString: string = row.getValue("createdAt")
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
        header: "Details",
        accessorKey: "params",
        cell: ({ row }) => {
            const dateString: string = row.getValue("createdAt")
            const date = parse(dateString, dateSchema, new Date())
            const runtime_metrics: Alert["params"] = row.getValue("params")
            console.log(runtime_metrics)

            return (
               <div>
                    {row.getValue("status") === "open" ? (
                        <AlertDialog>
                        <AlertDialogTrigger asChild>
                          <Button variant="outline">Acknowledge</Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Are you sure you want to acknowledge this alert?</AlertDialogTitle>
                            <AlertDialogDescription>
                                Alert ID: {row.getValue("id")}
                                Category : {row.getValue("category")}
                            
                                {/* Node Runtime Metrics:
                                <ul>
                                    <li>Number of Goroutines: {runtime_metrics.num_goroutine}</li>
                                    <li>Allocated Memory: {runtime_metrics.allocated_mem_bytes}</li>
                                    <li>Total Allocated Memory: {runtime_metrics.total_allocated_mem_bytes}</li>
                                    <li>System Memory: {runtime_metrics.sys_mem_bytes}</li>
                                </ul> */}
                                Created At: {format(date, "yyyy-MM-dd HH:mm:ss")}
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction>Continue</AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    ) : (
                        <Button disabled className="bg-green-100">
                            Acknowledged
                        </Button>
                    )}
               </div>
            )
        }
    }
]