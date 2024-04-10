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
                                        {/* <p className="text-center">
                                            Alert ID: {row.getValue("id")}
                                        </p>
                                        <p className="text-center">
                                            Category : {row.getValue("category")}
                                        </p>

                                        <p className="text-center">
                                            <p className="font-bold">Node Runtime Metrics:</p>
                                            <ul>
                                                <li>Number of Goroutines: {runtime_metrics.num_goroutine}</li>
                                                <li>CPU Usage: {runtime_metrics.cpu_usage}</li>
                                                <li>RAM Usage: {runtime_metrics.ram_usage}</li>
                                            </ul>
                                        </p>

                                        Created At: {format(date, "yyyy-MM-dd HH:mm:ss")} */}
                                        <div className="grid border rounded-lg gap-2 p-4 m-2">
                                            <div className="grid grid-cols-2 items-center gap-4">
                                                <div className="font-bold">Alert ID:</div>
                                                <Button variant="outline">
                                                    {row.getValue("id")}
                                                </Button>
                                            </div>
                                            <div className="grid grid-cols-2 items-center gap-4">
                                                <div className="font-bold">Category:</div>
                                                <Button variant="outline">
                                                    {row.getValue("category")}
                                                </Button>
                                            </div>
                                            <div className="flex flex-col items-center gap-4">
                                                <div className="font-bold">Node Runtime Metrics:</div>
                                                <div className="grid">
                                                    <ul className="grid gap-2">
                                                        <li className="grid grid-cols-2 items-center gap-4">
                                                            <div className="font-bold">Number of Goroutines:</div>
                                                            <Button variant="outline" > 
                                                                {runtime_metrics.num_goroutine}
                                                            </Button>
                                                        </li>
                                                        <li className="grid grid-cols-2 items-center gap-4">
                                                            <div className="font-bold">CPU Usage:</div>
                                                            <Button variant="outline">
                                                                {runtime_metrics.cpu_usage}
                                                            </Button>
                                                        </li>
                                                        <li className="grid grid-cols-2 gap-4 items-center">
                                                            <div className="font-bold">RAM Usage:</div>
                                                            <Button variant="outline">
                                                                {runtime_metrics.ram_usage}
                                                            </Button>
                                                        </li>
                                                    </ul>
                                                </div>
                                            </div>
                                            <div className="grid grid-cols-2 items-center gap-4">
                                                <div className="font-bold">Created At:</div>
                                                <Button variant="outline">
                                                    {format(date, "yyyy-MM-dd HH:mm:ss")}
                                                </Button>
                                            </div>
                                        </div>
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