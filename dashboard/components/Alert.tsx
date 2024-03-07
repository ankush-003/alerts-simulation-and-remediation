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
                            "text-xl font-medium",
                            severity === "critical" && "text-red-500",
                            severity === "warning" && "text-white",
                            severity === "info" && "text-blue-500",
                        )}
                    >
                        <div className="flex justify-between">
                            <p>{`Alert ${id}`}</p>
                            <AlertCircle className="h-10 w-10" />
                        </div>
                    </CardTitle>
                    <CardDescription className="text-lg">{`${description}`}</CardDescription>
                </CardHeader>
                <CardContent>
                    <p>NodeID</p><br />
                    <p>{nodeID}</p>
                </CardContent>
                <CardFooter className="grid">
                    
                    <div className="w-full">
                        <AlertDialog>
                            <AlertDialogTrigger asChild>
                                <Button variant="outline">view details</Button>
                            </AlertDialogTrigger>
                            <AlertDialogContent>
                                <AlertDialogHeader>
                                    <AlertDialogTitle><p>{`Alert ${id}`}</p> </AlertDialogTitle>
                                    <AlertDialogDescription>
                                    <p>{`createdAt : ${createdAt}`}</p>
                                    </AlertDialogDescription>
                                </AlertDialogHeader>
                                <AlertDialogFooter>
                                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                                    <AlertDialogAction>Resolve</AlertDialogAction>
                                </AlertDialogFooter>
                            </AlertDialogContent>
                        </AlertDialog>
                    </div>
                </CardFooter>
            </Card>
        </div>
    )
}

export default Alert