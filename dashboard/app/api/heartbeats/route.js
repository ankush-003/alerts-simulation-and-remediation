export const dynamic = 'force-dynamic'
import { Console } from 'console';

// import { createClient } from 'redis';
import { Redis } from 'ioredis';

const redisClient = new Redis(process.env.REDIS_URL);

redisClient.on('error', (err) => console.error('Redis Client Error', err));
redisClient.on('connect', () => console.log('Connected to Redis'));

export async function GET() {
    const stream = 'nodeHeartbeats';
    
    try {

        const customReadable = new ReadableStream({
            async start(controller) {
                controller.enqueue(`data: ${JSON.stringify({ message: 'Connected' })}\n\n`)
                let lastId = null;
                while (true) {
                    try {
                        
                        const messages = await redisClient.xread('BLOCK', 0, 'STREAMS', stream, lastId || '$');
                        
                        if (messages && messages.length > 0) {
                            const [streamName, messageEntries] = messages[0];
                            messageEntries.forEach(([id, fields]) => {
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