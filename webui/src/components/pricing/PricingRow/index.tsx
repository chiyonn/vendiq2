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

export const PricingRowHeader = () => {
    return (
        <div className={styles.headerRow}>
            <div>ASIN</div>
            <div>最低価格</div>
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
                <input
                    type="number"
                    value={item.minPrice ?? ''}
                    onChange={handleMinPriceChange}
                    className={styles.input}
                />
            </div>
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
