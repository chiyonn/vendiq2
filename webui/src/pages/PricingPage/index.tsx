import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import styles from './PricingPage.module.css';
import { PricingRow, PricingRowHeader, PricingItem } from '@/components/pricing/PricingRow';

const PricingPage = () => {
    const [items, setItems] = useState<PricingItem[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchItems = async () => {
            try {
                const res = await fetch('/pricer/pricings');
                if (!res.ok) throw new Error(`HTTP error ${res.status}`);
                const data: PricingItem[] = await res.json();
                setItems(data);
            } catch (err: any) {
                setError(err.message ?? 'データ取得に失敗しました');
            } finally {
                setLoading(false);
            }
        };

        fetchItems();
    }, []);

    const updateItem = (index: number, updated: Partial<PricingItem>) => {
        setItems(prev =>
            prev.map((item, i) => (i === index ? { ...item, ...updated } : item))
        );
    };

    const handlePricingNow = async () => {
        try {
            const res = await fetch(`/pricer/queue`, {
                method: 'POST',
            });
            if (!res.ok) throw new Error(`HTTP error ${res.status}`);
            alert("価格調整を予約しました。");
        } catch (err) {
            console.error("failed to post new queue", err);
            alert("価格調整のリクエストに失敗しました");
        }
    };

    const handleSaveItem = async (index: number, item: PricingItem) => {
        try {
            const res = await fetch(`/pricer/pricing/${item.ASIN}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(item),
            });
            if (!res.ok) throw new Error(`HTTP error ${res.status}`);
            const updatedItem = await res.json();
            setItems(prev =>
                prev.map((it, i) => (i === index ? updatedItem : it))
            );
        } catch (err) {
            console.error(`保存失敗: ${item.ASIN}`, err);
        }
    };

    if (loading) return <div>読み込み中...</div>;
    if (error) return <div>エラー: {error}</div>;

    return (
        <div className={styles.container}>
            <div className={styles.controller}>
                <button onClick={handlePricingNow}>今すぐ価格調整をする</button>
                <Link to="/pricing/queues">予約を見る</Link>
            </div>
            <PricingRowHeader />
            {items.map((item, index) => (
                <PricingRow
                    key={item.ASIN}
                    index={index}
                    item={item}
                    onChange={updateItem}
                    onSave={() => handleSaveItem(index, items[index])}
                />
            ))}
        </div>
    );
};

export default PricingPage;
