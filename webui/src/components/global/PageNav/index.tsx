import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import styles from './PageNav.module.css';

type Service = {
    name: string;
    host: string;
    link: string;
};

const PageNav = () => {
    const services: Service[] = [
        { name: '価格設定', host: '/pricer', link: '/pricings' },
        { name: 'リサーチ', host: '/researcher', link: '/research' },
    ];

    const [statuses, setStatuses] = useState<Record<string, 'ok' | 'fail'>>({});

    useEffect(() => {
        services.forEach(async (service) => {
            try {
                const res = await fetch(`${service.host}/health`);
                const data = await res.json();
                setStatuses((prev) => ({
                    ...prev,
                    [service.host]: data.status === 'ok' ? 'ok' : 'fail',
                }));
            } catch {
                setStatuses((prev) => ({
                    ...prev,
                    [service.host]: 'fail',
                }));
            }
        });
    }, []);

    return (
        <div className={styles.container}>
            {services.map((service) => (
                <div className={styles.row} key={service.host}>
                    <Link to={service.link}>{service.name}</Link>
                    <p className={statuses[service.host] === 'ok' ? styles.green : styles.red}>
                        ●
                    </p>
                </div>
            ))}
        </div>
    );
};

export default PageNav;
