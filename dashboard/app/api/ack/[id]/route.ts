import mongoClient from '@/lib/mongodb';
import { ObjectId } from 'mongodb';
import { type NextRequest , NextResponse } from 'next/server';

export async function GET(
    request: NextRequest,
    { params }: { params: { id: string } }
  ) {
    try {
        const client = await mongoClient;
        const db = client.db("AlertSimAndRemediation");
        const alertsCollection = db.collection("Alerts");
        // update and add Acknowledged at the same time
        const updatedAlert = await alertsCollection.findOneAndUpdate(
            { _id: new ObjectId(params.id) },
            { $set: { Acknowledged: true, AcknowledgedAt: new Date() } },
            { returnDocument: 'after' }
        );
        if (!updatedAlert) {
            return NextResponse.json({ message: 'Alert not found' }, { status: 404 });
        }
        // console.table(updatedAlert);
        console.log(`Alert ${params.id} acknowledged on the server`);
        return NextResponse.json(updatedAlert, { status: 200 });
        
    } catch (e) {
        console.error(e);
        return NextResponse.json({ message: 'An error occurred' }, { status: 500 });
    }
  }