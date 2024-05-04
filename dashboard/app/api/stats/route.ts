import { NextResponse } from 'next/server';
import clientPromise from "../../../lib/mongodb";
import { hash } from 'argon2';

interface TimeSeries {
    date: string;
    count: number;
}

export async function GET(req: Request) {
    try {
        const client = await clientPromise;
        const db = client.db("AlertSimAndRemediation");
        const alertsCollection = db.collection("Alerts");
        const alerts = await alertsCollection.aggregate([
            {
              $group: {
                _id: {
                  $dateToString: {
                    format: "%Y-%m-%d",
                    date: "$createdAt"
                  }
                },
                value: { $sum: 1 }
              }
            },
            {
              $project: {
                _id: 0,
                day: "$_id",
                value: 1,
              }
            },
            {
              $sort: {
                day: 1
              }
            }
          ]).toArray();
        return NextResponse.json(alerts, { status: 200 });
    } catch (e) {
        console.error(e);
        return NextResponse.json({ message: 'An error occurred' }, { status: 500 });
    }
}