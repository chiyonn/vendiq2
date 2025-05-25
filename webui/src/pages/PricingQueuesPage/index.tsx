import { useEffect, useState } from 'react';
import styles from './PricingQueuesPage.module.css';

type QueueInfo = {
    name: string;
    messages: number;
    messages_ready: number;
    messages_unacknowledged: number;
};

const PricingQueuesPage = () => {
    const [queues, setQueues] = useState<QueueInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchQueues = async () => {
            try {
                const res = await fetch("/pricer/queues");
                if (!res.ok) throw new Error(`HTTP error ${res.status}`);
                const data: QueueInfo[] = await res.json();
                setQueues(data);
            } catch (err: any) {
                setError(err.message || "Failed to fetch queue data");
            } finally {
                setLoading(false);
            }
        };

        fetchQueues();
    }, []);

    if (loading) return <div>Loading queue info...</div>;
    if (error) return <div>Error: {error}</div>;

    return (
        <div className={styles.container}>
            <h2>価格調整予約</h2>
            <table className={styles.table}>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Total Messages</th>
                        <th>Ready</th>
                        <th>Unacknowledged</th>
                    </tr>
                </thead>
                <tbody>
                    {queues.map((queue) => (
                        <tr key={queue.name}>
                            <td>{queue.name}</td>
                            <td>{queue.messages}</td>
                            <td>{queue.messages_ready}</td>
                            <td>{queue.messages_unacknowledged}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default PricingQueuesPage;
