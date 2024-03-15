export const dynamic = 'force-dynamic'

import { Kafka, logLevel } from "kafkajs"
// import { encode } from "punycode"
// import { json } from "stream/consumers"

const kafka = new Kafka({
    brokers: [process.env.KAFKA_URL],
    ssl: true,
    sasl: {
        mechanism: "scram-sha-256",
        username: process.env.KAFKA_USERNAME,
        password: process.env.KAFKA_PASSWORD
    },
    logLevel: logLevel.INFO,
})

// const consumer = kafka.consumer({ groupId: "dashboard-10" })

// consumer.on("consumer.group_join", async () => {
//     console.log("Consumer connected to kafka broker!");
// });

// // on disconnect from kafka broker
// consumer.on("consumer.stop", () => {
//     console.log("Consumer disconnected from kafka broker!");
// });

// await consumer.connect();
// await consumer.subscribe({ topics: ["alerts"] });
// await consumer.run({
//     eachMessage: ({ topic, partition, message }) => {
//         console.log("Received message", JSON.parse(message.value.toString()));
//         messages.push(JSON.parse(message.value.toString()));
//     }
// })

export async function GET() {
    const encoder = new TextEncoder();
    let consumer;
    try {
        consumer = kafka.consumer({ groupId: "dashboard-10" });
        await consumer.connect();
        await consumer.subscribe({ topics: ["alerts"] });

        const customReadable = new ReadableStream({
            async start(controller) {
                await consumer.run({
                    eachMessage: ({ topic, partition, message }) => {
                        console.log("Received message", JSON.parse(message.value.toString()));
                        controller.enqueue(encoder.encode(`data: ${JSON.stringify(JSON.parse(message.value.toString()))}\n\n`));
                    }
                })
            },
            async cancel() {
                await consumer.disconnect();
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
    } catch (error) {
        console.error(error);
        if (consumer) {
            await consumer.disconnect();
        }
        return new Response("Internal Server Error", { status: 500 });
    }

}