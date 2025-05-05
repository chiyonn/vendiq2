import styles from './PricingPage.module.css';
import { PricingRow, PricingRowHeader } from '@/components/pricing/PricingRow';
import { useState } from 'react';

const initialItems = [
    {
        asin: "B001",
        image_url: "https://placehold.jp/60x60.png",
        current_price: 1980,
        min_price: 1500,
        last_priced: "2024-12-01T00:00:00.000Z",
        buybox_seller_id: "A",
        num_of_sellers: 8,
        enable: true,
    },
    {
        asin: "B002",
        image_url: "https://placehold.jp/60x60.png",
        current_price: 2480,
        min_price: 2000,
        last_priced: "2024-12-02T00:00:00.000Z",
        buybox_seller_id: "A",
        num_of_sellers: 4,
        enable: false,
    },
    {
        asin: "B003",
        image_url: "https://placehold.jp/60x60.png",
        current_price: 1680,
        min_price: 1600,
        last_priced: "2024-12-03T00:00:00.000Z",
        buybox_seller_id: "B",
        num_of_sellers: 1,
        enable: true,
    },
];

const PricingPage = () => {
    const [items, setItems] = useState(initialItems);

    const updateItem = (index: number, updated: Partial<typeof initialItems[0]>) => {
        setItems(prev =>
            prev.map((item, i) => (i === index ? { ...item, ...updated } : item))
        );
    };

    const handleSaveItem = async (index: number, item: typeof initialItems[0]) => {
        try {
            // 擬似API呼び出し（fetch でもOK）
            const response = await fakeApiUpdate(item);
            // 最新データで行だけ更新
            setItems(prev =>
                prev.map((it, i) => (i === index ? response : it))
            );
        } catch (err) {
            console.error(`保存失敗: ${item.asin}`, err);
            // 任意：エラー表示や再試行処理
        }
    };

    return (
        <div className={styles.container}>
            <PricingRowHeader />
            {items.map((item, index) => (
                <PricingRow
                    key={item.asin}
                    index={index}
                    item={item}
                    onChange={updateItem}
                    onSave={(i: number) => handleSaveItem(i, items[i])}
                />
            ))}
        </div>
    );
};

// モックAPI
const fakeApiUpdate = async (item: typeof initialItems[0]) => {
    // 500ms 待ってから成功レスポンス返す
    return new Promise<typeof initialItems[0]>((resolve) =>
        setTimeout(() => {
            resolve({
                ...item,
                last_priced: new Date().toISOString().split("T")[0], // 更新日だけ変化
            });
        }, 500)
    );
};

export default PricingPage;
