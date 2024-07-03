import { NextResponse, NextRequest } from "next/server";
import clientPromise from "@/lib/mongodb";

interface TimeSeries {
  date: string;
  count: number;
}

interface SeverityCount {
  severity: string;
  count: number;
}

interface CategoryCount {
  category: string;
  count: number;
}

async function getCalendarData(db: any): Promise<TimeSeries[]> {
  return await db.collection("Alerts").aggregate([
    {
      $group: {
        _id: {
          $dateToString: {
            format: "%Y-%m-%d",
            date: "$createdAt",
          },
        },
        value: { $sum: 1 },
      },
    },
    {
      $project: {
        _id: 0,
        day: "$_id",
        value: 1,
      },
    },
    {
      $sort: {
        day: 1,
      },
    },
  ]).toArray();
}

async function getSeverityData(db: any): Promise<SeverityCount[]> {
  return await db.collection("Alerts").aggregate([
    {
      $group: {
        _id: "$severity",
        count: { $sum: 1 }
      }
    },
    {
      $project: {
        severity: "$_id",
        count: 1,
        _id: 0
      }
    },
    {
      $sort: { severity: 1 }
    }
  ]).toArray();
}

async function getCategoryData(db: any): Promise<CategoryCount[]> {
  return await db.collection("Alerts").aggregate([
    {
      $group: {
        _id: "$category",
        count: { $sum: 1 }
      }
    },
    {
      $project: {
        category: "$_id",
        count: 1,
        _id: 0
      }
    },
    {
      $sort: { category: 1 }
    }
  ]).toArray();
}

export async function GET(req: NextRequest) {
  try {
    const client = await clientPromise;
    const type = req.nextUrl.searchParams.get("type");
    const db = client.db("AlertSimAndRemediation");

    let result;
    switch (type) {
      case "calendar":
        result = await getCalendarData(db);
        break;
      case "severity":
        result = await getSeverityData(db);
        break;
      case "category":
        result = await getCategoryData(db);
        break;
      default:
        return NextResponse.json({ message: "Invalid type" }, { status: 400 });
    }

    return NextResponse.json(result, { status: 200 });
  } catch (e) {
    console.error(e);
    return NextResponse.json({ message: "An error occurred" }, { status: 500 });
  }
}