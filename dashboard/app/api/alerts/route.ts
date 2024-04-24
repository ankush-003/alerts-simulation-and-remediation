import mongoClient from '@/lib/mongodb';
import { type NextRequest , NextResponse } from 'next/server';

export async function GET(req: NextRequest) {
    try {

        const client = await mongoClient;
        const db = client.db("AlertSimAndRemediation");
        const alertsCollection = db.collection("Alerts");

        const alerts = await alertsCollection.find().toArray();

        if (!alerts) {
            return NextResponse.json({ message: 'No alerts found' }, { status: 404 });
        }

        return NextResponse.json(alerts, { status: 200 });
    } catch (e) {
        console.error(e);
        return NextResponse.json({ message: 'An error occurred' }, { status: 500 });
    }
}