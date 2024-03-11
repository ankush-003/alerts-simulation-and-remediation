export const dynamic = 'force-dynamic'

import { Kafka } from "kafkajs"
import { encode } from "punycode"
import { json } from "stream/consumers"

console.table(process.env)

const kafka = new Kafka({
    brokers: [process.env.KAFKA_URL],
    ssl: true,
    sasl: {
        mechanism: "scram-sha-512",
        username: process.env.KAFKA_USERNAME,
        password: process.env.KAFKA_PASSWORD
    }
})

const c = kafka.consumer({ groupId: "dashboard" })

async function consume(consumer, controller, encoder) {
    console.log("Consuming messages")
    await consumer.connect();
    await consumer.subscribe({ topic: process.env.KAFKA_TOPIC, fromBeginning: false });
    await consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
            console.log("Received message", JSON.stringify(message));
            controller.enqueue(encoder.encode(`data: ${JSON.stringify(message)}\n\n`));
        }
    })
    // await consumer.disconnect();
}

export async function GET() {
    const encoder = new TextEncoder();
    const customReadable = new ReadableStream({
        async start(controller) {
            controller.enqueue(encoder.encode(`data: ${JSON.stringify({ message: "Connected" })}\n\n`));
            try {
                await consume(c, controller, encoder);
            } catch (error) {
                console.error("Error consuming messages", error);
            }
        },
        cancel() {
            // Handle cancellation if needed
            console.log("Stream was cancelled");
        }
    });

    return new Response(customReadable, {
        headers: {
            'Content-Type': 'text/event-stream',
            Connection: 'keep-alive',
            'Cache-Control': 'no-cache, no-transform',
            'Content-Encoding': 'none'
        }
    });
}

