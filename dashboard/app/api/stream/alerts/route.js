export const dynamic = 'force-dynamic'

import { Redis } from 'ioredis';
const redisClient = new Redis(process.env.REDIS_URL);

redisClient.on('error', (err) => console.error('Redis Client Error', err));
redisClient.on('connect', () => console.log('Connected to Redis'));

export async function GET() {
    const stream = 'alerts';
    
    try {
        // try {
        //     await redisClient.xgroup('CREATE', stream, groupName, '$', 'MKSTREAM');
        // } catch (error) {
        //     console.error('Error creating group', error);
        // }
        const customReadable = new ReadableStream({
            async start(controller) {
                controller.enqueue(`data: ${JSON.stringify({ message: 'Connected' })}\n\n`)
                let lastId = null;
                while (true) {
                    try {
                        // consumer = await redisClient.xreadgroup('GROUP', groupName, 'dashboard-consumer', 'STREAMS', stream, '>');
                        // if (consumer) {
                        //     console.log('Received messages', consumer);
                        //     const messages = consumer[0][1];
                        //     for (const message of messages) {
                        //         const [id, data] = message;
                        //         controller.enqueue(`data: ${JSON.stringify(data)}\n\n`);
                        //         await redisClient.xack(stream, groupName, id);
                        //     }
                        // }
                        const messages = await redisClient.xread('BLOCK', 0, 'STREAMS', stream, lastId || '$');
                        // console.log('Received messages', messages);
                        // console.table(messages[0]);
                        if (messages && messages.length > 0) {
                            const [streamName, messageEntries] = messages[0];
                            // console.log('Received messages', messageEntries);
                            messageEntries.forEach(([id, fields]) => {
                                // console.log('Processing message', id, fields);
                                const alert = fields.reduce((acc, field, index) => {
                                    if (index % 2 === 0) {
                                        acc[field] = fields[index + 1];
                                    }
                                    return acc;
                                }, {});
                                console.table(alert);
                                controller.enqueue(`data: ${JSON.stringify(alert)}\n\n`);
                                lastId = id; // Update the last processed ID
                            });
                        }
                    } catch (error) {
                        console.error('Error consuming messages', error);
                        controller.enqueue(`data: ${JSON.stringify({ message: 'Error consuming messages' })}\n\n`);
                    }
                }
            },
            async cancel() {
                // await consumer.disconnect();
                console.log('Stream closed');
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
        // if (consumer) {
        //     await consumer.disconnect();
        // }
        return new Response("Internal Server Error", { status: 500 });
    }

}