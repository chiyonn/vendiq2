import styles from './PricingRow.module.css';
import { ChangeEvent } from 'react';

export type PricingItem = {
    asin: string;
    mainImageUrl: string;
    minPrice?: number;
    maxPrice?: number;
    numOfSellers: number;
    buyboxPrice: number;
    buyboxSellerId: string;
    autoPricing: boolean;
    createdAt: string;
    updatedAt: string;
    deletedAt: string | null;
};

type Props = {
    index: number;
    item: PricingItem;
    onChange: (index: number, updated: Partial<PricingItem>) => void;
    onSave: (index: number, item: PricingItem) => void;
};

const my_seller_id = "A";

export const PricingRowHeader = () => {
    return (
        <div className={styles.headerRow}>
            <div>ASIN</div>
            <div>画像</div>
            <div>現在価格</div>
            <div>最低価格</div>
            <div>セラー数</div>
            <div>カート取得</div>
            <div>自動</div>
            <div>最終更新</div>
            <div></div>
        </div>
    );
};

export const PricingRow = ({ index, item, onChange, onSave }: Props) => {
    const handleMinPriceChange = (e: ChangeEvent<HTMLInputElement>) => {
        const value = Number(e.target.value);
        if (!isNaN(value)) {
            onChange(index, { minPrice: value });
        }
    };

    const handleEnableToggle = (e: ChangeEvent<HTMLInputElement>) => {
        onChange(index, { autoPricing: e.target.checked });
    };

    return (
        <div className={styles.itemRow}>
            <div>{item.asin}</div>
            <div>
                <img src={item.mainImageUrl} alt={item.asin} className={styles.image} />
            </div>
            <div>¥{item.buyboxPrice}</div>
            <div>
                <input
                    type="number"
                    value={item.buyboxPrice}
                    onChange={handleMinPriceChange}
                    className={styles.input}
                />
            </div>
            <div>{item.numOfSellers}</div>
            <div>{item.buyboxSellerId === my_seller_id ? "YES" : "NO"}</div>
            <div>
                <input
                    type="checkbox"
                    checked={item.autoPricing}
                    onChange={handleEnableToggle}
                />
            </div>
            <div>{new Date(item.updatedAt).toLocaleString()}</div>
            <div>
                <button onClick={() => onSave(index, item)}>保存</button>
            </div>
        </div>
    );
};
