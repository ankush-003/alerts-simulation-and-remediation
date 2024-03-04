import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { cn } from "@/lib/utils";
import { AlertCircle } from 'lucide-react';
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


interface AlertProp {
    id: string;
    nodeID: string;
    description: string;
    severity: string;
    source: string;
    createdAt: string;
}

function Alert({ id, nodeID, description, severity, source, createdAt }: AlertProp) {
    return (
        <div className="hover:border-2 hover:border-red-500 rounded-lg">
            <Card>
                <CardHeader>
                    <CardTitle
                        className={cn(
                            "text-sm font-medium",
                            severity === "critical" && "text-red-500",
                            severity === "warning" && "text-white",
                            severity === "info" && "text-blue-500",
                        )}
                    >
                        <div className="flex justify-between">
                            <p>{`Alert ${id}`}</p>
                            <AlertCircle className="h-4 w-4" />
                        </div>
                    </CardTitle>
                    <CardDescription>{`${description}`}</CardDescription>
                </CardHeader>
                <CardContent>
                    <p>{`NodeID : ${nodeID}`}</p><br />
                    <p>{`source : ${source}`}</p>
                </CardContent>
                <CardFooter>
                    <p>{`createdAt : ${createdAt}`}</p>
                </CardFooter>
                <AlertDialog>
                    <AlertDialogTrigger asChild>
                        <Button variant="outline">Show Dialog</Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                        <AlertDialogHeader>
                            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                            <AlertDialogDescription>
                                This action cannot be undone. This will permanently delete your
                                account and remove your data from our servers.
                            </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction>Continue</AlertDialogAction>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialog>
            </Card>
        </div>
    )
}

export default Alert