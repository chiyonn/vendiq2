import { useEffect, useState } from 'react';
import styles from './PricingPage.module.css';
import { PricingRow, PricingRowHeader, PricingItem } from '@/components/pricing/PricingRow';

const PricingPage = () => {
    const [items, setItems] = useState<PricingItem[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchItems = async () => {
            try {
                const res = await fetch('/api/pricings');
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

    const handleSaveItem = async (index: number, item: PricingItem) => {
        try {
            const res = await fetch(`/api/pricing/${item.ASIN}`, {
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
